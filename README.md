# 🎯 HireLens - AI-Powered Resume Analysis Platform

> Transform your resume into opportunities with AI-powered ATS scoring, skill analysis, and intelligent job recommendations.

[![Next.js](https://img.shields.io/badge/Next.js-16.0-black)](https://nextjs.org/)
[![Go](https://img.shields.io/badge/Go-1.21-00ADD8)](https://golang.org/)
[![Python](https://img.shields.io/badge/Python-3.10-3776AB)](https://www.python.org/)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-316192)](https://www.postgresql.org/)

---

## 📋 Table of Contents

- [Features](#-features)
- [Tech Stack](#-tech-stack)
- [Architecture](#-architecture)
- [Quick Start](#-quick-start)
- [API Documentation](#-api-documentation)
- [Database Schema](#-database-schema)
- [Job Recommendation System](#-job-recommendation-system)
- [Deployment](#-deployment)
- [Troubleshooting](#-troubleshooting)

---

## ✨ Features

### 🎯 Core Functionality
- **AI-Powered Resume Analysis** - Extract skills, experience, and education using NLP
- **ATS Score Calculation** - Comprehensive 0-100 scoring based on multiple factors
- **Job Description Matching** - Compare resume against job requirements
- **Skill Gap Analysis** - Identify matching and missing skills
- **Intelligent Job Recommendations** - Get relevant jobs from multiple sources
- **Real-Time Processing** - Fast analysis with parallel API architecture
- **Secure Authentication** - JWT-based auth with bcrypt password hashing
- **Cloud Storage** - Resume files stored securely on Appwrite

### 📊 Advanced Features
- **Parallel Job API System** - Fetch jobs from 3-4 sources simultaneously (2-3x faster)
- **Automatic Deduplication** - Remove duplicate job listings
- **7-Tier API Fallback** - Multiple job sources ensure 100% availability
- **Dark Mode UI** - Modern, eye-friendly interface
- **Responsive Design** - Works on desktop, tablet, and mobile
- **Animated Gauges** - Visual ATS score representation

---

## 🛠️ Tech Stack

### Frontend
- **Framework**: Next.js 16 (React 19, App Router)
- **Language**: TypeScript
- **Styling**: Tailwind CSS v4
- **State Management**: React Hooks
- **HTTP Client**: Fetch API

### Backend
- **Language**: Go 1.21+
- **Framework**: Gin Web Framework
- **Database**: PostgreSQL 15 with GORM
- **Authentication**: JWT (golang-jwt/jwt)
- **Password Hashing**: bcrypt
- **PDF Processing**: ledongthuc/pdf

### AI Analyzer
- **Language**: Python 3.10+
- **Framework**: FastAPI
- **NLP**: spaCy (en_core_web_sm)
- **PDF Extraction**: pdfplumber
- **Server**: Uvicorn (ASGI)

### External Services
- **File Storage**: Appwrite Cloud
- **Job APIs**: RemoteOK, Arbeitnow, The Muse, Adzuna, Jooble, JSearch, Findwork

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────────┐
│                    FRONTEND (Next.js - Port 3000)                    │
│  ┌────────────┐  ┌────────────┐  ┌────────────┐  ┌────────────┐   │
│  │  Landing   │  │  Register  │  │   Login    │  │ Dashboard  │   │
│  │    Page    │  │    Page    │  │    Page    │  │   (Auth)   │   │
│  └────────────┘  └────────────┘  └────────────┘  └────────────┘   │
└─────────────────────────┬────────────────────────────────────────────┘
                          │ REST API (JWT Auth)
                          │
┌─────────────────────────▼────────────────────────────────────────────┐
│                    BACKEND (Go/Gin - Port 8080)                       │
│  ┌──────────────────────────────────────────────────────────────┐   │
│  │  Routes:                                                       │   │
│  │  • POST /api/signup    • POST /api/login                     │   │
│  │  • GET  /api/profile   • POST /api/resume/upload (Protected) │   │
│  │  • GET  /api/resumes   • GET  /api/resume/:id                │   │
│  │  • DELETE /api/resume/:id   • GET /api/resume/:id/jobs       │   │
│  └──────────────────────────────────────────────────────────────┘   │
│                                                                       │
│  Services: PDF Extract | AI Analysis | Storage | Job Fetch           │
└──┬──────────────────┬──────────────────┬──────────────────┬─────────┘
   │                  │                  │                  │
   │                  │                  │                  │
   ▼                  ▼                  ▼                  ▼
┌─────────┐  ┌─────────────────┐  ┌──────────────┐  ┌──────────────┐
│ PostgreSQL│  │  AI Analyzer    │  │   Appwrite   │  │  Job APIs    │
│  Database │  │  (FastAPI 8000) │  │File Storage  │  │ (7-Tier)     │
│   :5432   │  │  Python+spaCy   │  │   (Cloud)    │  │ Parallel     │
└───────────┘  └─────────────────┘  └──────────────┘  └──────────────┘
```

### 🔄 Complete Data Flow

```
1. USER REGISTRATION/LOGIN
   Frontend → POST /api/signup → Backend
   Backend → Create user (bcrypt hash) → PostgreSQL
   Backend → Generate JWT token → Frontend
   Frontend → Store token → localStorage

2. RESUME UPLOAD & ANALYSIS
   Frontend → POST /api/resume/upload (with JWT) → Backend
   Backend → Save PDF temporarily
   Backend → Extract text (pdfplumber)
   Backend → POST http://localhost:8000/analyze → AI Analyzer
   AI Analyzer → spaCy NLP → Extract skills, calculate scores
   Backend → Upload file → Appwrite Storage
   Backend → Save resume + analysis → PostgreSQL
   Backend → Fetch job recommendations (Parallel APIs)
   Backend → Save jobs → PostgreSQL
   Backend → Return full analysis → Frontend
   Frontend → Display scores, skills, jobs

3. VIEW RESUMES
   Frontend → GET /api/resumes (with JWT) → Backend
   Backend → Query PostgreSQL → Return user's resumes
   Frontend → Display resume list with details
```

---

## 🚀 Quick Start

### Prerequisites

- **Go** 1.21+ - [Download](https://golang.org/dl/)
- **Python** 3.10+ - [Download](https://www.python.org/downloads/)
- **Node.js** 18+ - [Download](https://nodejs.org/)
- **PostgreSQL** 15+ - [Download](https://www.postgresql.org/download/)
- **Git** - [Download](https://git-scm.com/downloads/)

### Step 1: Clone Repository

```bash
git clone https://github.com/Sayan-467/HireLens.git
cd HireLens
```

### Step 2: Database Setup

```bash
# Start PostgreSQL and create database
psql -U postgres
CREATE DATABASE smart_resume_db;
\q
```

### Step 3: Backend Setup

```bash
cd backend

# Install dependencies
go mod download

# Create .env file
cp .env.example .env

# Edit .env with your credentials:
# - Database credentials
# - JWT secret (min 32 characters)
# - Appwrite credentials
# - Optional: Job API keys
```

**Required `.env` variables:**
```env
# Database (PostgreSQL URI)
DATABASE_URL=postgres://postgres:your_password@localhost:5432/smart_resume_db?sslmode=disable

# JWT (generate random 32+ char string)
JWT_SECRET=your_super_secret_jwt_key_min_32_characters

# Appwrite Storage
APPWRITE_ENDPOINT=https://cloud.appwrite.io/v1
APPWRITE_PROJECT_ID=your_project_id
APPWRITE_API_KEY=your_api_key
APPWRITE_BUCKET_ID=your_bucket_id

# Optional: Job APIs
ADZUNA_APP_ID=your_adzuna_app_id
ADZUNA_APP_KEY=your_adzuna_app_key
JOOBLE_API_KEY=your_jooble_key
RAPIDAPI_KEY=your_rapidapi_key
```

### Step 4: AI Analyzer Setup

```bash
cd backend/analyzer

# Create virtual environment
python -m venv venv

# Activate virtual environment
# Windows:
.\venv\Scripts\Activate.ps1
# Linux/Mac:
source venv/bin/activate

# Install dependencies
pip install -r requirements.txt

# Download spaCy model
python -m spacy download en_core_web_sm
```

### Step 5: Frontend Setup

```bash
cd frontend

# Install dependencies
npm install

# Create .env.local file
echo "NEXT_PUBLIC_API_URL=http://localhost:8080" > .env.local
```

### Step 6: Start All Services

**Terminal 1 - Backend (Port 8080):**
```bash
cd backend
go run main.go
```

**Terminal 2 - AI Analyzer (Port 8000):**
```bash
cd backend/analyzer
.\venv\Scripts\Activate.ps1  # Windows
# source venv/bin/activate    # Linux/Mac
python main.py
```

**Terminal 3 - Frontend (Port 3000):**
```bash
cd frontend
npm run dev
```

### Step 7: Access Application

Open browser: **http://localhost:3000**

---

## 📡 API Documentation

### Authentication Endpoints

#### Register User
```http
POST /api/signup
Content-Type: application/json

{
  "name": "John Doe",
  "email": "john@example.com",
  "password": "secure123"
}

Response: 200 OK
{
  "message": "user registered successfully"
}
```

#### Login User
```http
POST /api/login
Content-Type: application/json

{
  "email": "john@example.com",
  "password": "secure123"
}

Response: 200 OK
{
  "message": "login successful",
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

#### Get Profile (Protected)
```http
GET /api/profile
Authorization: Bearer <jwt_token>

Response: 200 OK
{
  "id": 1,
  "name": "John Doe",
  "email": "john@example.com",
  "created_at": "2025-10-31T10:00:00Z"
}
```

### Resume Endpoints

#### Upload Resume (Protected)
```http
POST /api/resume/upload
Authorization: Bearer <jwt_token>
Content-Type: multipart/form-data

Form Data:
- title: "My Software Engineer Resume"
- resume: <pdf_file>
- job_description: "Looking for Python developer..." (optional)

Response: 200 OK
{
  "message": "Resume uploaded successfully",
  "file_url": "https://cloud.appwrite.io/v1/storage/...",
  "analysis_result": "{...full_json...}",
  "ats_score": 85,
  "jd_match_score": 78,
  "matching_skills": "[\"Python\",\"React\",\"Docker\"]",
  "missing_skills": "[\"Kubernetes\",\"AWS\"]",
  "recommended_jobs": [
    {
      "id": 1,
      "title": "Senior Python Developer",
      "company": "Tech Corp",
      "location": "Remote",
      "description": "...",
      "salary": "$120k-$160k",
      "job_url": "https://...",
      "posted_date": "2 days ago",
      "job_type": "Full-time"
    }
  ]
}
```

#### Get All Resumes (Protected)
```http
GET /api/resumes
Authorization: Bearer <jwt_token>

Response: 200 OK
[
  {
    "id": 1,
    "title": "Software Engineer Resume",
    "file_url": "https://...",
    "ats_score": 85,
    "jd_match_score": 78,
    "uploaded_at": "2025-10-31T10:00:00Z"
  }
]
```

#### Get Resume Details (Protected)
```http
GET /api/resume/:id
Authorization: Bearer <jwt_token>

Response: 200 OK
{
  "id": 1,
  "title": "Software Engineer Resume",
  "file_url": "https://...",
  "ats_score": 85,
  "jd_match_score": 78,
  "matching_skills": ["Python", "React"],
  "missing_skills": ["AWS", "Docker"],
  "uploaded_at": "2025-10-31T10:00:00Z"
}
```

#### Delete Resume (Protected)
```http
DELETE /api/resume/:id
Authorization: Bearer <jwt_token>

Response: 200 OK
{
  "message": "Resume deleted successfully"
}
```

#### Get Job Recommendations (Protected)
```http
GET /api/resume/:id/jobs
Authorization: Bearer <jwt_token>

Response: 200 OK
[
  {
    "id": 1,
    "title": "Senior Developer",
    "company": "Tech Corp",
    "location": "Remote",
    "salary": "$120k-$160k",
    "job_url": "https://..."
  }
]
```

---

## 🗄️ Database Schema

### Users Table
```sql
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    password VARCHAR(255) NOT NULL,  -- bcrypt hashed
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Resumes Table
```sql
CREATE TABLE resumes (
    id SERIAL PRIMARY KEY,
    user_id INTEGER REFERENCES users(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    file_url TEXT NOT NULL,
    analysis_result JSONB,
    ats_score INTEGER,              -- 0-100
    jd_match_score INTEGER,          -- 0-100
    matching_skills JSONB,           -- ["Python", "React"]
    missing_skills JSONB,            -- ["AWS", "Docker"]
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Job Recommendations Table
```sql
CREATE TABLE job_recommendations (
    id SERIAL PRIMARY KEY,
    resume_id INTEGER REFERENCES resumes(id) ON DELETE CASCADE,
    title VARCHAR(255) NOT NULL,
    company VARCHAR(255) NOT NULL,
    location VARCHAR(255),
    description TEXT,
    salary VARCHAR(100),
    job_url TEXT,
    posted_date VARCHAR(50),
    job_type VARCHAR(50),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## 🎯 Job Recommendation System

### Parallel API Architecture

The system uses a **2-tier parallel architecture** for maximum speed and reliability:

#### 🔵 Priority 1 (3-4 APIs in Parallel)
- **RemoteOK** - Free, tech jobs, remote positions
- **Arbeitnow** - Free, EU + US jobs
- **The Muse** - Free, curated quality jobs
- **Adzuna** - Optional (with credentials), large database

#### 🟡 Priority 2 (Backup APIs)
- **Findwork** - Free, tech focus
- **Jooble** - Optional (with key), worldwide
- **JSearch/RapidAPI** - Optional (with key), Google Jobs

### How It Works

```
1. Priority 1 APIs launch simultaneously (parallel execution)
2. Collect all results (5-10 seconds)
3. If enough jobs → Deduplicate → Return
4. If insufficient → Launch Priority 2 in parallel
5. Combine all results → Deduplicate → Return
```

### Performance

- **2-3x faster** than sequential approach
- **3-4 different sources** per request
- **Automatic deduplication** (title + company)
- **100% reliability** (sample jobs as final fallback)

### Configuration

```env
# Optional - Better results with API keys
ADZUNA_APP_ID=your_app_id
ADZUNA_APP_KEY=your_app_key
JOOBLE_API_KEY=your_key
RAPIDAPI_KEY=your_key
```

**Works without any API keys!** Free APIs provide excellent coverage.

---

## 🧪 Testing

### Manual Testing Flow

1. **Start all services** (Backend, Analyzer, Frontend)
2. **Register** new account at `/register`
3. **Login** with credentials
4. **Upload Resume**:
   - Add title
   - Select PDF file
   - Optionally add job description
   - Click "Analyze Resume"
5. **View Results**:
   - ATS score gauge (0-100%)
   - JD match score (if JD provided)
   - Matching skills (green badges)
   - Missing skills (red badges)
   - Job recommendations (cards)

### API Testing with curl

```bash
# Register
curl -X POST http://localhost:8080/api/signup \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com","password":"test123"}'

# Login
curl -X POST http://localhost:8080/api/login \
  -H "Content-Type: application/json" \
  -d '{"email":"test@example.com","password":"test123"}'

# Get Profile (use token from login)
curl -X GET http://localhost:8080/api/profile \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## 📦 Deployment

### Frontend (Vercel)

```bash
cd frontend
vercel deploy --prod
```

Set environment variable:
- `NEXT_PUBLIC_API_URL` = Your backend URL

### Backend (Railway/Render)

1. Push code to GitHub
2. Connect repository to Railway/Render
3. Set environment variables (.env)
4. Deploy

### Analyzer (Python Hosting)

- Deploy to Heroku, Railway, or DigitalOcean
- Ensure port 8000 is exposed
- Install Python dependencies
- Download spaCy model during build

### Database (Cloud PostgreSQL)

- Use Supabase, Railway, or AWS RDS
- Update `DATABASE_URL` in backend .env with your cloud database URI
- Run migrations (auto-handled by GORM)

---

## 🐛 Troubleshooting

### Backend Won't Start

**Issue**: Database connection failed
```
Solution:
1. Check PostgreSQL is running: pg_isready
2. Verify credentials in .env
3. Ensure database exists: psql -l
4. Check port 5432 is not blocked
```

**Issue**: Port 8080 already in use
```
Solution (Windows):
netstat -ano | findstr :8080
taskkill /PID <PID> /F
```

### Analyzer Won't Start

**Issue**: spaCy model not found
```
Solution:
python -m spacy download en_core_web_sm
```

**Issue**: Import errors
```
Solution:
pip install --upgrade -r requirements.txt
```

### Frontend Issues

**Issue**: API connection failed / CORS errors
```
Solution:
1. Check backend running on :8080
2. Verify NEXT_PUBLIC_API_URL in .env.local
3. Check backend CORS config allows localhost:3000
```

**Issue**: Build fails
```
Solution:
rm -rf .next node_modules
npm install
npm run dev
```

### Upload Issues

**Issue**: Resume upload fails
```
Solution:
1. Check Appwrite credentials in backend .env
2. Verify bucket exists and has public read permissions
3. Check file size < 10MB
4. Ensure file is PDF format
```

**Issue**: No job recommendations
```
Solution:
1. Check analyzer extracted skills correctly
2. Backend logs show which APIs are tried
3. Verify network connectivity
4. Sample jobs will be returned if all APIs fail
```

---

## 📝 Environment Variables Summary

### Backend `.env`
```env
# Database (PostgreSQL URI)
DATABASE_URL=postgres://postgres:your_password@localhost:5432/smart_resume_db?sslmode=disable

# JWT
JWT_SECRET=your_very_secret_key_min_32_chars

# Appwrite
APPWRITE_ENDPOINT=https://cloud.appwrite.io/v1
APPWRITE_PROJECT_ID=your_project_id
APPWRITE_API_KEY=your_api_key
APPWRITE_BUCKET_ID=your_bucket_id

# Optional Job APIs
ADZUNA_APP_ID=your_app_id
ADZUNA_APP_KEY=your_app_key
JOOBLE_API_KEY=your_key
RAPIDAPI_KEY=your_key
JOB_COUNTRY=us
```

### Frontend `.env.local`
```env
NEXT_PUBLIC_API_URL=http://localhost:8080
```

---

## 🤝 Contributing

Contributions are welcome! Please follow these steps:

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/AmazingFeature`)
3. Commit your changes (`git commit -m 'Add some AmazingFeature'`)
4. Push to the branch (`git push origin feature/AmazingFeature`)
5. Open a Pull Request

---

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

---

## 👥 Authors

- **Syed Sayan** - [GitHub](https://github.com/Sayan-467)

---

## 🙏 Acknowledgments

- Next.js team for the amazing framework
- spaCy for the NLP library
- All job API providers
- Open source community

---

## 📞 Support

For support, email: syedsayan467@gmail.com or open an issue on GitHub.

---

## 🎉 Features Coming Soon

- [ ] Resume templates and suggestions
- [ ] Multi-language support
- [ ] Resume comparison feature
- [ ] Interview preparation tips
- [ ] Cover letter generation
- [ ] LinkedIn profile optimization
- [ ] Email notifications for new jobs
- [ ] Resume version history
- [ ] Team collaboration features

---

Made with ❤️ by the HireLens Team
