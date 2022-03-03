CREATE TABLE meetups(
    id SERIAL PRIMARY KEY,
    name VARCHAR(127) UNIQUE NOT NULL,
    description TEXT,
    user_id SERIAL REFERENCES users(id) ON DELETE CASCADE NOT NULL
);