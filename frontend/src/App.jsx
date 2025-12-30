import React from "react";
import { BrowserRouter as Router, Routes, Route } from "react-router-dom";
import Login from "./pages/Login";
import Dashboard from "./pages/Dashboard";
import RepoReviews from "./pages/RepoReviews";
import RepoDetails from "./pages/RepoDetails";
import './index.css';

function App() {
  return (
    <Router>
      <Routes>
        <Route path="/" element={<Login />} />
        <Route path="/dashboard" element={<Dashboard />} />
        {/* <Route path="/repo/:id" element={<RepoReviews />} /> */}
        <Route path="/repo/:id" element={<RepoDetails />} />
      </Routes>
    </Router>
  );
}

export default App;