CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    firstname TEXT NOT NULL,
    lastname TEXT NOT NULL,
    age INT NOT NULL
);