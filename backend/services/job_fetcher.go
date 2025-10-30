package services

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"
)

// Job represents a job posting from external API
type Job struct {
	Title       string `json:"title"`
	Company     string `json:"company"`
	Location    string `json:"location"`
	Description string `json:"description"`
	Salary      string `json:"salary"`
	JobUrl      string `json:"job_url"`
	PostedDate  string `json:"posted_date"`
	JobType     string `json:"job_type"`
}

// AdzunaResponse represents the response from Adzuna API
type AdzunaResponse struct {
	Results []struct {
		Title   string `json:"title"`
		Company struct {
			DisplayName string `json:"display_name"`
		} `json:"company"`
		Location struct {
			DisplayName string `json:"display_name"`
		} `json:"location"`
		Description  string    `json:"description"`
		SalaryMin    float64   `json:"salary_min"`
		SalaryMax    float64   `json:"salary_max"`
		RedirectUrl  string    `json:"redirect_url"`
		Created      time.Time `json:"created"`
		ContractType string    `json:"contract_type"`
	} `json:"results"`
}

// APIResult holds the result from an API call
type APIResult struct {
	Jobs   []Job
	Source string
	Error  error
}

// FetchJobRecommendations fetches real-time jobs using parallel API calls
func FetchJobRecommendations(skills []string, limit int) ([]Job, error) {
	if limit <= 0 || limit > 10 {
		limit = 5
	}

	fmt.Println("\n🚀 Starting Parallel Job Fetch System")
	fmt.Printf("📊 Target: %d jobs from skills: %v\n", limit, skills)

	// Priority 1: Fast, reliable APIs (run in parallel)
	priority1APIs := []func([]string, int) ([]Job, error){
		fetchFromRemoteOK,          // Free, no auth, tech jobs
		fetchFromArbeitnow,         // Free, no auth, EU + US jobs
		fetchFromTheMuse,           // Free, no auth, curated jobs
		fetchFromAdzunaIfAvailable, // Only if credentials available
	}

	// Priority 2: Backup APIs (run in parallel if Priority 1 fails)
	priority2APIs := []func([]string, int) ([]Job, error){
		fetchFromFindwork,           // Free, no auth, tech focus
		fetchFromJoobleIfAvailable,  // Only if API key available
		fetchFromJSearchIfAvailable, // Only if RapidAPI key available
	}

	// Try Priority 1 APIs in parallel
	fmt.Println("\n🔵 Priority 1: Fetching from 3-4 APIs simultaneously...")
	allJobs, totalFetched := fetchParallel(priority1APIs, skills, limit)

	// If we got enough jobs, return them
	if len(allJobs) >= limit {
		fmt.Printf("\n✅ SUCCESS: Got %d jobs from Priority 1 APIs\n", len(allJobs))
		return deduplicateJobs(allJobs, limit), nil
	}

	// If Priority 1 didn't get enough jobs, try Priority 2
	if totalFetched < limit {
		fmt.Printf("\n🟡 Priority 1 only got %d jobs, trying Priority 2 APIs...\n", totalFetched)
		moreJobs, _ := fetchParallel(priority2APIs, skills, limit-totalFetched)
		allJobs = append(allJobs, moreJobs...)
	}

	// Deduplicate and return
	if len(allJobs) > 0 {
		fmt.Printf("\n✅ TOTAL: Fetched %d jobs from all APIs\n", len(allJobs))
		return deduplicateJobs(allJobs, limit), nil
	}

	// Final fallback: generate sample jobs
	fmt.Println("\n📝 All APIs failed, generating sample jobs")
	return generateSampleJobs(skills, limit), nil
}

