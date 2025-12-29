import React from "react";
import { Link } from "react-router-dom";

const Dashboard = () => {
  // Hardcoded for now - we will fetch this from API next
  const repos = [
    { id: 1, name: "DHRUVV23/OmniScribe", status: "Active" },
  ];

  return (
    <div style={{ padding: "20px" }}>
      <h2>📂 Your Repositories</h2>
      <ul>
        {repos.map((repo) => (
          <li key={repo.id} style={{ marginBottom: "10px" }}>
            <Link to={`/repo/${repo.id}`} style={{ textDecoration: "none", fontSize: "18px" }}>
              {repo.name}
            </Link>
             <span style={{ marginLeft: "10px", color: "green" }}>● {repo.status}</span>
          </li>
        ))}
      </ul>
    </div>
  );
};

export default Dashboard;