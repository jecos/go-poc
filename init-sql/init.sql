CREATE DATABASE IF NOT EXISTS test_db;

USE test_db;

CREATE TABLE IF NOT EXISTS test_table (
    id INT,
    name VARCHAR(255)
) PRIMARY KEY (id);

INSERT INTO test_table (id, name) VALUES (1, 'test1'), (2, 'test2');