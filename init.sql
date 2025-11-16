CREATE TABLE IF NOT EXISTS teams (
    id SERIAL PRIMARY KEY,
    name VARCHAR(50) UNIQUE NOT NULL
);
CREATE TABLE IF NOT EXISTS users (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(30) UNIQUE NOT NULL,
    username VARCHAR(30) UNIQUE NOT NULL,
    team_id INTEGER REFERENCES teams (id),
    is_active BOOLEAN NOT NULL
);
CREATE TYPE pr_status AS ENUM ('OPEN', 'MERGED');
CREATE TABLE IF NOT EXISTS pull_requests (
    pr_id VARCHAR(30) PRIMARY KEY,
    pr_name VARCHAR(50) NOT NULL,
    author_id VARCHAR(30) REFERENCES users (user_id),
    assigned_reviewers VARCHAR(30)[] NOT NULL,
    status pr_status DEFAULT 'OPEN',
    merged_at TIMESTAMP WITH TIME ZONE
);