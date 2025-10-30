# analyzer/app.py
from fastapi import FastAPI
from pydantic import BaseModel
from typing import List, Dict, Any
import spacy
import re
import os
from fastapi.middleware.cors import CORSMiddleware

# Load model once
MODEL = os.getenv("SPACY_MODEL", "en_core_web_sm")
nlp = spacy.load(MODEL)

app = FastAPI(title="Resume Analyzer (spaCy)")

app.add_middleware(
    CORSMiddleware,
    allow_origins=["*"],  # adjust for production
    allow_methods=["*"],
    allow_headers=["*"],
)

class AnalyzeRequest(BaseModel):
    text: str
    job_description: str = ""  # Optional job description for matching

class Entity(BaseModel):
    text: str
    label: str
    start: int
    end: int
    score: float = 1.0

class AnalyzeResponse(BaseModel):
    entities: List[Entity]
    skills: List[str]
    summary: str
    ats_score: int  # 0-100 percentage score
    jd_match_score: int  # 0-100 percentage score for JD matching
    matching_skills: List[str]  # Skills found in both resume and JD
    missing_skills: List[str]  # Skills in JD but not in resume

# simple default tech keywords (extend as needed)
DEFAULT_SKILLS = [
    # Programming Languages
    "python", "java", "c", "c++", "c#", "javascript", "typescript", "go", "golang",
    "rust", "ruby", "swift", "kotlin", "scala", "dart", "php", "r", "perl", "objective-c",
    "bash", "shell", "powershell", "haskell", "elixir", "lua", "matlab", "fortran", 
    # Frontend Frameworks & Libraries
    "react", "reactjs", "nextjs", "angular", "vue", "nuxtjs", "svelte", "solidjs",
    "jquery", "bootstrap", "tailwindcss", "chakraui", "materialui", "redux", "mobx",
    "lit", "astro", "vite", "webpack", "parcel", "babel",
    # Backend Frameworks
    "node", "express", "nestjs", "fastify", "django", "flask", "fastapi", "spring", 
    "springboot", "laravel", "rails", "ruby on rails", "asp.net", "dotnet", "gin", "fiber",
    "echo", "phoenix", "hapi", "adonisjs", "koajs",
    # Databases
    "sql", "postgresql", "mysql", "mariadb", "sqlite", "mongodb", "redis", "oracle",
    "cassandra", "elasticsearch", "dynamodb", "couchdb", "neo4j", "firebase", "supabase",
    "prisma", "typeorm", "sequelize", "hibernate", "mongoose", "realm", "influxdb",
    # Cloud & DevOps
    "aws", "gcp", "azure", "digitalocean", "heroku", "vercel", "netlify", "render",
    "docker", "kubernetes", "terraform", "ansible", "jenkins", "github actions",
    "gitlab ci", "circleci", "travisci", "argo cd", "helm", "prometheus", "grafana",
    "nginx", "apache", "loadbalancer", "cdn", "serverless", "lambda", "cloudformation",
    # Version Control & Collaboration
    "git", "github", "gitlab", "bitbucket", "svn", "mercurial",
    # Data Science & Machine Learning
    "numpy", "pandas", "scikit-learn", "tensorflow", "pytorch", "keras", "matplotlib",
    "seaborn", "xgboost", "lightgbm", "catboost", "opencv", "nlp", "spacy", "transformers",
    "huggingface", "statsmodels", "jupyter", "notebook", "colab", "data visualization",
    "mlflow", "kubeflow", "pytorch lightning", "deep learning", "computer vision",
    "machine learning", "artificial intelligence", "reinforcement learning",
    # Data Engineering & Big Data
    "hadoop", "spark", "pyspark", "kafka", "airflow", "luigi", "snowflake", "bigquery",
    "databricks", "redshift", "data lake", "data pipeline", "etl", "elt", "presto",
    "hive", "flink", "storm",
    # Mobile & Cross-Platform
    "react native", "flutter", "swiftui", "android", "ios", "xcode", "kivy", "ionic",
    "cordova", "capacitor",
    # AI / NLP / CV
    "openai", "langchain", "llm", "chatgpt", "gpt", "bert", "gpt-4", "t5", "transformer",
    "yolo", "cnn", "rnn", "gans", "stable diffusion", "speech recognition", "ocr",
    "image classification", "nlp pipeline", "text generation",
    # Testing & QA
    "jest", "mocha", "chai", "enzyme", "cypress", "playwright", "puppeteer", "pytest",
    "unittest", "postman", "newman", "selenium", "robot framework",
    # Cybersecurity & Networking
    "penetration testing", "ethical hacking", "owasp", "burpsuite", "metasploit",
    "firewall", "wireshark", "nmap", "ssl", "tls", "encryption", "jwt", "oauth",
    "sso", "networking", "vpn", "zero trust", "iam",
    # Blockchain & Web3
    "blockchain", "ethereum", "solidity", "web3", "smart contracts", "nft", "defi",
    "metamask", "ethersjs", "hardhat", "truffle", "ipfs", "polygon", "solana",
    # Misc Tools / Others
    "restapi", "graphql", "grpc", "websocket", "mqtt", "rabbitmq", "kafka", "celery",
    "redis queue", "microservices", "monorepo", "turborepo", "api gateway",
    "swagger", "openapi", "postman", "insomnia", "linux", "ubuntu", "windows server",
    "macos", "bash scripting", "automation", "devops", "agile", "scrum", "jira",
    "confluence", "figma", "adobe xd", "ui/ux", "design systems",
    # Game Development
    "unity", "unreal engine", "godot", "blender", "threejs", "babylonjs",
    # Emerging Technologies
    "genai", "rag", "autogen", "agentic ai", "ai agent", "digital twin",
    "iot", "embedded systems", "arduino", "raspberry pi", "robotics", "edge computing",
    # Analytics & BI
    "tableau", "powerbi", "looker", "metabase", "superset", "google data studio",
    # Misc Development Skills
    "performance optimization", "scalability", "system design", "api design",
    "distributed systems", "event-driven architecture", "observability", "logging",
    "monitoring", "tracing"
]

