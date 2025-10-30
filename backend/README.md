# Smart Resume Backend

A powerful Go-based backend system for AI-powered resume analysis, ATS scoring, job description matching, and intelligent job recommendations.

## 🚀 Features

### Core Functionality
- **Resume Upload & Storage**: Upload PDF resumes to Appwrite cloud storage
- **PDF Text Extraction**: Extract clean text from PDF resumes using `ledongthuc/pdf` library
- **AI-Powered Analysis**: Analyze resumes using FastAPI + spaCy NLP for entity and skill extraction
- **User Authentication**: JWT-based authentication for secure access
- **PostgreSQL Database**: Store resumes, users, and job recommendations with GORM ORM

### ATS (Applicant Tracking System) Scoring
- **Comprehensive Scoring (0-100)**: Evaluate resumes based on:
  - Skills Match (30-40 points)
  - Experience/Education Detection (25-30 points)
  - Resume Structure (15-30 points)
  - Job Description Match (30 points when JD provided)
- **Adaptive Scoring**: Automatically adjusts scoring weights based on whether a job description is provided
- **40+ Technical Skills Detection**: Recognizes skills like Python, JavaScript, React, Docker, AWS, etc.

### Job Description Matching
- **Optional JD Parameter**: Provide job description for better ATS scoring
- **NLP-Based Keyword Extraction**: Uses spaCy for intelligent keyword analysis
- **Skill Gap Analysis**: 
  - Returns matching skills from resume
  - Identifies missing skills required for the job
  - Calculates JD match score (0-100)

### Intelligent Job Recommendations - Parallel API System 🚀

#### Key Features
- **Parallel API Architecture**: 3-4 APIs run simultaneously (2-3x faster than sequential)
- **2-Tier Priority System**: Smart fallback ensures 100% reliability
- **Automatic Deduplication**: Removes duplicate jobs by title + company
- **Skill-Based Filtering**: Matches jobs based on extracted resume skills
- **Database Storage**: Saves job recommendations linked to resumes
- **Zero Dependencies**: Works perfectly without any API keys!

#### 🔵 Priority 1 (Run in Parallel)
- **RemoteOK** (Free, no auth) - Tech jobs, remote positions ⭐⭐⭐⭐⭐
- **Arbeitnow** (Free, no auth) - European + US jobs ⭐⭐⭐⭐
- **The Muse** (Free, no auth) - Curated quality jobs ⭐⭐⭐⭐
- **Adzuna** (Optional) - Large database if credentials provided ⭐⭐⭐⭐⭐

#### 🟡 Priority 2 (Backup, Parallel)
- **Findwork** (Free, test token) - Tech focus ⭐⭐⭐
- **Jooble** (Optional) - Worldwide if API key provided ⭐⭐⭐⭐
- **JSearch/RapidAPI** (Optional) - Google Jobs if key provided ⭐⭐⭐⭐⭐

#### Performance Benefits
- **Speed**: 2-3x faster than sequential approach
- **Diversity**: Jobs from 3-4 different sources per request
- **Reliability**: If Priority 1 fails, Priority 2 kicks in automatically
- **Quality**: Deduplication ensures unique, high-quality results

#### How It Works
```
1. Extract skills from resume
2. Launch Priority 1 APIs in parallel (3-4 simultaneously)
3. Collect results (5-10 seconds)
4. If enough jobs → Deduplicate → Return
5. If insufficient → Launch Priority 2 in parallel
6. Combine all results → Deduplicate → Return
7. If all fail → Sample jobs (final fallback)
```

## 📁 Project Structure

```
backend/
├── analyzer/               # FastAPI NLP service
│   ├── app.py             # Main analyzer service (port 8000)
│   └── requirements.txt   # Python dependencies
├── config/
│   └── config.go          # Database & environment configuration
├── controllers/
│   ├── auth_controller.go      # User registration & login
│   ├── resume_controller.go    # Resume upload & analysis
│   └── user_controller.go      # User profile management
├── middlewares/
│   └── auth_middleware.go      # JWT authentication middleware
├── models/
│   ├── user.go                 # User model
│   ├── resume.go               # Resume model with ATS fields
│   └── jobRecommendation.go    # Job recommendation model
├── routes/
│   └── routes.go               # API route definitions
├── services/
│   ├── analyzer.go             # AI analysis & PDF extraction
│   ├── job_fetcher.go          # Job API integrations
│   └── storage.go              # Appwrite storage service
├── utils/
│   └── token.go                # JWT token generation & validation
├── .env                        # Environment variables
├── go.mod                      # Go dependencies
├── go.sum                      # Dependency checksums
└── main.go                     # Application entry point
```

## 🛠️ Tech Stack

