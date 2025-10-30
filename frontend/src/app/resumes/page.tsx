'use client';

import { useState, useEffect } from 'react';
import { useRouter } from 'next/navigation';
import Header from '@/components/Header';
import { resumeAPI, parseSkills, parseAnalysisResult, type Resume, type JobRecommendation } from '@/lib/api';

export default function ResumesPage() {
  const [resumes, setResumes] = useState<Resume[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [selectedResume, setSelectedResume] = useState<Resume | null>(null);
  const [jobs, setJobs] = useState<JobRecommendation[]>([]);
  const [loadingJobs, setLoadingJobs] = useState(false);
  const [showDeleteConfirm, setShowDeleteConfirm] = useState<number | null>(null);
  const router = useRouter();

  useEffect(() => {
    fetchResumes();
  }, []);

  const fetchResumes = async () => {
    try {
      setLoading(true);
      const response = await resumeAPI.getUserResumes();
      setResumes(response.resumes || []);
    } catch (err: any) {
      setError(err.message || 'Failed to load resumes');
      if (err.message?.includes('unauthorized')) {
        router.push('/login');
      }
    } finally {
      setLoading(false);
    }
  };

  const fetchResumeJobs = async (resumeId: number) => {
    try {
      setLoadingJobs(true);
      const response = await resumeAPI.getResumeJobs(resumeId);
      setJobs(response.jobs || []);
    } catch (err: any) {
      console.error('Failed to load jobs:', err);
      setJobs([]);
    } finally {
      setLoadingJobs(false);
    }
  };

  const handleViewDetails = (resume: Resume) => {
    setSelectedResume(resume);
    fetchResumeJobs(resume.id);
  };

  const handleCloseDetails = () => {
    setSelectedResume(null);
    setJobs([]);
  };

  const handleDeleteResume = async (id: number) => {
    try {
      await resumeAPI.deleteResume(id);
      setResumes(resumes.filter(r => r.id !== id));
      setShowDeleteConfirm(null);
      if (selectedResume?.id === id) {
        handleCloseDetails();
      }
    } catch (err: any) {
      alert(err.message || 'Failed to delete resume');
    }
  };

  const formatDate = (dateString: string) => {
    return new Date(dateString).toLocaleDateString('en-US', {
      year: 'numeric',
      month: 'long',
      day: 'numeric',
    });
  };

  if (loading) {
    return (
      <div className="min-h-screen bg-slate-950">
        <Header />
        <main className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
          <div className="flex items-center justify-center h-64">
            <div className="animate-spin rounded-full h-12 w-12 border-b-2 border-primary-500"></div>
          </div>
        </main>
      </div>
    );
  }

  return (
    <div className="min-h-screen bg-slate-950">
      <Header />
      <main className="container mx-auto px-4 sm:px-6 lg:px-8 py-8">
        {/* Page Header */}
        <div className="mb-8">
          <h1 className="text-3xl font-bold gradient-text mb-2">My Resumes</h1>
          <p className="text-slate-400">
            View and manage all your uploaded resumes and their job recommendations
          </p>
        </div>

        {error && (
          <div className="mb-6 p-4 bg-red-500/10 border border-red-500/50 rounded-lg">
            <p className="text-red-400">{error}</p>
          </div>
        )}

        {resumes.length === 0 ? (
          <div className="text-center py-12">
            <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-800 mb-4">
              <svg className="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12h6m-6 4h6m2 5H7a2 2 0 01-2-2V5a2 2 0 012-2h5.586a1 1 0 01.707.293l5.414 5.414a1 1 0 01.293.707V19a2 2 0 01-2 2z" />
              </svg>
            </div>
            <h3 className="text-xl font-semibold text-white mb-2">No resumes yet</h3>
            <p className="text-slate-400 mb-6">Upload your first resume to get started</p>
            <button
              onClick={() => router.push('/dashboard')}
              className="px-6 py-3 bg-gradient-to-r from-primary-500 to-secondary-500 text-white rounded-lg hover:shadow-lg transition-all duration-300"
            >
              Go to Dashboard
            </button>
          </div>
        ) : (
          <div className="grid grid-cols-1 lg:grid-cols-2 gap-6">
            {/* Resumes List */}
            <div className="space-y-4">
              {resumes.map((resume) => {
                const matchingSkills = parseSkills(resume.matching_skills);
                const missingSkills = parseSkills(resume.missing_skills);
                const analysis = parseAnalysisResult(resume.analysis_result);

                return (
                  <div
                    key={resume.id}
                    className={`p-6 rounded-lg border transition-all cursor-pointer ${
                      selectedResume?.id === resume.id
                        ? 'bg-slate-800/50 border-primary-500'
                        : 'bg-slate-900/50 border-slate-800 hover:border-slate-700'
                    }`}
                    onClick={() => handleViewDetails(resume)}
                  >
                    <div className="flex items-start justify-between mb-4">
                      <div className="flex-1">
                        <h3 className="text-lg font-semibold text-white mb-1">{resume.title}</h3>
                        <p className="text-sm text-slate-400">Uploaded {formatDate(resume.uploaded_at)}</p>
                      </div>
                      <button
                        onClick={(e) => {
                          e.stopPropagation();
                          setShowDeleteConfirm(resume.id);
                        }}
                        className="p-2 text-slate-400 hover:text-red-400 transition-colors"
                      >
                        <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                          <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M19 7l-.867 12.142A2 2 0 0116.138 21H7.862a2 2 0 01-1.995-1.858L5 7m5 4v6m4-6v6m1-10V4a1 1 0 00-1-1h-4a1 1 0 00-1 1v3M4 7h16" />
                        </svg>
                      </button>
                    </div>

                    {/* ATS Score */}
                    <div className="flex items-center gap-4 mb-4">
                      <div className="flex-1">
                        <div className="flex items-center justify-between mb-1">
                          <span className="text-sm text-slate-400">ATS Score</span>
                          <span className="text-sm font-semibold text-white">{resume.ats_score}%</span>
                        </div>
                        <div className="w-full bg-slate-800 rounded-full h-2">
                          <div
                            className={`h-2 rounded-full transition-all duration-500 ${
                              resume.ats_score >= 80
                                ? 'bg-green-500'
                                : resume.ats_score >= 60
                                ? 'bg-yellow-500'
                                : 'bg-red-500'
                            }`}
                            style={{ width: `${resume.ats_score}%` }}
                          ></div>
                        </div>
                      </div>
                      {resume.jd_match_score > 0 && (
                        <div className="flex-1">
                          <div className="flex items-center justify-between mb-1">
                            <span className="text-sm text-slate-400">JD Match</span>
                            <span className="text-sm font-semibold text-white">{resume.jd_match_score}%</span>
                          </div>
                          <div className="w-full bg-slate-800 rounded-full h-2">
                            <div
                              className={`h-2 rounded-full transition-all duration-500 ${
                                resume.jd_match_score >= 80
                                  ? 'bg-green-500'
                                  : resume.jd_match_score >= 60
                                  ? 'bg-yellow-500'
                                  : 'bg-red-500'
                              }`}
                              style={{ width: `${resume.jd_match_score}%` }}
                            ></div>
                          </div>
                        </div>
                      )}
                    </div>

                    {/* Skills Preview */}
                    {matchingSkills.length > 0 && (
                      <div className="flex flex-wrap gap-2">
                        {matchingSkills.slice(0, 4).map((skill, idx) => (
                          <span
                            key={idx}
                            className="px-2 py-1 text-xs bg-green-500/20 text-green-400 rounded border border-green-500/30"
                          >
                            {skill}
                          </span>
                        ))}
                        {matchingSkills.length > 4 && (
                          <span className="px-2 py-1 text-xs bg-slate-800 text-slate-400 rounded">
                            +{matchingSkills.length - 4} more
                          </span>
                        )}
                      </div>
                    )}

                    {/* Delete Confirmation */}
                    {showDeleteConfirm === resume.id && (
                      <div
                        className="mt-4 p-3 bg-red-500/10 border border-red-500/50 rounded"
                        onClick={(e) => e.stopPropagation()}
                      >
                        <p className="text-sm text-red-400 mb-3">Are you sure you want to delete this resume?</p>
                        <div className="flex gap-2">
                          <button
                            onClick={() => handleDeleteResume(resume.id)}
                            className="px-3 py-1 bg-red-500 text-white text-sm rounded hover:bg-red-600 transition-colors"
                          >
                            Delete
                          </button>
                          <button
                            onClick={(e) => {
                              e.stopPropagation();
                              setShowDeleteConfirm(null);
                            }}
                            className="px-3 py-1 bg-slate-700 text-white text-sm rounded hover:bg-slate-600 transition-colors"
                          >
                            Cancel
                          </button>
                        </div>
                      </div>
                    )}
                  </div>
                );
              })}
            </div>

            {/* Details Panel */}
            <div className="lg:sticky lg:top-24 h-fit">
              {selectedResume ? (
                <div className="p-6 bg-slate-900/50 border border-slate-800 rounded-lg">
                  <div className="flex items-start justify-between mb-6">
                    <h2 className="text-xl font-bold text-white">Resume Details</h2>
                    <button
                      onClick={handleCloseDetails}
                      className="p-1 text-slate-400 hover:text-white transition-colors lg:hidden"
                    >
                      <svg className="w-5 h-5" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                        <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M6 18L18 6M6 6l12 12" />
                      </svg>
                    </button>
                  </div>

                  <div className="space-y-6">
                    {/* Skills */}
                    <div>
                      <h3 className="text-sm font-semibold text-slate-300 mb-3">Matching Skills</h3>
                      <div className="flex flex-wrap gap-2">
                        {parseSkills(selectedResume.matching_skills).map((skill, idx) => (
                          <span
                            key={idx}
                            className="px-3 py-1 text-sm bg-green-500/20 text-green-400 rounded border border-green-500/30"
                          >
                            {skill}
                          </span>
                        ))}
                      </div>
                    </div>

                    {parseSkills(selectedResume.missing_skills).length > 0 && (
                      <div>
                        <h3 className="text-sm font-semibold text-slate-300 mb-3">Skills to Improve</h3>
                        <div className="flex flex-wrap gap-2">
                          {parseSkills(selectedResume.missing_skills).map((skill, idx) => (
                            <span
                              key={idx}
                              className="px-3 py-1 text-sm bg-red-500/20 text-red-400 rounded border border-red-500/30"
                            >
                              {skill}
                            </span>
                          ))}
                        </div>
                      </div>
                    )}

                    {/* Job Recommendations */}
                    <div>
                      <h3 className="text-sm font-semibold text-slate-300 mb-3">
                        Job Recommendations ({jobs.length})
                      </h3>
                      {loadingJobs ? (
                        <div className="flex items-center justify-center py-8">
                          <div className="animate-spin rounded-full h-8 w-8 border-b-2 border-primary-500"></div>
                        </div>
                      ) : jobs.length === 0 ? (
                        <p className="text-sm text-slate-400">No job recommendations available</p>
                      ) : (
                        <div className="space-y-3 max-h-96 overflow-y-auto">
                          {jobs.map((job) => (
                            <div
                              key={job.id}
                              className="p-4 bg-slate-800/50 border border-slate-700 rounded-lg hover:border-primary-500/50 transition-all"
                            >
                              <h4 className="font-semibold text-white mb-1">{job.title}</h4>
                              <p className="text-sm text-slate-300 mb-2">{job.company}</p>
                              <div className="flex items-center gap-3 text-xs text-slate-400 mb-3">
                                <span className="flex items-center gap-1">
                                  <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M17.657 16.657L13.414 20.9a1.998 1.998 0 01-2.827 0l-4.244-4.243a8 8 0 1111.314 0z" />
                                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 11a3 3 0 11-6 0 3 3 0 016 0z" />
                                  </svg>
                                  {job.location}
                                </span>
                                {job.salary && (
                                  <span className="flex items-center gap-1">
                                    <svg className="w-3 h-3" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M12 8c-1.657 0-3 .895-3 2s1.343 2 3 2 3 .895 3 2-1.343 2-3 2m0-8c1.11 0 2.08.402 2.599 1M12 8V7m0 1v8m0 0v1m0-1c-1.11 0-2.08-.402-2.599-1M21 12a9 9 0 11-18 0 9 9 0 0118 0z" />
                                    </svg>
                                    {job.salary}
                                  </span>
                                )}
                              </div>
                              <a
                                href={job.job_url}
                                target="_blank"
                                rel="noopener noreferrer"
                                className="inline-flex items-center gap-1 text-sm text-primary-400 hover:text-primary-300 transition-colors"
                              >
                                View Job
                                <svg className="w-4 h-4" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                                  <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M10 6H6a2 2 0 00-2 2v10a2 2 0 002 2h10a2 2 0 002-2v-4M14 4h6m0 0v6m0-6L10 14" />
                                </svg>
                              </a>
                            </div>
                          ))}
                        </div>
                      )}
                    </div>

                    {/* Download Resume */}
                    <a
                      href={selectedResume.file_url}
                      target="_blank"
                      rel="noopener noreferrer"
                      className="block w-full px-4 py-3 bg-primary-500 text-white text-center rounded-lg hover:bg-primary-600 transition-colors"
                    >
                      Download Resume
                    </a>
                  </div>
                </div>
              ) : (
                <div className="p-12 bg-slate-900/50 border border-slate-800 rounded-lg text-center">
                  <div className="inline-flex items-center justify-center w-16 h-16 rounded-full bg-slate-800 mb-4">
                    <svg className="w-8 h-8 text-slate-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M15 12a3 3 0 11-6 0 3 3 0 016 0z" />
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M2.458 12C3.732 7.943 7.523 5 12 5c4.478 0 8.268 2.943 9.542 7-1.274 4.057-5.064 7-9.542 7-4.477 0-8.268-2.943-9.542-7z" />
                    </svg>
                  </div>
                  <p className="text-slate-400">Select a resume to view details</p>
                </div>
              )}
            </div>
          </div>
        )}
      </main>
    </div>
  );
}