# load additional keywords from file if present
SKILLS_LIST = []
SKILLS_FILE = os.getenv("SKILLS_FILE", "")  # optional path to newline-separated keywords
if SKILLS_FILE and os.path.isfile(SKILLS_FILE):
    with open(SKILLS_FILE, "r", encoding="utf-8") as f:
        SKILLS_LIST = [line.strip().lower() for line in f if line.strip()]
else:
    SKILLS_LIST = DEFAULT_SKILLS

# extract the skills from the resume pdf
def extract_skills(text: str) -> List[str]:
    text_l = text.lower()
    found = set()
    # simple substring match — robust and fast
    for kw in SKILLS_LIST:
        if kw in text_l:
            found.add(kw)
    return sorted(found)

# extract keywords from resume
def extract_keywords(text: str, min_word_length: int = 3) -> set:
    """Extract important keywords from text (nouns, proper nouns, skills)"""
    text_lower = text.lower()
    doc = nlp(text_lower)
    
    keywords = set()
    
    # Extract nouns and proper nouns
    for token in doc:
        if token.pos_ in ['NOUN', 'PROPN'] and len(token.text) >= min_word_length:
            if not token.is_stop:
                keywords.add(token.lemma_)
    
    # Extract named entities
    for ent in doc.ents:
        if ent.label_ in ['ORG', 'PRODUCT', 'SKILL', 'LANGUAGE']:
            keywords.add(ent.text.lower())
    
    # Add known skills
    for skill in SKILLS_LIST:
        if skill in text_lower:
            keywords.add(skill)
    
    return keywords

# give matching score for job description
def calculate_jd_match_score(resume_text: str, job_description: str, resume_skills: List[str]) -> tuple:
    """
    Calculate job description match score (0-30 points for ATS, 0-100 for JD match)
    Returns: (ats_points, jd_match_percentage, matching_skills, missing_skills)
    """
    if not job_description or len(job_description.strip()) < 10:
        # No JD provided, return neutral score
        return (0, 0, [], [])
    
    score = 0
    
    # Extract keywords from both texts
    resume_keywords = extract_keywords(resume_text)
    jd_keywords = extract_keywords(job_description)
    
    # Extract skills from JD
    jd_skills = set(extract_skills(job_description))
    resume_skills_set = set(resume_skills)
    
    # Calculate matching and missing skills
    matching_skills = sorted(list(resume_skills_set & jd_skills))
    missing_skills = sorted(list(jd_skills - resume_skills_set))
    
    # 1. Skill matching (15 points max for ATS)
    skill_match_ratio = 0
    if jd_skills:
        skill_match_ratio = len(matching_skills) / len(jd_skills)
        score += int(skill_match_ratio * 15)
    
    # 2. Keyword matching (15 points max for ATS)
    keyword_match_ratio = 0
    if jd_keywords:
        keyword_match_ratio = len(resume_keywords & jd_keywords) / len(jd_keywords)
        score += int(keyword_match_ratio * 15)
    
    # Calculate overall JD match percentage (0-100)
    jd_match_percentage = int(((skill_match_ratio + keyword_match_ratio) / 2) * 100)
    
    return (min(score, 30), jd_match_percentage, matching_skills, missing_skills)

# generate a simple summary of the resume
def simple_summary(text: str, max_chars=600) -> str:
    # naive summary: first N chars / first 3 sentences
    # We'll pick first 3 sentences as short summary
    sentences = re.split(r'(?<=[.!?])\s+', text.strip())
    summary = " ".join(sentences[:3])
    if len(summary) > max_chars:
        return summary[:max_chars].rsplit(' ',1)[0] + "..."
    return summary

