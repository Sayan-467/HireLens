// API Configuration
const API_BASE_URL = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8080';

// Token management
export const getToken = (): string | null => {
  if (typeof window === 'undefined') return null;
  return localStorage.getItem('token');
};

export const setToken = (token: string): void => {
  localStorage.setItem('token', token);
};

export const removeToken = (): void => {
  localStorage.removeItem('token');
};

// API client with auth headers
const apiClient = async (
  endpoint: string,
  options: RequestInit = {}
): Promise<Response> => {
  const token = getToken();
  const headers: HeadersInit = {
    ...options.headers,
  };

  if (token) {
    headers['Authorization'] = `Bearer ${token}`;
  }

  // Don't set Content-Type for FormData (browser will set it with boundary)
  if (!(options.body instanceof FormData)) {
    headers['Content-Type'] = 'application/json';
  }

  const response = await fetch(`${API_BASE_URL}${endpoint}`, {
    ...options,
    headers,
  });

  return response;
};

// Types
export interface User {
  id: number;
  name: string;
  email: string;
  created_at: string;
}

export interface LoginResponse {
  message: string;
  token: string;
}

export interface RegisterResponse {
  message: string;
}

export interface ProfileResponse {
  name: string;
  email: string;
}

export interface RegisterData {
  name: string;
  email: string;
  password: string;
}

export interface LoginData {
  email: string;
  password: string;
}

export interface Resume {
  id: number;
  user_id: number;
  title: string;
  file_url: string;
  analysis_result: string;
  ats_score: number;
  jd_match_score: number;
  matching_skills: string; // JSON string array
  missing_skills: string; // JSON string array
  uploaded_at: string;
}

export interface JobRecommendation {
  id: number;
  resume_id: number;
  title: string;
  company: string;
  location: string;
  description: string;
  salary: string;
  job_url: string;
  posted_date: string;
  job_type: string;
  created_at: string;
}

export interface UploadResumeResponse {
  message: string;
  file_url: string;
  analysis_result: string;
  ats_score: number;
  jd_match_score: number;
  matching_skills: string;
  missing_skills: string;
  recommended_jobs: JobRecommendation[];
}

export interface AnalysisResult {
  ats_score: number;
  jd_match_score: number;
  skills: string[];
  matching_skills: string[];
  missing_skills: string[];
  experience_years: number;
  education_level: string;
  recommendations: string[];
}

// Auth API
export const authAPI = {
  register: async (data: RegisterData): Promise<RegisterResponse> => {
    const response = await apiClient('/api/signup', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Registration failed');
    }

    return response.json();
  },

  login: async (data: LoginData): Promise<LoginResponse> => {
    const response = await apiClient('/api/login', {
      method: 'POST',
      body: JSON.stringify(data),
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Login failed');
    }

    return response.json();
  },

  getProfile: async (): Promise<ProfileResponse> => {
    const response = await apiClient('/api/profile');

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to fetch profile');
    }

    return response.json();
  },

  logout: (): void => {
    removeToken();
  },
};

// Resume API
export const resumeAPI = {
  upload: async (
    file: File,
    title: string,
    jobDescription?: string
  ): Promise<UploadResumeResponse> => {
    const formData = new FormData();
    formData.append('resume', file);
    formData.append('title', title);
    if (jobDescription) {
      formData.append('job_description', jobDescription);
    }

    const response = await apiClient('/api/resume/upload', {
      method: 'POST',
      body: formData,
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Resume upload failed');
    }

    return response.json();
  },

  getUserResumes: async (): Promise<{ resumes: Resume[] }> => {
    const response = await apiClient('/api/resumes');

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to fetch resumes');
    }

    return response.json();
  },

  getResumeById: async (id: number): Promise<{ resume: Resume }> => {
    const response = await apiClient(`/api/resume/${id}`);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to fetch resume');
    }

    return response.json();
  },

  deleteResume: async (id: number): Promise<{ message: string }> => {
    const response = await apiClient(`/api/resume/${id}`, {
      method: 'DELETE',
    });

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to delete resume');
    }

    return response.json();
  },

  getResumeJobs: async (id: number): Promise<{ jobs: JobRecommendation[] }> => {
    const response = await apiClient(`/api/resume/${id}/jobs`);

    if (!response.ok) {
      const error = await response.json();
      throw new Error(error.error || 'Failed to fetch job recommendations');
    }

    return response.json();
  },
};

// Helper function to parse JSON strings from backend
export const parseSkills = (skillsJson: string): string[] => {
  try {
    return JSON.parse(skillsJson);
  } catch (e) {
    return [];
  }
};

// Helper function to parse analysis result
export const parseAnalysisResult = (analysisJson: string): AnalysisResult | null => {
  try {
    return JSON.parse(analysisJson);
  } catch (e) {
    return null;
  }
};
