CREATE TABLE user_provider (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO user_provider (name) VALUES
('google'),
('facebook');