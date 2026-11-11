ALTER TABLE bookings ADD COLUMN status TEXT NOT NULL DEFAULT 'active';

ALTER TABLE bookings DROP CONSTRAINT bookings_no_overlap;
ALTER TABLE bookings ADD CONSTRAINT bookings_no_overlap
    EXCLUDE USING gist (barber_id WITH =, tstzrange(starts_at, ends_at) WITH &&) WHERE (status = 'active');
