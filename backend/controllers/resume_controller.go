package controllers

import (
	"backend/config"
	"backend/models"
	"backend/services"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

func UploadResume(c *gin.Context) {
	// Extract authenticated user ID from context (set by AuthMiddleware)
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	// The middleware stores user_id as uint; validate the type
	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	title := c.PostForm("title")
	jobDescription := c.PostForm("job_description") // Optional job description for better ATS matching
	file, err := c.FormFile("resume")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "resume file required"})
		return
	}

	// Validate PDF file
	if file.Header.Get("Content-Type") != "application/pdf" {
		fmt.Printf("⚠️  Warning: Content-Type is %s, expected application/pdf\n", file.Header.Get("Content-Type"))
	}

	// save temporarily with unique name to avoid conflicts
	tempPath := fmt.Sprintf("./temp_resume_%d_%s", time.Now().UnixNano(), file.Filename)
	if err := c.SaveUploadedFile(file, tempPath); err != nil {
		fmt.Println("File save error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save file"})
		return
	}
	defer os.Remove(tempPath) // Clean up temp file after processing

	fmt.Printf("📁 Saved temp file: %s (size: %d bytes)\n", tempPath, file.Size)

	// Extract text from PDF BEFORE uploading
	pdfText, err := services.ExtractTextFromPdfFile(tempPath)
	if err != nil {
		fmt.Println("PDF Extraction Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to extract text from PDF"})
		return
	}
	fmt.Println("Extracted text length:", len(pdfText))

	// Analyze the extracted text with optional job description
	analysis, err := services.AnalyzeResumeText(pdfText, jobDescription)
	if err != nil {
		fmt.Println("AI Analysis Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to analyze resume with AI"})
		return
	}
	fmt.Println("Analysis completed")

	// upload to Appwrite (new storage service)
	url, err := services.UploadResume(tempPath)
	if err != nil {
		fmt.Println("Upload Error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to upload to Appwrite"})
		return
	}
	fmt.Println("Uploaded to:", url)

	// store in database
	resume := models.Resume{
		UserId:         uid,
		Title:          title,
		FileUrl:        url,
		AnalysisResult: "{}",
		AtsScore:       0,
		JdMatchScore:   0,
		MatchingSkills: "[]",
		MissingSkills:  "[]",
		UploadedAt:     time.Now(),
	}
	if err := config.DB.Create(&resume).Error; err != nil {
		fmt.Println("Db error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to save resume"})
		return
	}

	// Save full analysis JSON and extract fields to their own columns
	resume.AnalysisResult = analysis
	var parsed map[string]interface{}
	parseErr := json.Unmarshal([]byte(analysis), &parsed)
	if parseErr != nil {
		// non-fatal: keep default values if parsing fails
		fmt.Println("⚠️  Warning: failed to parse analysis JSON:", parseErr)
	} else {
		fmt.Println("✅ Successfully parsed analysis JSON")

		// Extract ats_score
		if v, ok := parsed["ats_score"]; ok && v != nil {
			switch t := v.(type) {
			case float64:
				resume.AtsScore = int(t)
			case int:
				resume.AtsScore = t
			}
		}

		// Extract jd_match_score
		if v, ok := parsed["jd_match_score"]; ok && v != nil {
			switch t := v.(type) {
			case float64:
				resume.JdMatchScore = int(t)
			case int:
				resume.JdMatchScore = t
			}
		}

		// Extract matching_skills (array)
		if v, ok := parsed["matching_skills"]; ok && v != nil {
			if skillsBytes, err := json.Marshal(v); err == nil {
				resume.MatchingSkills = string(skillsBytes)
			}
		}

		// Extract missing_skills (array)
		if v, ok := parsed["missing_skills"]; ok && v != nil {
			if skillsBytes, err := json.Marshal(v); err == nil {
				resume.MissingSkills = string(skillsBytes)
			}
		}
	}

	if err := config.DB.Save(&resume).Error; err != nil {
		fmt.Println("Db save error:", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update resume analysis"})
		return
	}

	// Fetch job recommendations based on extracted skills
	var recommendedJobs []services.Job
	// Extract skills array from analysis (use parsed map if available)
	var skills []string

	if parseErr == nil && parsed != nil {
		fmt.Println("🔍 Looking for skills in parsed JSON...")
		if skillsInterface, ok := parsed["skills"]; ok && skillsInterface != nil {
			fmt.Printf("✅ Found skills field: %v\n", skillsInterface)
			if skillsArray, ok := skillsInterface.([]interface{}); ok {
				for _, skill := range skillsArray {
					if skillStr, ok := skill.(string); ok {
						skills = append(skills, skillStr)
					}
				}
				fmt.Printf("✅ Extracted %d skills: %v\n", len(skills), skills)
			} else {
				fmt.Println("⚠️  Skills field is not an array")
			}
		} else {
			fmt.Println("⚠️  No skills field found in parsed JSON")
		}
	} else {
		fmt.Println("⚠️  Parsed JSON is nil or parse error occurred")
	}

	// Fetch 5-10 jobs based on skills
	if len(skills) > 0 {
		fmt.Printf("🔍 Fetching job recommendations for %d skills\n", len(skills))
		jobs, err := services.FetchJobRecommendations(skills, 8)
		if err != nil {
			fmt.Println("⚠️  Job fetch error (non-fatal):", err)
			// Don't fail the entire upload if job fetch fails
		} else {
			recommendedJobs = jobs
			fmt.Printf("✅ Fetched %d job recommendations\n", len(recommendedJobs))

			// Save job recommendations to database
			if len(recommendedJobs) > 0 {
				fmt.Println("💾 Saving job recommendations to database...")
				for _, job := range recommendedJobs {
					jobRec := models.JobRecommendation{
						ResumeId:    resume.Id,
						Title:       job.Title,
						Company:     job.Company,
						Location:    job.Location,
						Description: job.Description,
						Salary:      job.Salary,
						JobUrl:      job.JobUrl,
						PostedDate:  job.PostedDate,
						JobType:     job.JobType,
					}
					if err := config.DB.Create(&jobRec).Error; err != nil {
						fmt.Printf("⚠️  Failed to save job recommendation: %v\n", err)
						// Continue saving other jobs even if one fails
					}
				}
				fmt.Printf("✅ Saved %d job recommendations to database\n", len(recommendedJobs))
			}
		}
	} else {
		fmt.Println("⚠️  No skills extracted, skipping job recommendations")
	}

	c.JSON(http.StatusOK, gin.H{
		"message":          "Resume uploaded successfully",
		"file_url":         url,
		"analysis_result":  resume.AnalysisResult,
		"ats_score":        resume.AtsScore,
		"jd_match_score":   resume.JdMatchScore,
		"matching_skills":  resume.MatchingSkills,
		"missing_skills":   resume.MissingSkills,
		"recommended_jobs": recommendedJobs,
	})
}

