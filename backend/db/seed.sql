-- 1. Fake User
INSERT INTO users (github_id, username, email) 
VALUES (999, 'testuser', 'test@example.com')
ON CONFLICT (github_id) DO NOTHING;

-- 2. Fake Installation
INSERT INTO installations (github_installation_id, account_id, account_type, account_login)
VALUES (8888, 999, 'User', 'testuser')
ON CONFLICT (github_installation_id) DO NOTHING;

-- 3. Fake Repository
INSERT INTO repositories (github_repo_id, installation_id, name, full_name, private)
VALUES (12345, (SELECT id FROM installations WHERE github_installation_id=8888), 'my-awesome-project', 'testuser/my-awesome-project', FALSE)
ON CONFLICT (github_repo_id) DO NOTHING;

-- 4. Fake Configuration (The table we just added!)
INSERT INTO configurations (repository_id, ignore_patterns, review_style)
VALUES ((SELECT id FROM repositories WHERE github_repo_id=12345), ARRAY['*.md', 'assets/*'], 'socratic');