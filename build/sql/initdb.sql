DROP TABLE IF EXISTS url;

------------------------------------------------------------------------
CREATE TABLE IF NOT EXISTS urls (
    short VARCHAR(255) PRIMARY KEY,
    original VARCHAR(255),
    expiration_time TIMESTAMP
);