// GetUserResumes fetches all resumes for the authenticated user
func GetUserResumes(c *gin.Context) {
	// Extract authenticated user ID from context
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	var resumes []models.Resume
	if err := config.DB.Where("user_id = ?", uid).Order("uploaded_at DESC").Find(&resumes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch resumes"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resumes": resumes,
	})
}

// GetResumeById fetches a single resume by ID (only if it belongs to the user)
func GetResumeById(c *gin.Context) {
	// Extract authenticated user ID from context
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	resumeId := c.Param("id")
	var resume models.Resume
	if err := config.DB.Where("id = ? AND user_id = ?", resumeId, uid).First(&resume).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"resume": resume,
	})
}

// DeleteResume deletes a resume (only if it belongs to the user)
func DeleteResume(c *gin.Context) {
	// Extract authenticated user ID from context
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	resumeId := c.Param("id")
	var resume models.Resume
	if err := config.DB.Where("id = ? AND user_id = ?", resumeId, uid).First(&resume).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
		return
	}

	// Delete associated job recommendations first
	config.DB.Where("resume_id = ?", resumeId).Delete(&models.JobRecommendation{})

	// Delete the resume
	if err := config.DB.Delete(&resume).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete resume"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "resume deleted successfully",
	})
}

// GetResumeJobs fetches all job recommendations for a specific resume
func GetResumeJobs(c *gin.Context) {
	// Extract authenticated user ID from context
	uidVal, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return
	}

	uid, ok := uidVal.(uint)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid user context"})
		return
	}

	resumeId := c.Param("id")

	// Verify the resume belongs to the user
	var resume models.Resume
	if err := config.DB.Where("id = ? AND user_id = ?", resumeId, uid).First(&resume).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "resume not found"})
		return
	}

	// Fetch job recommendations
	var jobs []models.JobRecommendation
	if err := config.DB.Where("resume_id = ?", resumeId).Order("created_at DESC").Find(&jobs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch job recommendations"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"jobs": jobs,
	})
}
