CREATE DATABASE maxbot;

\c maxbot;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name TEXT
);
