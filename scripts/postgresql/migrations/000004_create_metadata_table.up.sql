BEGIN;

CREATE TABLE IF NOT EXISTS metadata (
    id serial primary KEY,
    player_id  INT NOT NULL,
    played_game_id INT NOT NULL,
    play_time INT NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_player
        FOREIGN KEY(player_id)
            REFERENCES users(id)
            ON DELETE CASCADE,
    CONSTRAINT fk_game
        FOREIGN KEY(played_game_id)
            REFERENCES games(id)
            ON DELETE CASCADE
);

COMMIT;