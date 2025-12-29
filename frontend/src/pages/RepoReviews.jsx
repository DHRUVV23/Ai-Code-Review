import React from "react";
import { useParams, Link } from "react-router-dom";
import ReviewList from "../components/ReviewList";

const RepoReviews = () => {
  const { id } = useParams(); // Get ID from URL (e.g., /repo/1)

  return (
    <div>
      <div style={{ padding: "10px 20px", background: "#f4f4f4", borderBottom: "1px solid #ddd" }}>
        <Link to="/dashboard">← Back to Dashboard</Link>
      </div>
      <ReviewList repoId={id} />
    </div>
  );
};

export default RepoReviews;