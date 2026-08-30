CREATE TABLE moods (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE
);

INSERT INTO moods (name) VALUES
('Very Sad'),
('Sad'),
('Neutral'),
('Happy'),
('Very Happy');

