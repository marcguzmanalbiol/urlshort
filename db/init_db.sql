DROP TABLE IF EXISTS url_mapping;

CREATE TABLE url_mapping (
    short_url VARCHAR(255) PRIMARY KEY,
    original_url VARCHAR(255),
    used_count INT,
    expiration_time TIMESTAMP
);