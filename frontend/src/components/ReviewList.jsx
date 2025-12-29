import { useEffect, useState } from "react";
import axios from "axios";

const ReviewList = ({ repoId }) => {
  const [reviews, setReviews] = useState([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    // Note: Assuming repoId 1 for now if not passed
    const id = repoId || 1; 
    axios
      .get(`http://localhost:8080/api/v1/repositories/${id}/reviews`)
      .then((response) => {
        setReviews(response.data);
        setLoading(false);
      })
      .catch((err) => console.error(err));
  }, [repoId]);

  const parseContent = (content) => {
    try {
      const parsed = JSON.parse(content);
      return Array.isArray(parsed) ? parsed : []; 
    } catch (e) {
      return [{ message: content, type: "legacy" }];
    }
  };

  const getSeverityColor = (severity) => {
    switch (severity?.toLowerCase()) {
      case "high": return "#ffcccc";
      case "medium": return "#fff4cc";
      case "low": return "#ccffcc";
      default: return "#f9f9f9";
    }
  };

  if (loading) return <h2>⏳ Loading Reviews...</h2>;

  return (
    <div style={{ padding: "20px", maxWidth: "800px", margin: "0 auto" }}>
      <h1>🤖 AI Code Reviews</h1>
      {reviews.length === 0 && <p>No reviews found.</p>}

      {reviews.map((review) => {
        const issues = parseContent(review.content);
        return (
          <div key={review.id} style={{ marginBottom: "30px", borderTop: "2px solid #eee", paddingTop: "20px" }}>
            <h3>Review #{review.id}</h3>
            {issues.map((issue, idx) => (
              <div key={idx} style={{
                  backgroundColor: getSeverityColor(issue.severity),
                  border: "1px solid #ddd", borderRadius: "8px", padding: "15px", marginBottom: "10px"
              }}>
                <h4 style={{ margin: "0 0 5px 0" }}>📄 {issue.file}</h4>
                <p><strong>{issue.message}</strong></p>
                {issue.suggestion && <code style={{ display:"block", background:"#fff", padding:"5px" }}>{issue.suggestion}</code>}
              </div>
            ))}
          </div>
        );
      })}
    </div>
  );
};

export default ReviewList;