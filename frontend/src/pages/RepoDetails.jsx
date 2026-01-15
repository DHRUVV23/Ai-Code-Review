import React, { useEffect, useState } from "react";
import { useParams, Link } from "react-router-dom";
import axios from "axios";
import { ArrowLeft, Settings, FileCode, AlertCircle, CheckCircle, Loader } from "lucide-react";

const RepoDetails = () => {
  const { id } = useParams();
  const [repo, setRepo] = useState(null);
  const [config, setConfig] = useState(null);
  const [loading, setLoading] = useState(true);
  
  const [connecting, setConnecting] = useState(false);
  const [connected, setConnected] = useState(false);

  useEffect(() => {
    fetchRepoDetails();
  }, [id]);

  const fetchRepoDetails = async () => {
    const token = localStorage.getItem("auth_token");
    try {
      const listResp = await axios.get("http://localhost:8080/api/v1/user/repositories", {
        headers: { Authorization: `Bearer ${token}` }
      });
      const foundRepo = listResp.data.find(r => r.id === parseInt(id));
      setRepo(foundRepo);

      const configResp = await axios.get(`http://localhost:8080/api/v1/repositories/${id}`, {
        headers: { Authorization: `Bearer ${token}` }
      });
      setConfig(configResp.data);

    } catch (error) {
      console.error("Failed to load repo details", error);
    } finally {
      setLoading(false);
    }
  };


  const handleConnect = async () => {
    setConnecting(true);
    const token = localStorage.getItem("auth_token");
    
    try {
      
      await axios.post(`http://localhost:8080/api/v1/repositories/${id}/webhook`, {}, {
        headers: { Authorization: `Bearer ${token}` }
      });
      
      setConnected(true);
      alert("Webhook successfully created on GitHub!");
    } catch (error) {
      const msg = error.response?.data?.error || "Failed to connect";
      alert("Error: " + msg);
      
     
      if (msg.includes("token not found")) {
         alert("Please Log Out and Log In again to save your GitHub permissions.");
      }
    } finally {
      setConnecting(false);
    }
  };

  if (loading) return <div className="p-10 text-center">Loading Control Center...</div>;
  if (!repo) return <div className="p-10 text-center text-red-500">Repository not found.</div>;

  return (
    <div className="min-h-screen bg-gray-50 p-8">
      <div className="max-w-4xl mx-auto">
        <Link to="/dashboard" className="flex items-center text-gray-500 hover:text-gray-900 mb-6">
          <ArrowLeft className="w-4 h-4 mr-2" /> Back to Dashboard
        </Link>

        {/* Title Card */}
        <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-8 mb-6">
          <div className="flex justify-between items-start">
            <div>
              <h1 className="text-3xl font-bold text-gray-900">{repo.owner} / {repo.name}</h1>
              <p className="text-gray-500 mt-2 flex items-center gap-2">
                <span className="bg-green-100 text-green-700 px-2 py-0.5 rounded-full text-sm font-medium">Active</span>
                <span className="text-sm">• Added on {new Date(repo.created_at).toLocaleDateString()}</span>
              </p>
            </div>
            <div className="bg-purple-50 p-3 rounded-full">
              <FileCode className="w-8 h-8 text-purple-600" />
            </div>
          </div>
        </div>

        <div className="grid grid-cols-1 md:grid-cols-2 gap-6">
          {/* Configuration */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
            <div className="flex items-center gap-2 mb-4">
              <Settings className="w-5 h-5 text-gray-600" />
              <h2 className="text-xl font-semibold text-gray-900">AI Configuration</h2>
            </div>
            {config ? (
              <div className="space-y-4">
                 <div>
                    <label className="text-xs font-bold text-gray-500 uppercase tracking-wide">Review Style</label>
                    <div className="mt-1 font-medium text-gray-900 capitalize">{config.review_style}</div>
                 </div>
                 <div>
                    <label className="text-xs font-bold text-gray-500 uppercase tracking-wide">Ignored Files</label>
                    <div className="mt-1 font-mono text-sm bg-gray-100 p-2 rounded text-gray-700">
                      {config.ignore_patterns || "None"}
                    </div>
                 </div>
              </div>
            ) : <p>Loading...</p>}
          </div>

          {/* Connection Status */}
          <div className="bg-white rounded-xl shadow-sm border border-gray-200 p-6">
             <div className="flex items-center gap-2 mb-4">
              {connected ? (
                 <CheckCircle className="w-5 h-5 text-green-600" />
              ) : (
                 <AlertCircle className="w-5 h-5 text-orange-600" />
              )}
              <h2 className="text-xl font-semibold text-gray-900">Connection Status</h2>
            </div>
            
            <p className="text-gray-600 mb-6">
              {connected 
                ? "Your repository is connected! The AI is listening for Pull Requests." 
                : "The webhook is not connected yet. The AI cannot see your Pull Requests."}
            </p>

            {!connected ? (
              <button 
                onClick={handleConnect}
                disabled={connecting}
                className="w-full bg-black text-white font-bold py-2 rounded-lg hover:bg-gray-800 transition shadow-lg flex justify-center items-center gap-2"
              >
                 {connecting && <Loader className="w-4 h-4 animate-spin" />}
                 {connecting ? "Connecting..." : "Connect to GitHub"}
              </button>
            ) : (
              <button className="w-full bg-green-100 text-green-800 font-bold py-2 rounded-lg cursor-default border border-green-200">
                 Connected via Webhook
              </button>
            )}
          </div>
        </div>
      </div>
    </div>
  );
};

export default RepoDetails;