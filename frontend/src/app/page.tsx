'use client';

import Link from 'next/link';
import { useEffect, useState } from 'react';
import Header from '@/components/Header';

export default function Home() {
  const [mounted, setMounted] = useState(false);
  const [activeFeature, setActiveFeature] = useState(0);

  useEffect(() => {
    setMounted(true);
    
    // Auto-rotate features
    const interval = setInterval(() => {
      setActiveFeature((prev) => (prev + 1) % 3);
    }, 3000);

    return () => clearInterval(interval);
  }, []);

  return (
    <div className="min-h-screen bg-slate-950 relative overflow-hidden">
      {/* Animated Background */}
      <div className="absolute inset-0 overflow-hidden">
        <div className="absolute top-0 left-1/4 w-96 h-96 bg-primary-500/20 rounded-full blur-3xl animate-pulse"></div>
        <div className="absolute bottom-0 right-1/4 w-96 h-96 bg-secondary-500/20 rounded-full blur-3xl animate-pulse delay-1000"></div>
        <div className="absolute top-1/2 left-1/2 w-96 h-96 bg-purple-500/10 rounded-full blur-3xl animate-pulse delay-2000"></div>
      </div>

      <Header />
      
      <main className="container mx-auto px-4 sm:px-6 lg:px-8 relative z-10">
        {/* Hero Section */}
        <section className="py-20 md:py-32">
          <div className="max-w-5xl mx-auto text-center">
            {/* Animated Badge */}
            <div className={`inline-flex items-center gap-2 px-4 py-2 bg-primary-500/10 border border-primary-500/30 rounded-full mb-8 transition-all duration-1000 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 -translate-y-4'}`}>
              <span className="relative flex h-3 w-3">
                <span className="animate-ping absolute inline-flex h-full w-full rounded-full bg-primary-400 opacity-75"></span>
                <span className="relative inline-flex rounded-full h-3 w-3 bg-primary-500"></span>
              </span>
              <span className="text-sm text-primary-300 font-medium">AI-Powered Career Acceleration</span>
            </div>

            {/* Main Heading with Gradient Animation */}
            <h1 className={`text-5xl md:text-7xl lg:text-8xl font-bold mb-6 transition-all duration-1000 delay-100 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
              <span className="inline-block bg-gradient-to-r from-primary-400 via-secondary-400 to-purple-400 bg-clip-text text-transparent animate-gradient bg-[length:200%_auto]">
                Transform Your
              </span>
              <br />
              <span className="inline-block bg-gradient-to-r from-purple-400 via-pink-400 to-primary-400 bg-clip-text text-transparent animate-gradient bg-[length:200%_auto]">
                Career Journey
              </span>
            </h1>
            
            <p className={`text-xl md:text-2xl text-slate-300 mb-8 max-w-3xl mx-auto leading-relaxed transition-all duration-1000 delay-200 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
              Get <span className="text-primary-400 font-semibold">instant ATS scores</span>, 
              <span className="text-secondary-400 font-semibold"> skill gap analysis</span>, and 
              <span className="text-purple-400 font-semibold"> personalized job recommendations</span> powered by advanced AI
            </p>
            
            <div className={`flex flex-col sm:flex-row gap-4 justify-center transition-all duration-1000 delay-300 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
              <Link
                href="/register"
                className="group relative px-8 py-4 bg-gradient-to-r from-primary-500 to-secondary-500 text-white font-semibold rounded-lg overflow-hidden transition-all duration-300 hover:scale-105 hover:shadow-[0_0_40px_rgba(59,130,246,0.5)] text-lg"
              >
                <span className="relative z-10 flex items-center justify-center gap-2">
                  Get Started Free
                  <svg className="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                  </svg>
                </span>
                <div className="absolute inset-0 bg-gradient-to-r from-primary-600 to-secondary-600 opacity-0 group-hover:opacity-100 transition-opacity"></div>
              </Link>
              <Link
                href="/login"
                className="group px-8 py-4 bg-slate-800/50 backdrop-blur-sm text-slate-100 font-semibold rounded-lg border-2 border-slate-700 hover:border-primary-500 hover:bg-slate-800 transition-all duration-300 hover:scale-105 hover:shadow-[0_0_30px_rgba(59,130,246,0.3)] text-lg"
              >
                <span className="flex items-center justify-center gap-2">
                  Sign In
                  <svg className="w-5 h-5 group-hover:translate-x-1 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M11 16l-4-4m0 0l4-4m-4 4h14m-5 4v1a3 3 0 01-3 3H6a3 3 0 01-3-3V7a3 3 0 013-3h7a3 3 0 013 3v1" />
                  </svg>
                </span>
              </Link>
            </div>

            {/* Stats Section */}
            <div className={`mt-16 grid grid-cols-3 gap-8 max-w-3xl mx-auto transition-all duration-1000 delay-500 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
              <div className="text-center">
                <div className="text-4xl md:text-5xl font-bold text-primary-400 mb-2">
                  98%
                </div>
                <div className="text-sm text-slate-400">ATS Compatibility</div>
              </div>
              <div className="text-center border-x border-slate-800">
                <div className="text-4xl md:text-5xl font-bold text-secondary-400 mb-2">
                  10K+
                </div>
                <div className="text-sm text-slate-400">Jobs Analyzed</div>
              </div>
              <div className="text-center">
                <div className="text-4xl md:text-5xl font-bold text-purple-400 mb-2">
                  5s
                </div>
                <div className="text-sm text-slate-400">Analysis Time</div>
              </div>
            </div>
          </div>
        </section>

        {/* Features Section */}
        <section className="py-20">
          <div className="max-w-6xl mx-auto">
            <h2 className={`text-3xl md:text-5xl font-bold text-center mb-4 transition-all duration-1000 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
              <span className="text-primary-400">
                Powerful AI Features
              </span>
            </h2>
            <p className={`text-center text-slate-400 mb-16 text-lg transition-all duration-1000 delay-100 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
              Everything you need to land your dream job
            </p>
            
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8">
              {/* Feature 1 - ATS Score */}
              <div 
                className={`group relative bg-slate-900/50 backdrop-blur-sm rounded-2xl p-8 border border-slate-800 hover:border-primary-500 transition-all duration-500 hover:scale-105 hover:shadow-[0_0_50px_rgba(59,130,246,0.3)] cursor-pointer ${activeFeature === 0 ? 'border-primary-500 scale-105 shadow-[0_0_50px_rgba(59,130,246,0.3)]' : ''} ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}
                style={{ transitionDelay: '200ms' }}
                onMouseEnter={() => setActiveFeature(0)}
              >
                <div className="absolute inset-0 bg-gradient-to-br from-primary-500/10 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <div className="relative z-10">
                  <div className="w-16 h-16 rounded-xl bg-gradient-to-br from-blue-500 to-cyan-500 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                    <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 19v-6a2 2 0 00-2-2H5a2 2 0 00-2 2v6a2 2 0 002 2h2a2 2 0 002-2zm0 0V9a2 2 0 012-2h2a2 2 0 012 2v10m-6 0a2 2 0 002 2h2a2 2 0 002-2m0 0V5a2 2 0 012-2h2a2 2 0 012 2v14a2 2 0 01-2 2h-2a2 2 0 01-2-2z" />
                    </svg>
                  </div>
                  <h3 className="text-2xl font-bold text-white mb-4 group-hover:text-primary-400 transition-colors">
                    ATS Score Analysis
                  </h3>
                  <p className="text-slate-400 leading-relaxed">
                    Get instant ATS compatibility scores with detailed feedback on formatting, keywords, and structure optimization
                  </p>
                  <div className="mt-6 flex items-center text-primary-400 opacity-0 group-hover:opacity-100 transition-opacity">
                    <span className="text-sm font-semibold">Learn more</span>
                    <svg className="w-4 h-4 ml-1 group-hover:translate-x-2 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  </div>
                </div>
              </div>

              {/* Feature 2 - Skill Gap */}
              <div 
                className={`group relative bg-slate-900/50 backdrop-blur-sm rounded-2xl p-8 border border-slate-800 hover:border-purple-500 transition-all duration-500 hover:scale-105 hover:shadow-[0_0_50px_rgba(168,85,247,0.3)] cursor-pointer ${activeFeature === 1 ? 'border-purple-500 scale-105 shadow-[0_0_50px_rgba(168,85,247,0.3)]' : ''} ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}
                style={{ transitionDelay: '400ms' }}
                onMouseEnter={() => setActiveFeature(1)}
              >
                <div className="absolute inset-0 bg-gradient-to-br from-purple-500/10 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <div className="relative z-10">
                  <div className="w-16 h-16 rounded-xl bg-gradient-to-br from-purple-500 to-pink-500 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                    <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4M7.835 4.697a3.42 3.42 0 001.946-.806 3.42 3.42 0 014.438 0 3.42 3.42 0 001.946.806 3.42 3.42 0 013.138 3.138 3.42 3.42 0 00.806 1.946 3.42 3.42 0 010 4.438 3.42 3.42 0 00-.806 1.946 3.42 3.42 0 01-3.138 3.138 3.42 3.42 0 00-1.946.806 3.42 3.42 0 01-4.438 0 3.42 3.42 0 00-1.946-.806 3.42 3.42 0 01-3.138-3.138 3.42 3.42 0 00-.806-1.946 3.42 3.42 0 010-4.438 3.42 3.42 0 00.806-1.946 3.42 3.42 0 013.138-3.138z" />
                    </svg>
                  </div>
                  <h3 className="text-2xl font-bold text-white mb-4 group-hover:text-purple-400 transition-colors">
                    Skill Gap Analysis
                  </h3>
                  <p className="text-slate-400 leading-relaxed">
                    Identify matching skills and improvement areas to perfectly align your resume with target positions
                  </p>
                  <div className="mt-6 flex items-center text-purple-400 opacity-0 group-hover:opacity-100 transition-opacity">
                    <span className="text-sm font-semibold">Learn more</span>
                    <svg className="w-4 h-4 ml-1 group-hover:translate-x-2 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  </div>
                </div>
              </div>

              {/* Feature 3 - Job Recommendations */}
              <div 
                className={`group relative bg-slate-900/50 backdrop-blur-sm rounded-2xl p-8 border border-slate-800 hover:border-green-500 transition-all duration-500 hover:scale-105 hover:shadow-[0_0_50px_rgba(34,197,94,0.3)] cursor-pointer ${activeFeature === 2 ? 'border-green-500 scale-105 shadow-[0_0_50px_rgba(34,197,94,0.3)]' : ''} ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-8'}`}
                style={{ transitionDelay: '600ms' }}
                onMouseEnter={() => setActiveFeature(2)}
              >
                <div className="absolute inset-0 bg-gradient-to-br from-green-500/10 to-transparent rounded-2xl opacity-0 group-hover:opacity-100 transition-opacity duration-500"></div>
                <div className="relative z-10">
                  <div className="w-16 h-16 rounded-xl bg-gradient-to-br from-green-500 to-emerald-500 flex items-center justify-center mb-6 group-hover:scale-110 transition-transform duration-300 shadow-lg">
                    <svg className="w-8 h-8 text-white" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M21 13.255A23.931 23.931 0 0112 15c-3.183 0-6.22-.62-9-1.745M16 6V4a2 2 0 00-2-2h-4a2 2 0 00-2 2v2m4 6h.01M5 20h14a2 2 0 002-2V8a2 2 0 00-2-2H5a2 2 0 00-2 2v10a2 2 0 002 2z" />
                    </svg>
                  </div>
                  <h3 className="text-2xl font-bold text-white mb-4 group-hover:text-green-400 transition-colors">
                    Smart Job Matching
                  </h3>
                  <p className="text-slate-400 leading-relaxed">
                    Receive AI-powered job recommendations from multiple sources tailored to your unique skills and experience
                  </p>
                  <div className="mt-6 flex items-center text-green-400 opacity-0 group-hover:opacity-100 transition-opacity">
                    <span className="text-sm font-semibold">Learn more</span>
                    <svg className="w-4 h-4 ml-1 group-hover:translate-x-2 transition-transform" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 5l7 7-7 7" />
                    </svg>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* How It Works */}
        <section className="py-20">
          <div className="max-w-5xl mx-auto">
            <h2 className={`text-3xl md:text-5xl font-bold text-center mb-4 transition-all duration-1000 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
              <span className="text-secondary-400">
                How It Works
              </span>
            </h2>
            <p className={`text-center text-slate-400 mb-16 text-lg transition-all duration-1000 delay-100 ${mounted ? 'opacity-100 translate-y-0' : 'opacity-0 translate-y-4'}`}>
              Get results in 3 simple steps
            </p>
            
            <div className="relative">
              {/* Connection Line */}
              <div className="absolute left-8 top-12 bottom-12 w-0.5 bg-gradient-to-b from-primary-500 via-purple-500 to-green-500 hidden md:block"></div>
              
              <div className="space-y-12">
                {/* Step 1 */}
                <div className={`relative flex items-start gap-6 transition-all duration-1000 delay-200 ${mounted ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-8'}`}>
                  <div className="flex-shrink-0 w-16 h-16 rounded-2xl bg-gradient-to-br from-primary-500 to-blue-600 flex items-center justify-center text-white font-bold text-2xl shadow-lg shadow-primary-500/50 relative z-10">
                    <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M7 16a4 4 0 01-.88-7.903A5 5 0 1115.9 6L16 6a5 5 0 011 9.9M15 13l-3-3m0 0l-3 3m3-3v12" />
                    </svg>
                  </div>
                  <div className="flex-1 bg-slate-900/50 backdrop-blur-sm rounded-2xl p-6 border border-slate-800 hover:border-primary-500 transition-all duration-300 group">
                    <h3 className="text-2xl font-bold text-white mb-2 group-hover:text-primary-400 transition-colors">
                      Upload Your Resume
                    </h3>
                    <p className="text-slate-400 leading-relaxed">
                      Simply drag and drop your resume in PDF format. Add an optional job description for targeted analysis and better matching results.
                    </p>
                  </div>
                </div>

                {/* Step 2 */}
                <div className={`relative flex items-start gap-6 transition-all duration-1000 delay-400 ${mounted ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-8'}`}>
                  <div className="flex-shrink-0 w-16 h-16 rounded-2xl bg-gradient-to-br from-purple-500 to-pink-600 flex items-center justify-center text-white font-bold text-2xl shadow-lg shadow-purple-500/50 relative z-10">
                    <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 10V3L4 14h7v7l9-11h-7z" />
                    </svg>
                  </div>
                  <div className="flex-1 bg-slate-900/50 backdrop-blur-sm rounded-2xl p-6 border border-slate-800 hover:border-purple-500 transition-all duration-300 group">
                    <h3 className="text-2xl font-bold text-white mb-2 group-hover:text-purple-400 transition-colors">
                      AI Analysis
                    </h3>
                    <p className="text-slate-400 leading-relaxed">
                      Our advanced AI engine analyzes your resume for ATS compatibility, extracts skills, evaluates experience, and identifies improvement areas in seconds.
                    </p>
                  </div>
                </div>

                {/* Step 3 */}
                <div className={`relative flex items-start gap-6 transition-all duration-1000 delay-600 ${mounted ? 'opacity-100 translate-x-0' : 'opacity-0 -translate-x-8'}`}>
                  <div className="flex-shrink-0 w-16 h-16 rounded-2xl bg-gradient-to-br from-green-500 to-emerald-600 flex items-center justify-center text-white font-bold text-2xl shadow-lg shadow-green-500/50 relative z-10">
                    <svg className="w-8 h-8" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M9 12l2 2 4-4m6 2a9 9 0 11-18 0 9 9 0 0118 0z" />
                    </svg>
                  </div>
                  <div className="flex-1 bg-slate-900/50 backdrop-blur-sm rounded-2xl p-6 border border-slate-800 hover:border-green-500 transition-all duration-300 group">
                    <h3 className="text-2xl font-bold text-white mb-2 group-hover:text-green-400 transition-colors">
                      Get Results & Jobs
                    </h3>
                    <p className="text-slate-400 leading-relaxed">
                      Receive your ATS score, comprehensive skill analysis, actionable recommendations, and curated job opportunities that match your profile perfectly.
                    </p>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        {/* CTA Section */}
        <section className="py-20">
          <div className={`max-w-5xl mx-auto relative transition-all duration-1000 ${mounted ? 'opacity-100 scale-100' : 'opacity-0 scale-95'}`}>
            {/* Glow Effect */}
            <div className="absolute inset-0 bg-gradient-to-r from-primary-500/20 via-purple-500/20 to-secondary-500/20 blur-3xl"></div>
            
            <div className="relative bg-slate-900/80 backdrop-blur-xl rounded-3xl p-12 md:p-16 border border-slate-800 overflow-hidden">
              {/* Animated Background Pattern */}
              <div className="absolute inset-0 opacity-10">
                <div className="absolute top-0 left-0 w-40 h-40 bg-primary-500 rounded-full blur-3xl animate-pulse"></div>
                <div className="absolute bottom-0 right-0 w-40 h-40 bg-secondary-500 rounded-full blur-3xl animate-pulse delay-1000"></div>
              </div>
              
              <div className="relative z-10 text-center">
                <h2 className="text-3xl md:text-5xl font-bold text-white mb-4">
                  Ready to <span className="text-primary-400">Accelerate</span> Your Career?
                </h2>
                <p className="text-xl text-slate-300 mb-8 max-w-2xl mx-auto">
                  Join thousands of professionals who transformed their job search with AI-powered resume optimization
                </p>
                
                {/* Benefits */}
                <div className="flex flex-wrap justify-center gap-6 mb-10">
                  <div className="flex items-center gap-2 text-slate-300">
                    <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span>Free Forever</span>
                  </div>
                  <div className="flex items-center gap-2 text-slate-300">
                    <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span>Instant Results</span>
                  </div>
                  <div className="flex items-center gap-2 text-slate-300">
                    <svg className="w-5 h-5 text-green-400" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                      <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M5 13l4 4L19 7" />
                    </svg>
                    <span>No Credit Card</span>
                  </div>
                </div>
                
                <Link
                  href="/register"
                  className="group inline-flex items-center gap-3 px-10 py-5 bg-gradient-to-r from-primary-500 via-purple-500 to-secondary-500 text-white font-bold rounded-xl hover:shadow-[0_0_60px_rgba(59,130,246,0.6)] transition-all duration-300 hover:scale-105 text-lg relative overflow-hidden"
                >
                  <span className="relative z-10">Start Free Analysis Now</span>
                  <svg className="w-6 h-6 group-hover:translate-x-2 transition-transform relative z-10" fill="none" stroke="currentColor" viewBox="0 0 24 24">
                    <path strokeLinecap="round" strokeLinejoin="round" strokeWidth={2} d="M13 7l5 5m0 0l-5 5m5-5H6" />
                  </svg>
                  <div className="absolute inset-0 bg-gradient-to-r from-primary-600 via-purple-600 to-secondary-600 opacity-0 group-hover:opacity-100 transition-opacity"></div>
                </Link>
                
                <p className="mt-6 text-sm text-slate-400">
                  ⚡ Takes less than 30 seconds to get started
                </p>
              </div>
            </div>
          </div>
        </section>
      </main>

      {/* Footer */}
      <footer className="relative z-10 border-t border-slate-800 py-12 mt-20">
        <div className="container mx-auto px-4 sm:px-6 lg:px-8">
          <div className="max-w-6xl mx-auto">
            <div className="grid grid-cols-1 md:grid-cols-4 gap-8 mb-8">
              {/* Brand */}
              <div className="col-span-1 md:col-span-2">
                <div className="flex items-center space-x-2 mb-4">
                  <div className="w-10 h-10 rounded-lg bg-gradient-to-r from-primary-500 to-secondary-500 flex items-center justify-center">
                    <span className="text-white font-bold text-xl">H</span>
                  </div>
                  <span className="text-2xl font-bold text-primary-400">
                    HireLens
                  </span>
                </div>
                <p className="text-slate-400 mb-4 max-w-md">
                  AI-powered resume analysis and job matching platform helping professionals land their dream careers.
                </p>
                <div className="flex gap-4">
                  <a href="#" className="w-10 h-10 rounded-full bg-slate-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors">
                    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M24 12.073c0-6.627-5.373-12-12-12s-12 5.373-12 12c0 5.99 4.388 10.954 10.125 11.854v-8.385H7.078v-3.47h3.047V9.43c0-3.007 1.792-4.669 4.533-4.669 1.312 0 2.686.235 2.686.235v2.953H15.83c-1.491 0-1.956.925-1.956 1.874v2.25h3.328l-.532 3.47h-2.796v8.385C19.612 23.027 24 18.062 24 12.073z"/>
                    </svg>
                  </a>
                  <a href="#" className="w-10 h-10 rounded-full bg-slate-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors">
                    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M23.953 4.57a10 10 0 01-2.825.775 4.958 4.958 0 002.163-2.723c-.951.555-2.005.959-3.127 1.184a4.92 4.92 0 00-8.384 4.482C7.69 8.095 4.067 6.13 1.64 3.162a4.822 4.822 0 00-.666 2.475c0 1.71.87 3.213 2.188 4.096a4.904 4.904 0 01-2.228-.616v.06a4.923 4.923 0 003.946 4.827 4.996 4.996 0 01-2.212.085 4.936 4.936 0 004.604 3.417 9.867 9.867 0 01-6.102 2.105c-.39 0-.779-.023-1.17-.067a13.995 13.995 0 007.557 2.209c9.053 0 13.998-7.496 13.998-13.985 0-.21 0-.42-.015-.63A9.935 9.935 0 0024 4.59z"/>
                    </svg>
                  </a>
                  <a href="#" className="w-10 h-10 rounded-full bg-slate-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors">
                    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M20.447 20.452h-3.554v-5.569c0-1.328-.027-3.037-1.852-3.037-1.853 0-2.136 1.445-2.136 2.939v5.667H9.351V9h3.414v1.561h.046c.477-.9 1.637-1.85 3.37-1.85 3.601 0 4.267 2.37 4.267 5.455v6.286zM5.337 7.433c-1.144 0-2.063-.926-2.063-2.065 0-1.138.92-2.063 2.063-2.063 1.14 0 2.064.925 2.064 2.063 0 1.139-.925 2.065-2.064 2.065zm1.782 13.019H3.555V9h3.564v11.452zM22.225 0H1.771C.792 0 0 .774 0 1.729v20.542C0 23.227.792 24 1.771 24h20.451C23.2 24 24 23.227 24 22.271V1.729C24 .774 23.2 0 22.222 0h.003z"/>
                    </svg>
                  </a>
                  <a href="#" className="w-10 h-10 rounded-full bg-slate-800 hover:bg-slate-700 flex items-center justify-center text-slate-400 hover:text-white transition-colors">
                    <svg className="w-5 h-5" fill="currentColor" viewBox="0 0 24 24">
                      <path d="M12 0C5.374 0 0 5.373 0 12c0 5.302 3.438 9.8 8.207 11.387.599.111.793-.261.793-.577v-2.234c-3.338.726-4.033-1.416-4.033-1.416-.546-1.387-1.333-1.756-1.333-1.756-1.089-.745.083-.729.083-.729 1.205.084 1.839 1.237 1.839 1.237 1.07 1.834 2.807 1.304 3.492.997.107-.775.418-1.305.762-1.604-2.665-.305-5.467-1.334-5.467-5.931 0-1.311.469-2.381 1.236-3.221-.124-.303-.535-1.524.117-3.176 0 0 1.008-.322 3.301 1.23A11.509 11.509 0 0112 5.803c1.02.005 2.047.138 3.006.404 2.291-1.552 3.297-1.23 3.297-1.23.653 1.653.242 2.874.118 3.176.77.84 1.235 1.911 1.235 3.221 0 4.609-2.807 5.624-5.479 5.921.43.372.823 1.102.823 2.222v3.293c0 .319.192.694.801.576C20.566 21.797 24 17.3 24 12c0-6.627-5.373-12-12-12z"/>
                    </svg>
                  </a>
                </div>
              </div>

              {/* Quick Links */}
              <div>
                <h4 className="text-white font-semibold mb-4">Quick Links</h4>
                <ul className="space-y-2">
                  <li><Link href="/dashboard" className="text-slate-400 hover:text-primary-400 transition-colors">Dashboard</Link></li>
                  <li><Link href="/resumes" className="text-slate-400 hover:text-primary-400 transition-colors">My Resumes</Link></li>
                  <li><Link href="/register" className="text-slate-400 hover:text-primary-400 transition-colors">Get Started</Link></li>
                  <li><Link href="/login" className="text-slate-400 hover:text-primary-400 transition-colors">Sign In</Link></li>
                </ul>
              </div>

              {/* Resources */}
              <div>
                <h4 className="text-white font-semibold mb-4">Resources</h4>
                <ul className="space-y-2">
                  <li><a href="#" className="text-slate-400 hover:text-primary-400 transition-colors">Help Center</a></li>
                  <li><a href="#" className="text-slate-400 hover:text-primary-400 transition-colors">Privacy Policy</a></li>
                  <li><a href="#" className="text-slate-400 hover:text-primary-400 transition-colors">Terms of Service</a></li>
                  <li><a href="#" className="text-slate-400 hover:text-primary-400 transition-colors">Contact Us</a></li>
                </ul>
              </div>
            </div>

            <div className="border-t border-slate-800 pt-8">
              <p className="text-center text-slate-400">
                © 2025 <span className="text-primary-400">HireLens</span>. All rights reserved. Built with ❤️ for job seekers worldwide.
              </p>
            </div>
          </div>
        </div>
      </footer>
    </div>
  );
}
