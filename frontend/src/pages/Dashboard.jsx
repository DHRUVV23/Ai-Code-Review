import React, { useEffect, useState } from "react";
import { Link, useSearchParams, useNavigate } from "react-router-dom";
import axios from "axios";
import { GitPullRequest, Search, Bell, Plus, LogOut, X, Trash2 } from "lucide-react";

const Dashboard = () => {
  const [searchParams] = useSearchParams();
  const navigate = useNavigate();
  
  const [user, setUser] = useState(localStorage.getItem("user_name") || "Guest");
  const [repos, setRepos] = useState([]);
  const [loading, setLoading] = useState(true);

  const [isModalOpen, setIsModalOpen] = useState(false);
  const [newRepoOwner, setNewRepoOwner] = useState("");
  const [newRepoName, setNewRepoName] = useState("");

  useEffect(() => {
    const tokenFromURL = searchParams.get("token");
    const userFromURL = searchParams.get("user");
    if (tokenFromURL) {
      localStorage.setItem("auth_token", tokenFromURL);
      localStorage.setItem("user_name", userFromURL);
      setUser(userFromURL);
      window.history.replaceState({}, document.title, "/dashboard");
    }
    fetchRepos();
  }, [searchParams]);

  const fetchRepos = async () => {
    const token = localStorage.getItem("auth_token");
    if (!token) return;
    try {
      const response = await axios.get("http://localhost:8080/api/v1/user/repositories", {
        headers: { Authorization: `Bearer ${token}` }
      });
      setRepos(response.data || []);
    } catch (error) {
      console.error("Failed to fetch repos", error);
    } finally {
      setLoading(false);
    }
  };

  const handleAddRepo = async (e) => {
    e.preventDefault();
    const token = localStorage.getItem("auth_token");
    
    try {
      await axios.post("http://localhost:8080/api/v1/repositories", 
        { owner: newRepoOwner, name: newRepoName },
        { headers: { Authorization: `Bearer ${token}` } }
      );
      
      setIsModalOpen(false);
      setNewRepoName("");
      setNewRepoOwner("");
      fetchRepos(); 
      alert("Repository Added Successfully!");
    } catch (err) {
      alert("Failed to add repository: " + err.message);
    }
  };

  const handleDelete = async (repoId) => {
    if (!confirm("Are you sure? This will delete the webhook and unlink the repo.")) {
      return;
    }

    const token = localStorage.getItem("auth_token");
    try {
      await axios.delete(`http://localhost:8080/api/v1/repositories/${repoId}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      
      setRepos(repos.filter((r) => r.id !== repoId));
      alert("Repository deleted successfully!");
    } catch (err) {
      console.error(err);
      alert("Failed to delete repository.");
    }
  };

  const handleLogout = () => {
    localStorage.clear();
    window.location.href = "/";
  };

  return (
    <div className="min-h-screen bg-gray-50 relative">
      <nav className="bg-white border-b border-gray-200 px-6 py-4 flex items-center justify-between sticky top-0 z-10">
        <div className="flex items-center gap-3">
          <div className="bg-purple-600 p-2 rounded-lg">
            <GitPullRequest className="w-6 h-6 text-white" />
          </div>
          <h1 className="text-xl font-bold text-gray-900">AI Code Reviewer</h1>
        </div>
        <div className="flex items-center gap-4">
           <span className="text-gray-700 font-medium">Hi , Dhruv</span>
           <button onClick={handleLogout} className="p-2 hover:bg-gray-100 rounded-full"><LogOut className="w-5 h-5 text-gray-600" /></button>
        </div>
      </nav>

      <main className="max-w-7xl mx-auto px-6 py-8">
        <div className="flex justify-between items-end mb-8">
          <h2 className="text-2xl font-bold text-gray-900">Repositories</h2>
          <button 
            onClick={() => setIsModalOpen(true)}
            className="bg-gray-900 hover:bg-gray-800 text-white px-5 py-2.5 rounded-lg font-medium flex items-center gap-2"
          >
            <Plus className="w-4 h-4" /> Add Repository
          </button>
        </div>

        {loading ? <p>Loading...</p> : repos.length === 0 ? (
          <div className="text-center py-20 bg-white rounded-xl border border-dashed border-gray-300">
             <p className="text-gray-500">No repositories yet.</p>
          </div>
        ) : (
          <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-6">
            {repos.map((repo) => (
              
              <div key={repo.id} className="bg-white p-6 rounded-xl border border-gray-200 shadow-sm hover:shadow-md transition-all relative group">
                <Link to={`/repo/${repo.id}`} className="block">
                    <h3 className="text-lg font-bold text-gray-900 hover:text-purple-600 transition-colors">
                        {repo.owner} / {repo.name}
                    </h3>
                    <span className="text-xs bg-green-100 text-green-700 px-2 py-1 rounded-full mt-2 inline-block">Active</span>
                </Link>

               
                <button
                    onClick={() => handleDelete(repo.id)}
                    className="absolute top-4 right-4 p-2 text-gray-400 hover:text-red-600 hover:bg-red-50 rounded-full transition-colors opacity-0 group-hover:opacity-100"
                    title="Delete Repository"
                >
                    <Trash2 className="w-5 h-5" />
                </button>
              </div>
            ))}
          </div>
        )}
      </main>

    
      {isModalOpen && (
        <div className="fixed inset-0 bg-black bg-opacity-50 flex items-center justify-center z-50">
          <div className="bg-white rounded-xl p-6 w-full max-w-md shadow-2xl">
            <div className="flex justify-between items-center mb-4">
              <h3 className="text-xl font-bold">Add New Repository</h3>
              <button onClick={() => setIsModalOpen(false)}><X className="w-5 h-5 text-gray-500" /></button>
            </div>
            
            <form onSubmit={handleAddRepo}>
              <div className="space-y-4">
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Owner (Username)</label>
                  <input 
                    type="text" 
                    placeholder="e.g. DHRUVV23"
                    className="w-full border border-gray-300 rounded-lg p-2 focus:ring-2 focus:ring-purple-500"
                    value={newRepoOwner}
                    onChange={(e) => setNewRepoOwner(e.target.value)}
                    required
                  />
                </div>
                <div>
                  <label className="block text-sm font-medium text-gray-700 mb-1">Repository Name</label>
                  <input 
                    type="text" 
                    placeholder="e.g. ai-code-review"
                    className="w-full border border-gray-300 rounded-lg p-2 focus:ring-2 focus:ring-purple-500"
                    value={newRepoName}
                    onChange={(e) => setNewRepoName(e.target.value)}
                    required
                  />
                </div>
                <button type="submit" className="w-full bg-purple-600 hover:bg-purple-700 text-white font-bold py-2 rounded-lg transition-colors">
                  Add Repository
                </button>
              </div>
            </form>
          </div>
        </div>
      )}
    </div>
  );
};

export default Dashboard;