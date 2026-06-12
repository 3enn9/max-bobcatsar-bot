CREATE DATABASE maxbot;

\c maxbot;

CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name TEXT
);

CREATE TABLE operations (
                            id BIGSERIAL PRIMARY KEY,
                            telegram_group_id BIGINT NOT NULL,
                            operation VARCHAR(50) NOT NULL,
                            description TEXT,
                            amount NUMERIC(15, 2) NOT NULL,
                            created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP
);