// fetchParallel runs multiple API fetchers in parallel and collects results
func fetchParallel(apis []func([]string, int) ([]Job, error), skills []string, limit int) ([]Job, int) {
	var wg sync.WaitGroup
	results := make(chan APIResult, len(apis))

	// Launch all API calls in parallel
	for _, apiFunc := range apis {
		wg.Add(1)
		go func(fn func([]string, int) ([]Job, error)) {
			defer wg.Done()
			jobs, err := fn(skills, limit)
			results <- APIResult{
				Jobs:  jobs,
				Error: err,
			}
		}(apiFunc)
	}

	// Wait for all goroutines to complete
	go func() {
		wg.Wait()
		close(results)
	}()

	// Collect all results
	var allJobs []Job
	successCount := 0
	failCount := 0

	for result := range results {
		if result.Error == nil && len(result.Jobs) > 0 {
			allJobs = append(allJobs, result.Jobs...)
			successCount++
			fmt.Printf("  ✓ API returned %d jobs\n", len(result.Jobs))
		} else {
			failCount++
			fmt.Printf("  ✗ API failed or returned 0 jobs\n")
		}
	}

	fmt.Printf("📈 Parallel fetch complete: %d succeeded, %d failed\n", successCount, failCount)
	return allJobs, len(allJobs)
}

// deduplicateJobs removes duplicate jobs based on title + company and limits to desired count
func deduplicateJobs(jobs []Job, limit int) []Job {
	seen := make(map[string]bool)
	unique := make([]Job, 0, limit)

	for _, job := range jobs {
		// Create a unique key based on title and company
		key := strings.ToLower(fmt.Sprintf("%s|%s", job.Title, job.Company))

		if !seen[key] {
			seen[key] = true
			unique = append(unique, job)

			// Stop when we reach the desired limit
			if len(unique) >= limit {
				break
			}
		}
	}

	fmt.Printf("🔧 Deduplication: %d jobs → %d unique jobs\n", len(jobs), len(unique))
	return unique
}

// fetchFromAdzunaIfAvailable fetches from Adzuna only if credentials are available
func fetchFromAdzunaIfAvailable(skills []string, limit int) ([]Job, error) {
	appId := os.Getenv("ADZUNA_APP_ID")
	appKey := os.Getenv("ADZUNA_APP_KEY")

	if appId == "" || appKey == "" {
		return nil, fmt.Errorf("Adzuna credentials not available")
	}

	return fetchFromAdzuna(skills, limit, appId, appKey)
}

// fetchFromJoobleIfAvailable fetches from Jooble only if API key is available
func fetchFromJoobleIfAvailable(skills []string, limit int) ([]Job, error) {
	apiKey := os.Getenv("JOOBLE_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("Jooble API key not available")
	}

	return fetchFromJooble(skills, limit, apiKey)
}

// fetchFromJSearchIfAvailable fetches from JSearch only if RapidAPI key is available
func fetchFromJSearchIfAvailable(skills []string, limit int) ([]Job, error) {
	apiKey := os.Getenv("RAPIDAPI_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("RapidAPI key not available")
	}

	return fetchFromJSearch(skills, limit, apiKey)
}

// fetchFromAdzuna fetches jobs from Adzuna API
func fetchFromAdzuna(skills []string, limit int, appId, appKey string) ([]Job, error) {
	// Build search query from skills
	query := strings.Join(skills, " OR ")
	if len(query) > 200 {
		query = strings.Join(skills[:5], " OR ")
	}

	country := os.Getenv("JOB_COUNTRY")
	if country == "" {
		country = "us"
	}

	apiURL := fmt.Sprintf("https://api.adzuna.com/v1/api/jobs/%s/search/1", country)

	params := url.Values{}
	params.Add("app_id", appId)
	params.Add("app_key", appKey)
	params.Add("results_per_page", fmt.Sprintf("%d", limit))
	params.Add("what", query)
	params.Add("content-type", "application/json")

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var adzunaResp AdzunaResponse
	if err := json.Unmarshal(body, &adzunaResp); err != nil {
		return nil, err
	}

	if len(adzunaResp.Results) == 0 {
		return nil, fmt.Errorf("no jobs found")
	}

	// Convert to our Job struct
	jobs := make([]Job, 0, len(adzunaResp.Results))
	for _, result := range adzunaResp.Results {
		salary := ""
		if result.SalaryMin > 0 || result.SalaryMax > 0 {
			if result.SalaryMin > 0 && result.SalaryMax > 0 {
				salary = fmt.Sprintf("$%.0f - $%.0f", result.SalaryMin, result.SalaryMax)
			} else if result.SalaryMin > 0 {
				salary = fmt.Sprintf("From $%.0f", result.SalaryMin)
			} else {
				salary = fmt.Sprintf("Up to $%.0f", result.SalaryMax)
			}
		}

		description := cleanDescription(result.Description, 200)

		jobs = append(jobs, Job{
			Title:       result.Title,
			Company:     result.Company.DisplayName,
			Location:    result.Location.DisplayName,
			Description: description,
			Salary:      salary,
			JobUrl:      result.RedirectUrl,
			PostedDate:  result.Created.Format("2006-01-02"),
			JobType:     result.ContractType,
		})
	}

	fmt.Printf("  ✓ Adzuna API: %d jobs\n", len(jobs))
	return jobs, nil
}

