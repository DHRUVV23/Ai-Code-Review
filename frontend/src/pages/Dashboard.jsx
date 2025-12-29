import React from "react";
import { Link } from "react-router-dom";
import { GitPullRequest, Search, Bell, Settings, LogOut } from "lucide-react";

const Dashboard = () => {
  // Mock data for now
  const repos = [
    { id: 1, name: "DHRUVV23/OmniScribe", status: "Active", pendingReviews: 2 },
    { id: 2, name: "DHRUVV23/ai-code-review", status: "Inactive", pendingReviews: 0 },
  ];

  return (
    <div className="min-h-screen bg-gray-50">
      {/* Top Navigation */}
      <nav className="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between sticky top-0 z-10">
        <div className="flex items-center gap-3">
          <div className="bg-purple-600 p-2 rounded-lg">
            <GitPullRequest className="w-6 h-6 text-white" />
          </div>
          <h1 className="text-xl font-bold text-gray-900">AI Code Reviewer</h1>
        </div>
        <div className="flex items-center gap-4">
          <div className="relative hidden md:block">
            <Search className="w-5 h-5 absolute left-3 top-1/2 transform -translate-y-1/2 text-gray-400" />
            <input 
              type="text" 
              placeholder="Search repositories..." 
              className="pl-10 pr-4 py-2 border border-gray-200 rounded-lg focus:outline-none focus:ring-2 focus:ring-purple-500 w-64"
            />
          </div>
          <button className="p-2 hover:bg-gray-100 rounded-full relative">
            <Bell className="w-5 h-5 text-gray-600" />
            <span className="absolute top-1 right-1 w-2 h-2 bg-red-500 rounded-full"></span>
          </button>
          <div className="w-8 h-8 bg-gradient-to-tr from-purple-500 to-blue-500 rounded-full border-2 border-white shadow-sm"></div>
        </div>
      </nav>

      {/* Main Content */}
      <main className="max-w-7xl mx-auto px-6 py-8">
        <div className="flex justify-between items-end mb-8">
          <div>
            <h2 className="text-2xl font-bold text-gray-900">Repositories</h2>
            <p className="text-gray-500 mt-1">Manage your connected projects and reviews</p>
          </div>
          <button className="bg-gray-900 hover:bg-gray-800 text-white px-5 py-2.5 rounded-lg font-medium transition-colors shadow-sm flex items-center gap-2">
            <Settings className="w-4 h-4" />
            Configure New Repo
          </button>
        </div>

        {/* Repository Grid */}
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
          {repos.map((repo) => (
            <Link 
              key={repo.id} 
              to={`/repo/${repo.id}`}
              className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm hover:shadow-md hover:border-purple-300 transition-all duration-300 group"
            >
              <div className="flex justify-between items-start mb-4">
                <div className="bg-purple-50 p-3 rounded-lg group-hover:bg-purple-100 transition-colors">
                  <GitPullRequest className="w-6 h-6 text-purple-600" />
                </div>
                <span className={`px-3 py-1 rounded-full text-xs font-medium ${
                  repo.status === 'Active' 
                    ? 'bg-green-100 text-green-700' 
                    : 'bg-gray-100 text-gray-600'
                }`}>
                  {repo.status}
                </span>
              </div>
              
              <h3 className="text-lg font-bold text-gray-900 mb-2 group-hover:text-purple-600 transition-colors">
                {repo.name}
              </h3>
              
              <div className="flex items-center justify-between text-sm text-gray-500 mt-4 pt-4 border-t border-gray-100">
                <span>Last checked: 2h ago</span>
                <span className="flex items-center gap-1">
                  {repo.pendingReviews} pending
                </span>
              </div>
            </Link>
          ))}
        </div>
      </main>
    </div>
  );
};

export default Dashboard;