CREATE TABLE IF NOT EXISTS configurations (
    id SERIAL PRIMARY KEY,
    repository_id INT REFERENCES repositories(id) ON DELETE CASCADE,
    ignore_patterns TEXT[], -- Array of strings, e.g. ["*.css", "vendor/*"]
    review_style VARCHAR(50) DEFAULT 'concise', -- 'concise', 'detailed', 'socratic'
    enabled BOOLEAN DEFAULT TRUE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);