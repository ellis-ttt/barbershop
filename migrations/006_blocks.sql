CREATE TABLE blocks (
    id BIGSERIAL PRIMARY KEY,
    barber_id BIGINT NOT NULL REFERENCES barbers(id),
    starts_at TIMESTAMPTZ NOT NULL,
    ends_at TIMESTAMPTZ NOT NULL,
    reason TEXT NOT NULL DEFAULT ''
);
