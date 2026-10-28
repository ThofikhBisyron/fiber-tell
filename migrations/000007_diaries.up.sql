CREATE EXTENSION IF NOT EXISTS pg_trgm;
CREATE TABLE diaries (
    id BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    mood_id BIGINT NOT NULL REFERENCES moods(id),
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    date DATE NOT NULL,
    search_vector TSVECTOR GENERATED ALWAYS AS (
        setweight(
            to_tsvector('simple', COALESCE(title, '')),
            'A'
        ) ||
        setweight(
            to_tsvector('simple', COALESCE(content, '')),
            'B'
        )
    ) STORED,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_diaries_user_date
ON diaries(user_id, date DESC);

CREATE INDEX idx_diaries_mood_id
ON diaries(mood_id);

CREATE INDEX idx_diaries_search_vector
ON diaries USING GIN (search_vector);

CREATE INDEX idx_diaries_title_trgm
ON diaries USING GIN (title gin_trgm_ops);

CREATE INDEX idx_diaries_content_trgm
ON diaries USING GIN (content gin_trgm_ops);