- **Language**: Go 1.x
- **Framework**: Gin (HTTP web framework)
- **Database**: PostgreSQL with GORM ORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Storage**: Appwrite Cloud Storage
- **PDF Processing**: ledongthuc/pdf
- **AI/NLP**: FastAPI + spaCy (Python analyzer service)
- **Environment**: godotenv

## 📦 Installation

### Prerequisites
- Go 1.19 or higher
- PostgreSQL 12 or higher
- Python 3.8+ (for analyzer service)
- Appwrite account (for file storage)

### Step 1: Clone and Setup

```bash
cd backend
go mod download
```

### Step 2: Configure Environment Variables

Create a `.env` file in the backend directory:

```env
# Database Configuration
DB_USER=postgres
DB_PASSWORD=your_password
DB_NAME=smart_resume_db
DB_HOST=localhost
DB_PORT=5432

# JWT Secret (generate a secure random string)
JWT_SECRET=your_jwt_secret_key

# Appwrite Configuration
APPWRITE_ENDPOINT=https://cloud.appwrite.io/v1
APPWRITE_PROJECT_ID=your_project_id
APPWRITE_API_KEY=your_api_key
APPWRITE_BUCKET_ID=your_bucket_id

# Optional: AI API Keys
GEMINI_API_KEY=your_gemini_api_key

### Job API Credentials (Optional - Enhances Results)

**Without any API keys**: System works perfectly using RemoteOK, Arbeitnow, The Muse, and Findwork (all free)

**With credentials** (better coverage):
```env
# Adzuna (Priority 1 - recommended)
ADZUNA_APP_ID=your_adzuna_app_id
ADZUNA_APP_KEY=your_adzuna_app_key
JOB_COUNTRY=us

# Jooble (Priority 2)
JOOBLE_API_KEY=your_jooble_api_key

# RapidAPI JSearch (Priority 2)
RAPIDAPI_KEY=your_rapidapi_key
```
```

### Step 3: Setup PostgreSQL Database

```bash
# Create database
psql -U postgres
CREATE DATABASE smart_resume_db;
\q
```

The application will auto-migrate tables on first run.

### Step 4: Setup Python Analyzer Service

```bash
cd analyzer
pip install -r requirements.txt
python -m spacy download en_core_web_sm
```

### Step 5: Build and Run

**Terminal 1 - Start Analyzer Service:**
```bash
cd analyzer
python -m uvicorn app:app --host 0.0.0.0 --port 8000
```

**Terminal 2 - Start Go Backend:**
```bash
cd backend
go build -o backend.exe
./backend.exe
```

Backend will run on `http://localhost:8080`

## 🔌 API Endpoints

### Authentication
- `POST /api/auth/register` - Register new user
- `POST /api/auth/login` - Login user

### Resume Management (Protected)
- `POST /api/resume/upload` - Upload and analyze resume
  - **Form Data**:
    - `title`: Resume title
    - `resume`: PDF file
    - `job_description` (optional): Job description for better matching

### User Profile (Protected)
- `GET /api/user/profile` - Get user profile
- `GET /api/user/resumes` - Get all user resumes

## 📊 Database Schema

### Users Table
- `id` (primary key)
- `name`
- `email` (unique)
- `password` (bcrypt hashed)
- `created_at`

### Resumes Table
- `id` (primary key)
- `user_id` (foreign key)
- `title`
- `file_url` (Appwrite storage URL)
- `analysis_result` (JSONB)
- `ats_score` (integer, 0-100)
- `jd_match_score` (integer, 0-100)
- `matching_skills` (JSONB array)
- `missing_skills` (JSONB array)
- `uploaded_at`

### Job Recommendations Table
- `id` (primary key)
- `resume_id` (foreign key)
- `title`
- `company`
- `location`
- `description`
- `salary`
- `job_url`
- `posted_date`
- `job_type`
- `created_at`

## 📤 API Response Example

```json
{
  "message": "Resume uploaded successfully",
  "file_url": "https://cloud.appwrite.io/v1/storage/buckets/.../view",
  "analysis_result": "{...}",
  "ats_score": 85,
  "jd_match_score": 72,
  "matching_skills": "[\"python\", \"react\", \"docker\"]",
  "missing_skills": "[\"kubernetes\", \"aws\"]",
  "recommended_jobs": [
    {
      "title": "Senior Python Developer",
      "company": "Tech Corp",
      "location": "Remote",
      "description": "Looking for experienced Python developer...",
      "salary": "$120,000 - $160,000",
      "job_url": "https://...",
      "posted_date": "2025-10-28",
      "job_type": "Full-time"
    }
  ]
}
```

## 🔧 Configuration Notes

