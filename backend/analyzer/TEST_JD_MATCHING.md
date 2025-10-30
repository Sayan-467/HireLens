# Enhanced ATS Scoring with Job Description Matching

## Overview
The ATS scoring system now includes **Job Description (JD) matching** to provide more realistic and reliable scores based on how well a resume matches a specific job posting.

## Scoring System (Total: 100 points)

### With Job Description Provided:
1. **Skills Match** (30 points) - Technical skills found in resume
2. **Experience & Education** (25 points) - Work history, education, certifications
3. **Resume Structure** (15 points) - Presence of standard resume sections
4. **Job Description Match** (30 points) - NEW! How well resume matches JD

### Without Job Description (Backward Compatible):
1. **Skills Match** (40 points) - Increased weight
2. **Experience & Education** (30 points) - Increased weight
3. **Resume Structure** (30 points) - Increased weight
4. **Job Description Match** (0 points) - Not applicable

## Job Description Matching Logic

The JD matching score (max 30 points) is calculated based on:

### 1. Skill Matching (15 points)
- Extracts technical skills from both resume and JD
- Calculates overlap percentage
- Score = (Matched Skills / Required Skills) × 15

**Example:**
- JD requires: Python, React, Docker, AWS, MongoDB (5 skills)
- Resume has: Python, React, Docker, AWS (4 skills)
- Match ratio: 4/5 = 80%
- Points: 0.80 × 15 = 12 points

### 2. Keyword Matching (15 points)
- Uses NLP to extract important keywords (nouns, named entities)
- Compares resume keywords with JD keywords
- Score = (Matched Keywords / JD Keywords) × 15

**Example:**
- JD keywords: engineer, software, senior, cloud, microservices, api, development
- Resume keywords: engineer, software, senior, cloud, api, full-stack, development
- Match ratio: 6/7 = 85.7%
- Points: 0.857 × 15 = 12.8 points (rounded to 12)

## API Request Format

### Without Job Description (Old Behavior)
```json
POST http://localhost:8000/analyze
{
  "text": "Resume text here..."
}
```

### With Job Description (New Feature)
```json
POST http://localhost:8000/analyze
{
  "text": "Resume text here...",
  "job_description": "Job posting text here..."
}
```

## Example Test Results

### Test 1: Without JD
**Resume:** Software engineer with Python, Java, React, Node, Docker, AWS, Git
**Score:** 78/100
**Breakdown:**
- Skills (40 pts): 7 skills found → 35 points
- Experience/Education (30 pts): Has experience + education → 24 points
- Structure (30 pts): 4 sections found → 25 points
- JD Match (0 pts): Not applicable

### Test 2: With Matching JD
**Resume:** Same as above
**JD:** "Looking for Senior Software Engineer with Python, React, Docker, AWS experience. Must have Computer Science degree and worked on e-commerce projects."
**Score:** 84/100
**Breakdown:**
- Skills (30 pts): 7 skills found → 30 points
- Experience/Education (25 pts): Has experience + education → 20 points
- Structure (15 pts): 4 sections found → 12 points
- JD Match (30 pts): High overlap → 22 points

### Test 3: With Non-Matching JD
**Resume:** Same as above
**JD:** "Looking for Mobile Developer with Swift, Kotlin, Flutter, React Native. iOS and Android development required."
**Score:** 58/100
**Breakdown:**
- Skills (30 pts): 7 skills found → 30 points
- Experience/Education (25 pts): Has experience + education → 20 points
- Structure (15 pts): 4 sections found → 12 points
- JD Match (30 pts): Very low overlap → 3 points

## Using in Your Backend

### Postman Request
```
POST http://localhost:8080/api/resume/upload
Content-Type: multipart/form-data

Fields:
- title: "My Resume"
- resume: [PDF file]
- job_description: "Looking for Python developer with React and AWS experience..." (OPTIONAL)
```

### Response
```json
{
    "message": "Resume uploaded successfully",
    "file_url": "https://fra.cloud.appwrite.io/...",
    "analysis_result": {
        "entities": [...],
        "skills": ["python", "react", "aws", ...],
        "summary": "...",
        "ats_score": 84
    },
    "ats_score": 84
}
```

## Benefits

1. **More Realistic Scoring**: Scores reflect actual job requirements
2. **Backward Compatible**: Works with or without JD
3. **Better Candidate Ranking**: Candidates can be ranked by JD match
4. **Actionable Insights**: Shows what skills/keywords are missing
5. **Fair Comparison**: All scores use the same 0-100 scale

## Tips for Better Scores

### For Candidates:
1. Read the job description carefully
2. Include relevant keywords from the JD
3. Match technical skills mentioned in JD
4. Use similar terminology as the job posting
5. Highlight relevant experience and projects

### For Recruiters:
1. Provide clear, detailed job descriptions
2. List specific technical skills required
3. Include key responsibilities and qualifications
4. Use industry-standard terminology
5. Be specific about experience level needed
