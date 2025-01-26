-- Create the `users` table
CREATE TABLE IF NOT EXISTS users (
    id VARCHAR PRIMARY KEY,
    first_name VARCHAR(50),
    last_name VARCHAR(50),
    email VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255)
);

-- Insert a sample record
INSERT INTO users (id, first_name, last_name, email, password_hash)
VALUES ('49738873-a77e-47f6-ba1c-0d4840f9eabc', 'John', 'Doe', 'john.doe@example.com', 'hashed_password');
