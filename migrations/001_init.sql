CREATE TABLE IF NOT EXISTS services (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE,
    price INT NOT NULL,
    duration_min INT NOT NULL
);

CREATE TABLE IF NOT EXISTS barbers (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL UNIQUE
);

INSERT INTO services (name, price, duration_min) VALUES
    ('Стрижка', 1500, 60),
    ('Борода', 800, 30)
ON CONFLICT (name) DO NOTHING;

INSERT INTO barbers (name) VALUES ('Иван'), ('Пётр')
ON CONFLICT (name) DO NOTHING;
