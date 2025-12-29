import React from "react";

const Login = () => {
  const handleLogin = () => {
    // We will connect this to real OAuth later
    window.location.href = "http://localhost:8080/auth/github/login";
  };

  return (
    <div style={{ textAlign: "center", marginTop: "100px" }}>
      <h1>🤖 AI Code Reviewer</h1>
      <p>Automated code quality checks powered by Gemini.</p>
      <button 
        onClick={handleLogin}
        style={{ padding: "10px 20px", fontSize: "16px", cursor: "pointer", backgroundColor: "#333", color: "#fff", border: "none", borderRadius: "5px" }}
      >
        Sign in with GitHub
      </button>
    </div>
  );
};

export default Login;