// fetchFromTheMuse uses The Muse API (free, no auth)
func fetchFromTheMuse(skills []string, limit int) ([]Job, error) {
	apiURL := "https://www.themuse.com/api/public/jobs"
	params := url.Values{}
	params.Add("page", "0")
	params.Add("descending", "true")
	params.Add("api_key", "public")

	// Add category based on skills
	if contains(skills, "javascript") || contains(skills, "react") || contains(skills, "node") {
		params.Add("category", "Software Engineering")
	} else if contains(skills, "python") || contains(skills, "java") {
		params.Add("category", "Data Science")
	}

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(fullURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var museResp struct {
		Results []struct {
			Name    string `json:"name"`
			Company struct {
				Name string `json:"name"`
			} `json:"company"`
			Locations []struct {
				Name string `json:"name"`
			} `json:"locations"`
			Contents        string `json:"contents"`
			PublicationDate string `json:"publication_date"`
			Refs            struct {
				LandingPage string `json:"landing_page"`
			} `json:"refs"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &museResp); err != nil {
		return nil, err
	}

	jobs := make([]Job, 0, min(limit, len(museResp.Results)))
	for i, result := range museResp.Results {
		if i >= limit {
			break
		}

		location := "Remote"
		if len(result.Locations) > 0 {
			location = result.Locations[0].Name
		}

		jobs = append(jobs, Job{
			Title:       result.Name,
			Company:     result.Company.Name,
			Location:    location,
			Description: cleanDescription(result.Contents, 200),
			Salary:      "",
			JobUrl:      result.Refs.LandingPage,
			PostedDate:  result.PublicationDate,
			JobType:     "Full-time",
		})
	}

	fmt.Printf("  ✓ The Muse API: %d jobs\n", len(jobs))
	return jobs, nil
}

// fetchFromJSearch uses JSearch API (RapidAPI)
func fetchFromJSearch(skills []string, limit int, apiKey string) ([]Job, error) {
	query := strings.Join(skills[:min(3, len(skills))], " ")
	apiURL := "https://jsearch.p.rapidapi.com/search"

	params := url.Values{}
	params.Add("query", query+" developer")
	params.Add("num_pages", "1")

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	req, err := http.NewRequest("GET", fullURL, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("X-RapidAPI-Key", apiKey)
	req.Header.Add("X-RapidAPI-Host", "jsearch.p.rapidapi.com")

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var jsearchResp struct {
		Data []struct {
			JobTitle          string `json:"job_title"`
			EmployerName      string `json:"employer_name"`
			JobCity           string `json:"job_city"`
			JobState          string `json:"job_state"`
			JobDescription    string `json:"job_description"`
			JobMinSalary      string `json:"job_min_salary"`
			JobMaxSalary      string `json:"job_max_salary"`
			JobApplyLink      string `json:"job_apply_link"`
			JobPostedDate     string `json:"job_posted_at_datetime_utc"`
			JobEmploymentType string `json:"job_employment_type"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &jsearchResp); err != nil {
		return nil, err
	}

	jobs := make([]Job, 0, min(limit, len(jsearchResp.Data)))
	for i, result := range jsearchResp.Data {
		if i >= limit {
			break
		}

		location := result.JobCity
		if result.JobState != "" {
			location = fmt.Sprintf("%s, %s", result.JobCity, result.JobState)
		}

		salary := ""
		if result.JobMinSalary != "" && result.JobMaxSalary != "" {
			salary = fmt.Sprintf("$%s - $%s", result.JobMinSalary, result.JobMaxSalary)
		}

		jobs = append(jobs, Job{
			Title:       result.JobTitle,
			Company:     result.EmployerName,
			Location:    location,
			Description: cleanDescription(result.JobDescription, 200),
			Salary:      salary,
			JobUrl:      result.JobApplyLink,
			PostedDate:  result.JobPostedDate,
			JobType:     result.JobEmploymentType,
		})
	}

	fmt.Printf("  ✓ JSearch API: %d jobs\n", len(jobs))
	return jobs, nil
}

// fetchFromJooble fetches jobs from Jooble API (requires API key)
func fetchFromJooble(skills []string, limit int, apiKey string) ([]Job, error) {
	// Build keywords from top skills
	keywords := strings.Join(skills[:min(5, len(skills))], " ")

	apiURL := "https://jooble.org/api/" + apiKey

	// Jooble expects POST request with JSON body
	requestBody := map[string]interface{}{
		"keywords": keywords,
		"location": "", // Empty for worldwide
		"radius":   "100",
		"page":     "1",
	}

	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return nil, err
	}

	client := &http.Client{Timeout: 15 * time.Second}
	req, err := http.NewRequest("POST", apiURL, strings.NewReader(string(jsonBody)))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("status code: %d, body: %s", resp.StatusCode, string(body))
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var joobleResp struct {
		TotalCount int `json:"totalCount"`
		Jobs       []struct {
			Title    string `json:"title"`
			Location string `json:"location"`
			Snippet  string `json:"snippet"`
			Salary   string `json:"salary"`
			Source   string `json:"source"`
			Type     string `json:"type"`
			Link     string `json:"link"`
			Company  string `json:"company"`
			Updated  string `json:"updated"`
		} `json:"jobs"`
	}

	if err := json.Unmarshal(body, &joobleResp); err != nil {
		return nil, err
	}

	if len(joobleResp.Jobs) == 0 {
		return nil, fmt.Errorf("no jobs found")
	}

	jobs := make([]Job, 0, min(limit, len(joobleResp.Jobs)))
	for i, result := range joobleResp.Jobs {
		if i >= limit {
			break
		}

		// Clean and format data
		company := result.Company
		if company == "" {
			company = result.Source
		}

		location := result.Location
		if location == "" {
			location = "Not specified"
		}

		jobType := result.Type
		if jobType == "" {
			jobType = "Full-time"
		}

		jobs = append(jobs, Job{
			Title:       result.Title,
			Company:     company,
			Location:    location,
			Description: cleanDescription(result.Snippet, 200),
			Salary:      result.Salary,
			JobUrl:      result.Link,
			PostedDate:  result.Updated,
			JobType:     jobType,
		})
	}

	fmt.Printf("  ✓ Jooble API: %d jobs (total available: %d)\n", len(jobs), joobleResp.TotalCount)
	return jobs, nil
}

// fetchFromArbeitnow fetches jobs from Arbeitnow API (free, no auth, EU + US)
func fetchFromArbeitnow(skills []string, limit int) ([]Job, error) {
	apiURL := "https://www.arbeitnow.com/api/job-board-api"

	client := &http.Client{Timeout: 10 * time.Second}
	resp, err := client.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var arbeitResp struct {
		Data []struct {
			Slug        string   `json:"slug"`
			CompanyName string   `json:"company_name"`
			Title       string   `json:"title"`
			Description string   `json:"description"`
			Remote      bool     `json:"remote"`
			URL         string   `json:"url"`
			Tags        []string `json:"tags"`
			JobTypes    []string `json:"job_types"`
			Location    string   `json:"location"`
			CreatedAt   int64    `json:"created_at"`
		} `json:"data"`
	}

	if err := json.Unmarshal(body, &arbeitResp); err != nil {
		return nil, err
	}

	// Filter jobs based on skills match
	jobs := make([]Job, 0, limit)
	skillsLower := make([]string, len(skills))
	for i, s := range skills {
		skillsLower[i] = strings.ToLower(s)
	}

	for _, item := range arbeitResp.Data {
		if len(jobs) >= limit {
			break
		}

		// Check if job matches any skill
		matchFound := false
		titleLower := strings.ToLower(item.Title)
		descLower := strings.ToLower(item.Description)

		for _, skill := range skillsLower {
			if strings.Contains(titleLower, skill) || strings.Contains(descLower, skill) {
				matchFound = true
				break
			}
			// Also check tags
			for _, tag := range item.Tags {
				if strings.Contains(strings.ToLower(tag), skill) {
					matchFound = true
					break
				}
			}
			if matchFound {
				break
			}
		}

		if !matchFound {
			continue
		}

		location := item.Location
		if item.Remote {
			location = "Remote"
		}

		jobType := "Full-time"
		if len(item.JobTypes) > 0 {
			jobType = item.JobTypes[0]
		}

		postedDate := time.Unix(item.CreatedAt, 0).Format("2006-01-02")

		jobs = append(jobs, Job{
			Title:       item.Title,
			Company:     item.CompanyName,
			Location:    location,
			Description: cleanDescription(item.Description, 200),
			Salary:      "",
			JobUrl:      item.URL,
			PostedDate:  postedDate,
			JobType:     jobType,
		})
	}

	fmt.Printf("  ✓ Arbeitnow API: %d jobs\n", len(jobs))
	return jobs, nil
}

// fetchFromFindwork fetches jobs from Findwork API (free, no auth, tech focus)
func fetchFromFindwork(skills []string, limit int) ([]Job, error) {
	// Findwork API - free tier, no auth
	apiURL := "https://findwork.dev/api/jobs/"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// Add required headers
	req.Header.Add("Authorization", "Token test-token") // Public test token

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var findworkResp struct {
		Results []struct {
			ID             int      `json:"id"`
			Role           string   `json:"role"`
			CompanyName    string   `json:"company_name"`
			Location       string   `json:"location"`
			Remote         bool     `json:"remote"`
			Description    string   `json:"text"`
			URL            string   `json:"url"`
			EmploymentType string   `json:"employment_type"`
			DatePosted     string   `json:"date_posted"`
			Keywords       []string `json:"keywords"`
		} `json:"results"`
	}

	if err := json.Unmarshal(body, &findworkResp); err != nil {
		return nil, err
	}

	// Filter jobs based on skills match
	jobs := make([]Job, 0, limit)
	skillsLower := make([]string, len(skills))
	for i, s := range skills {
		skillsLower[i] = strings.ToLower(s)
	}

	for _, item := range findworkResp.Results {
		if len(jobs) >= limit {
			break
		}

		// Check if job matches any skill
		matchFound := false
		roleLower := strings.ToLower(item.Role)
		descLower := strings.ToLower(item.Description)

		for _, skill := range skillsLower {
			if strings.Contains(roleLower, skill) || strings.Contains(descLower, skill) {
				matchFound = true
				break
			}
			// Also check keywords
			for _, keyword := range item.Keywords {
				if strings.Contains(strings.ToLower(keyword), skill) {
					matchFound = true
					break
				}
			}
			if matchFound {
				break
			}
		}

		if !matchFound {
			continue
		}

		location := item.Location
		if item.Remote {
			location = "Remote"
		}

		jobs = append(jobs, Job{
			Title:       item.Role,
			Company:     item.CompanyName,
			Location:    location,
			Description: cleanDescription(item.Description, 200),
			Salary:      "",
			JobUrl:      item.URL,
			PostedDate:  item.DatePosted,
			JobType:     item.EmploymentType,
		})
	}

	fmt.Printf("  ✓ Findwork API: %d jobs\n", len(jobs))
	return jobs, nil
}

// fetchFromRemoteOK fetches tech jobs from RemoteOK (free, no auth)
func fetchFromRemoteOK(skills []string, limit int) ([]Job, error) {
	apiURL := "https://remoteok.com/api"

	client := &http.Client{Timeout: 10 * time.Second}
	req, err := http.NewRequest("GET", apiURL, nil)
	if err != nil {
		return nil, err
	}

	// RemoteOK requires user agent
	req.Header.Add("User-Agent", "SmartResume/1.0")

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("status code: %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var remoteResp []map[string]interface{}
	if err := json.Unmarshal(body, &remoteResp); err != nil {
		return nil, err
	}

	// Filter jobs based on skills match
	jobs := make([]Job, 0, limit)
	skillsLower := make([]string, len(skills))
	for i, s := range skills {
		skillsLower[i] = strings.ToLower(s)
	}

	for _, item := range remoteResp {
		if len(jobs) >= limit {
			break
		}

		// Skip first item (it's metadata)
		if _, ok := item["position"]; !ok {
			continue
		}

		// Check if job matches any skill
		position := fmt.Sprintf("%v", item["position"])
		tags := fmt.Sprintf("%v", item["tags"])

		matchFound := false
		for _, skill := range skillsLower {
			if strings.Contains(strings.ToLower(position), skill) ||
				strings.Contains(strings.ToLower(tags), skill) {
				matchFound = true
				break
			}
		}

		if !matchFound {
			continue
		}

		company := ""
		if c, ok := item["company"]; ok {
			company = fmt.Sprintf("%v", c)
		}

		location := "Remote"
		if l, ok := item["location"]; ok && l != nil {
			location = fmt.Sprintf("%v", l)
		}

		description := ""
		if d, ok := item["description"]; ok {
			description = cleanDescription(fmt.Sprintf("%v", d), 200)
		}

		salary := ""
		if s, ok := item["salary_range"]; ok && s != nil {
			salary = fmt.Sprintf("%v", s)
		}

		url := ""
		if u, ok := item["url"]; ok {
			urlStr := fmt.Sprintf("%v", u)
			// Check if URL already starts with http (full URL) or just path
			if strings.HasPrefix(urlStr, "http") {
				url = urlStr
			} else {
				url = fmt.Sprintf("https://remoteok.com%v", urlStr)
			}
		}

		date := ""
		if d, ok := item["date"]; ok {
			date = fmt.Sprintf("%v", d)
		}

		jobs = append(jobs, Job{
			Title:       position,
			Company:     company,
			Location:    location,
			Description: description,
			Salary:      salary,
			JobUrl:      url,
			PostedDate:  date,
			JobType:     "Remote",
		})
	}

	fmt.Printf("  ✓ RemoteOK API: %d jobs\n", len(jobs))
	return jobs, nil
} // generateSampleJobs creates sample job listings based on skills
func generateSampleJobs(skills []string, limit int) []Job {
	fmt.Println("📝 Generating sample job recommendations based on skills")

	// Use top skills to generate relevant job titles
	topSkills := skills
	if len(topSkills) > 3 {
		topSkills = topSkills[:3]
	}

	jobs := []Job{
		{
			Title:       fmt.Sprintf("Senior %s Developer", capitalize(topSkills[0])),
			Company:     "Tech Innovations Inc.",
			Location:    "Remote",
			Description: fmt.Sprintf("Seeking experienced developer proficient in %s. Work on cutting-edge projects with modern technologies.", strings.Join(topSkills, ", ")),
			Salary:      "$120,000 - $160,000",
			JobUrl:      "https://www.linkedin.com/jobs/",
			PostedDate:  time.Now().AddDate(0, 0, -2).Format("2006-01-02"),
			JobType:     "Full-time",
		},
		{
			Title:       fmt.Sprintf("%s Software Engineer", capitalize(topSkills[0])),
			Company:     "Global Solutions Ltd.",
			Location:    "New York, NY",
			Description: fmt.Sprintf("Join our team working with %s and modern frameworks. Competitive benefits and growth opportunities.", topSkills[0]),
			Salary:      "$100,000 - $140,000",
			JobUrl:      "https://www.indeed.com/",
			PostedDate:  time.Now().AddDate(0, 0, -5).Format("2006-01-02"),
			JobType:     "Full-time",
		},
		{
			Title:       "Full Stack Developer",
			Company:     "StartUp Ventures",
			Location:    "San Francisco, CA",
			Description: "Build scalable applications using modern tech stack. Experience with our key technologies is a plus.",
			Salary:      "$110,000 - $150,000",
			JobUrl:      "https://www.glassdoor.com/Job/",
			PostedDate:  time.Now().AddDate(0, 0, -7).Format("2006-01-02"),
			JobType:     "Full-time",
		},
		{
			Title:       fmt.Sprintf("Mid-Level %s Developer", capitalize(topSkills[0])),
			Company:     "Enterprise Corp",
			Location:    "Austin, TX",
			Description: "Growing team seeking talented developers. Work on enterprise-level applications.",
			Salary:      "$90,000 - $120,000",
			JobUrl:      "https://www.monster.com/jobs/",
			PostedDate:  time.Now().AddDate(0, 0, -10).Format("2006-01-02"),
			JobType:     "Full-time",
		},
		{
			Title:       "Software Development Engineer",
			Company:     "Cloud Services Inc.",
			Location:    "Seattle, WA",
			Description: "Build and maintain cloud-based solutions. Strong technical skills required.",
			Salary:      "$115,000 - $145,000",
			JobUrl:      "https://www.dice.com/jobs/",
			PostedDate:  time.Now().AddDate(0, 0, -3).Format("2006-01-02"),
			JobType:     "Full-time",
		},
	}

	if limit < len(jobs) {
		return jobs[:limit]
	}
	return jobs
}

// Helper functions
func cleanDescription(desc string, maxLength int) string {
	// Remove all HTML tags using regex-like approach
	// Simple state machine to remove everything between < and >
	var result strings.Builder
	inTag := false

	for i := 0; i < len(desc); i++ {
		char := desc[i]

		if char == '<' {
			inTag = true
			continue
		}

		if char == '>' {
			inTag = false
			continue
		}

		if !inTag {
			result.WriteByte(char)
		}
	}

	desc = result.String()

	// Replace common line breaks and special characters
	desc = strings.ReplaceAll(desc, "\n", " ")
	desc = strings.ReplaceAll(desc, "\r", " ")
	desc = strings.ReplaceAll(desc, "\t", " ")

	// Decode common HTML entities
	desc = strings.ReplaceAll(desc, "&nbsp;", " ")
	desc = strings.ReplaceAll(desc, "&amp;", "&")
	desc = strings.ReplaceAll(desc, "&lt;", "<")
	desc = strings.ReplaceAll(desc, "&gt;", ">")
	desc = strings.ReplaceAll(desc, "&quot;", "\"")
	desc = strings.ReplaceAll(desc, "&#39;", "'")
	desc = strings.ReplaceAll(desc, "&apos;", "'")

	// Remove excessive whitespace
	desc = strings.Join(strings.Fields(desc), " ")

	// Truncate if too long
	if len(desc) > maxLength {
		desc = desc[:maxLength] + "..."
	}

	return desc
}

func capitalize(s string) string {
	if len(s) == 0 {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func contains(slice []string, item string) bool {
	item = strings.ToLower(item)
	for _, s := range slice {
		if strings.ToLower(s) == item {
			return true
		}
	}
	return false
}
