import React from "react";
import { Github, Sparkles, Shield, Zap } from "lucide-react";

const Login = () => {
  const handleLogin = () => {
    window.location.href = "http://localhost:8080/auth/github/login";
  };

  return (
    <div className="min-h-screen bg-white flex items-center justify-center p-4">
      {/* Subtle background pattern */}
      <div className="absolute inset-0 bg-gradient-to-br from-purple-50 via-white to-blue-50"></div>
      
      {/* Main content */}
      <div className="relative z-10 max-w-md w-full">
        {/* Logo/Icon section */}
        <div className="text-center mb-10">
          <div className="inline-flex items-center justify-center w-20 h-20 bg-gradient-to-br from-purple-600 to-blue-600 rounded-2xl mb-5 shadow-lg transform hover:scale-110 transition-transform duration-300">
            <Sparkles className="w-10 h-10 text-white" />
          </div>
          <h1 className="text-5xl font-bold text-gray-900 mb-3 tracking-tight">
            AI Code Reviewer
          </h1>
          <p className="text-gray-600 text-lg">
            Automated code quality checks powered by Gemini
          </p>
        </div>

        {/* Features cards */}
        <div className="space-y-3 mb-8">
          <div className="bg-white rounded-xl p-5 border border-gray-200 shadow-sm hover:shadow-md hover:border-purple-300 transition-all duration-300">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-purple-100 rounded-xl flex items-center justify-center flex-shrink-0">
                <Shield className="w-6 h-6 text-purple-600" />
              </div>
              <div>
                <h3 className="text-gray-900 font-semibold text-base">Secure Integration</h3>
                <p className="text-gray-600 text-sm">Connect directly with GitHub</p>
              </div>
            </div>
          </div>

          <div className="bg-white rounded-xl p-5 border border-gray-200 shadow-sm hover:shadow-md hover:border-blue-300 transition-all duration-300">
            <div className="flex items-center gap-4">
              <div className="w-12 h-12 bg-blue-100 rounded-xl flex items-center justify-center flex-shrink-0">
                <Zap className="w-6 h-6 text-blue-600" />
              </div>
              <div>
                <h3 className="text-gray-900 font-semibold text-base">Instant Analysis</h3>
                <p className="text-gray-600 text-sm">Real-time code quality feedback</p>
              </div>
            </div>
          </div>
        </div>

        {/* Login button */}
        <button
          onClick={handleLogin}
          className="w-full bg-gray-900 hover:bg-gray-800 text-white font-semibold py-4 px-6 rounded-xl shadow-lg transform hover:scale-105 hover:shadow-xl transition-all duration-300 flex items-center justify-center gap-3 group"
        >
          <Github className="w-5 h-5 group-hover:rotate-12 transition-transform duration-300" />
          <span>Sign in with GitHub</span>
        </button>

        {/* Footer text */}
        <p className="text-center text-gray-500 text-sm mt-6">
          By signing in, you agree to our terms and privacy policy
        </p>
      </div>
    </div>
  );
};

export default Login;