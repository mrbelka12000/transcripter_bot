CREATE TABLE IF NOT EXISTS messages (
    id SERIAL PRIMARY KEY,
    chat_id VARCHAR(50) NOT NULL,
    message_ID BIGINT NOT NULL,
    text    text  not null,
    search_vector tsvector,
    created_at  timestamp default now()
);

CREATE INDEX IF NOT EXISTS articles_search_vector_idx ON messages USING gin(search_vector);