### Appwrite Setup
1. Create account at [Appwrite Cloud](https://cloud.appwrite.io)
2. Create a new project
3. Create a storage bucket with **public read** permissions
4. Get your Project ID, API Key, and Bucket ID

### Job API Setup

The parallel job system works out of the box with free APIs! Optional credentials improve coverage:

#### Free APIs (No Setup Required)
- **RemoteOK**: Automatic, no auth needed
- **Arbeitnow**: Automatic, no auth needed  
- **The Muse**: Automatic, no auth needed
- **Findwork**: Automatic, uses test token

#### Optional APIs (Better Results)
- **Adzuna**: Sign up at [developer.adzuna.com](https://developer.adzuna.com) - Large job database
- **Jooble**: Get API key from [Jooble API](https://jooble.org/api/about) - Worldwide coverage
- **RapidAPI**: Subscribe to JSearch at [RapidAPI](https://rapidapi.com) - Google Jobs access

### Parallel System Architecture

```
Resume Upload → Extract Skills
        ↓
┌───────────────────────────────────┐
│  PRIORITY 1 (Parallel Execution)  │
├───────────────────────────────────┤
│  • RemoteOK    (Free) ✓           │
│  • Arbeitnow   (Free) ✓           │
│  • The Muse    (Free) ✓           │
│  • Adzuna      (If creds)         │
└───────────────────────────────────┘
        ↓
   Enough jobs?
   ↙Yes    No↘
   ↓           ↓
   ↓    ┌─────────────────────────┐
   ↓    │ PRIORITY 2 (Parallel)   │
   ↓    ├─────────────────────────┤
   ↓    │ • Findwork  (Free) ✓    │
   ↓    │ • Jooble    (If key)    │
   ↓    │ • JSearch   (If key)    │
   ↓    └─────────────────────────┘
   ↓              ↓
   └──────────────┘
          ↓
   Deduplication
    (Title+Company)
          ↓
   Return Unique Jobs
```

### Performance Metrics

| Scenario | Response Time | Sources | Unique Jobs |
|----------|--------------|---------|-------------|
| Priority 1 Success | ~5s | 3-4 | 20+ |
| With Priority 2 | ~12s | 5-7 | 30+ |
| All APIs Fail | ~2s | Sample | 8 |

## 🐛 Troubleshooting

### Analyzer Service Not Connecting
- Ensure analyzer is running on port 8000
- Check `http://localhost:8000/health` endpoint
- On Windows, run analyzer in separate CMD window (not PowerShell background)

### PDF Extraction Issues
- Ensure PDF is not password-protected
- Check if PDF contains selectable text (not scanned image)
- Temp files are automatically cleaned up after processing

### Database Connection Failed
- Verify PostgreSQL is running: `pg_isready`
- Check credentials in `.env` file
- Ensure database exists and is accessible

### Job Recommendations Empty or Slow
**Issue**: No jobs returned or taking too long
```
Solution:
1. Check if skills were extracted from resume (backend logs)
2. Verify analyzer is running and extracting skills correctly
3. Check backend logs to see which APIs are being tried:
   - "🚀 Starting Parallel Job Fetch System"
   - "✓ RemoteOK API: X jobs" (should see multiple)
   - "✅ SUCCESS: Got X jobs from Priority 1"
4. Optional: Add Adzuna credentials for better results
5. Sample jobs will be returned if all APIs fail (never truly fails)

**Performance**:
- With free APIs: ~5-10 seconds for 20+ jobs
- With Priority 2 backup: ~10-15 seconds for 30+ jobs
- Watch backend logs for parallel execution details
```

### Backend Logs Show API Details
When you upload a resume, you'll see real-time parallel execution:
```
🚀 Starting Parallel Job Fetch System
📊 Target: 8 jobs from skills: [javascript, react, node]

🔵 Priority 1: Fetching from 3-4 APIs simultaneously...
  ✓ RemoteOK API: 5 jobs
  ✓ Arbeitnow API: 4 jobs  
  ✓ The Muse API: 3 jobs
  ✗ Adzuna API: credentials not available
📈 Parallel fetch complete: 3 succeeded, 1 failed

✅ SUCCESS: Got 12 jobs from Priority 1 APIs
🔧 Deduplication: 12 jobs → 8 unique jobs
```

## 🚦 Health Check

Backend health: `GET http://localhost:8080/`
Analyzer health: `GET http://localhost:8000/health`

## 📝 License

MIT License

## 👥 Contributors

- Syed Sayan

## 🔮 Future Enhancements

- [ ] Resume templates and suggestions
- [ ] Multi-language support
- [ ] Resume comparison feature
- [ ] Interview preparation tips based on skills
- [ ] Cover letter generation
- [ ] LinkedIn profile optimization
