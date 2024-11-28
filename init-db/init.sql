-- init-db/init.sql

CREATE TABLE IF NOT EXISTS time_log (
    id INT AUTO_INCREMENT PRIMARY KEY,
    timestamp TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Insert an initial record if needed
INSERT INTO time_log (timestamp) VALUES (NOW());
