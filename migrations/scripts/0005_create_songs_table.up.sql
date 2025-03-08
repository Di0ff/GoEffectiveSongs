CREATE TABLE songs (
    id BIGSERIAL PRIMARY KEY,
    group VARCHAR(255) NOT NULL,
    song VARCHAR(255) NOT NULL,
    release_date TIMESTAMP NOT NULL,
    link VARCHAR(32779) UNIQUE,
    created_at TIMESTAMP NOT NULL,
    updated_at TIMESTAMP NOT NULL,
    deleted_at TIMESTAMP
);

CREATE INDEX idx_songs_group ON songs (group);
CREATE INDEX idx_songs_song ON songs (song);
CREATE INDEX idx_songs_deleted_at ON songs (deleted_at);