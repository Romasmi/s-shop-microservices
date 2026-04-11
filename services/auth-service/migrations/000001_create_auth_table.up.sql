-- 000001_create_auth_table.up.sql
CREATE TABLE IF NOT EXISTS auth_logs (
    id SERIAL PRIMARY KEY,
    user_id VARCHAR(255),
    login VARCHAR(255),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