# calculate ATS score
def calculate_ats_score(text: str, skills: List[str], entities: List[Entity], job_description: str = "") -> tuple:
    """
    Calculate ATS score (0-100) based on:
    1. Skills match (30 points max) - reduced from 40 to make room for JD match
    2. Experience/Education detection (25 points max) - reduced from 30
    3. Resume structure/sections (15 points max) - reduced from 30
    4. Job Description match (30 points max) - NEW
    
    If no JD provided, scores are redistributed to maintain fairness
    
    Returns: (ats_score, jd_match_score, matching_skills, missing_skills)
    """
    score = 0
    text_lower = text.lower()
    has_jd = job_description and len(job_description.strip()) >= 10
    
    # Initialize JD match variables
    jd_match_score = 0
    matching_skills = []
    missing_skills = []
    
    # 1. Skills Match (30 points with JD, 40 points without JD)
    skills_count = len(skills)
    max_skills_points = 30 if has_jd else 40
    
    if skills_count >= 10:
        score += max_skills_points
    elif skills_count >= 7:
        score += int(max_skills_points * 0.875)
    elif skills_count >= 5:
        score += int(max_skills_points * 0.75)
    elif skills_count >= 3:
        score += int(max_skills_points * 0.5)
    elif skills_count >= 1:
        score += int(max_skills_points * 0.25)
    
    # 2. Experience & Education Detection (25 points with JD, 30 points without JD)
    max_exp_edu_points = 25 if has_jd else 30
    experience_keywords = ['experience', 'work history', 'employment', 'worked at', 'position', 'role', 'job']
    education_keywords = ['education', 'degree', 'university', 'college', 'bachelor', 'master', 'phd', 'diploma', 'graduated']
    certification_keywords = ['certification', 'certified', 'certificate', 'license']
    
    has_experience = any(keyword in text_lower for keyword in experience_keywords)
    has_education = any(keyword in text_lower for keyword in education_keywords)
    has_certifications = any(keyword in text_lower for keyword in certification_keywords)
    
    if has_experience:
        score += int(max_exp_edu_points * 0.48)
    if has_education:
        score += int(max_exp_edu_points * 0.48)
    if has_certifications:
        score += int(max_exp_edu_points * 0.24)
    
    # 3. Resume Structure (15 points with JD, 30 points without JD)
    max_structure_points = 15 if has_jd else 30
    sections = {
        'summary': ['summary', 'objective', 'profile', 'about'],
        'projects': ['projects', 'portfolio', 'work samples'],
        'skills': ['skills', 'technical skills', 'competencies', 'expertise'],
        'contact': ['email', 'phone', 'linkedin', 'github', 'contact'],
        'achievements': ['achievements', 'awards', 'honors', 'accomplishments']
    }
    
    sections_found = 0
    for section_name, keywords in sections.items():
        if any(keyword in text_lower for keyword in keywords):
            sections_found += 1
    
    # Award points based on sections found
    if sections_found >= 5:
        score += max_structure_points
    elif sections_found >= 4:
        score += int(max_structure_points * 0.833)
    elif sections_found >= 3:
        score += int(max_structure_points * 0.667)
    elif sections_found >= 2:
        score += int(max_structure_points * 0.5)
    elif sections_found >= 1:
        score += int(max_structure_points * 0.333)
    
    # 4. Job Description Match (30 points max) - only if JD provided
    if has_jd:
        jd_ats_points, jd_match_score, matching_skills, missing_skills = calculate_jd_match_score(text, job_description, skills)
        score += jd_ats_points
    
    # Ensure score is within 0-100 range
    ats_score = min(max(score, 0), 100)
    
    return (ats_score, jd_match_score, matching_skills, missing_skills)

@app.post("/analyze", response_model=AnalyzeResponse)
def analyze(req: AnalyzeRequest):
    text = req.text or ""
    
    # Limit text size to prevent memory issues
    MAX_TEXT_LENGTH = 100000  # 100KB
    if len(text) > MAX_TEXT_LENGTH:
        text = text[:MAX_TEXT_LENGTH]
    
    try:
        doc = nlp(text)

        entities = []
        for ent in doc.ents:
            entities.append(Entity(
                text=ent.text,
                label=ent.label_,
                start=ent.start_char,
                end=ent.end_char,
                score=1.0
            ))

        skills = extract_skills(text)
        summary = simple_summary(text)
        
        # Calculate ATS score with optional job description
        ats_score, jd_match_score, matching_skills, missing_skills = calculate_ats_score(text, skills, entities, req.job_description)

        return AnalyzeResponse(
            entities=entities, 
            skills=skills, 
            summary=summary, 
            ats_score=ats_score,
            jd_match_score=jd_match_score,
            matching_skills=matching_skills,
            missing_skills=missing_skills
        )
    except Exception as e:
        print(f"Error analyzing text: {e}")
        # Return empty result instead of error
        return AnalyzeResponse(
            entities=[], 
            skills=[], 
            summary="Error analyzing resume", 
            ats_score=0,
            jd_match_score=0,
            matching_skills=[],
            missing_skills=[]
        )
