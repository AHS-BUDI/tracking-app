-- Tabel untuk tracking mouse
CREATE TABLE IF NOT EXISTS mouse_events (
    id SERIAL PRIMARY KEY,
    x_position INTEGER NOT NULL,
    y_position INTEGER NOT NULL,
    event_type VARCHAR(20) NOT NULL,
    timestamp TIMESTAMP NOT NULL
);

-- Tabel untuk tracking keyboard
CREATE TABLE IF NOT EXISTS keyboard_events (
    id SERIAL PRIMARY KEY,
    key_pressed VARCHAR(50) NOT NULL,
    event_type VARCHAR(20) NOT NULL,
    timestamp TIMESTAMP NOT NULL
);