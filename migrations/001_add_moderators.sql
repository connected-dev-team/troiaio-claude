-- Migration: Add moderators table
-- Run this SQL on your PostgreSQL database before starting the dashboard

CREATE TABLE IF NOT EXISTS moderators (
    id SERIAL PRIMARY KEY,
    username TEXT UNIQUE NOT NULL,
    passwd_hash TEXT NOT NULL,
    name TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT NOW()
);

-- Insert default admin user
-- Username: admin
-- Password: admin123 (MD5 hash)
INSERT INTO moderators (username, passwd_hash, name)
VALUES ('admin', '0192023a7bbd73250516f069df18b500', 'Amministratore')
ON CONFLICT (username) DO NOTHING;
