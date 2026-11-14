package store

import (
	"context"
	"errors"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

var ErrSlotTaken = errors.New("slot taken")

func (s *Store) Services(ctx context.Context) ([]Service, error) {
	rows, _ := s.db.Query(ctx, "SELECT id, name, price, duration_min FROM services ORDER BY id")
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Service])
}

func (s *Store) Service(ctx context.Context, id int64) (Service, error) {
	rows, _ := s.db.Query(ctx, "SELECT id, name, price, duration_min FROM services WHERE id = $1", id)
	return pgx.CollectExactlyOneRow(rows, pgx.RowToStructByPos[Service])
}

func (s *Store) Barbers(ctx context.Context) ([]Barber, error) {
	rows, _ := s.db.Query(ctx, "SELECT id, name FROM barbers ORDER BY id")
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Barber])
}

func (s *Store) CreateBooking(ctx context.Context, b Booking) error {
	_, err := s.db.Exec(ctx, `INSERT INTO bookings (service_id, barber_id, starts_at, ends_at, client_name, client_phone)
		VALUES ($1, $2, $3, $4, $5, $6)`,
		b.ServiceID, b.BarberID, b.StartsAt, b.EndsAt, b.ClientName, b.ClientPhone)
	return slotErr(err)
}

func (s *Store) CancelBooking(ctx context.Context, id int64) error {
	_, err := s.db.Exec(ctx, "UPDATE bookings SET status = 'cancelled' WHERE id = $1", id)
	return err
}

func (s *Store) RescheduleBooking(ctx context.Context, id, barberID int64, start time.Time) error {
	_, err := s.db.Exec(ctx, `UPDATE bookings SET barber_id = $2, ends_at = $3::timestamptz + (ends_at - starts_at), starts_at = $3
		WHERE id = $1`, id, barberID, start)
	return slotErr(err)
}

func slotErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
		return ErrSlotTaken
	}
	return err
}

func (s *Store) Bookings(ctx context.Context, from, to time.Time) ([]Booking, error) {
	rows, _ := s.db.Query(ctx, `SELECT b.id, b.service_id, b.barber_id, b.starts_at, b.ends_at, b.client_name, b.client_phone, s.name, br.name
		FROM bookings b
		JOIN services s ON s.id = b.service_id
		JOIN barbers br ON br.id = b.barber_id
		WHERE b.status = 'active' AND b.starts_at < $2 AND b.ends_at > $1
		ORDER BY b.starts_at`, from, to)
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Booking])
}

func (s *Store) SaveAdmin(ctx context.Context, login, hash string) error {
	_, err := s.db.Exec(ctx, `INSERT INTO admins (login, password_hash) VALUES ($1, $2)
		ON CONFLICT (login) DO UPDATE SET password_hash = EXCLUDED.password_hash`, login, hash)
	return err
}

func (s *Store) AdminHash(ctx context.Context, login string) (string, error) {
	var hash string
	err := s.db.QueryRow(ctx, "SELECT password_hash FROM admins WHERE login = $1", login).Scan(&hash)
	return hash, err
}

func (s *Store) Blocks(ctx context.Context, from, to time.Time) ([]Block, error) {
	rows, _ := s.db.Query(ctx, `SELECT id, barber_id, starts_at, ends_at, reason FROM blocks
		WHERE starts_at < $2 AND ends_at > $1 ORDER BY starts_at`, from, to)
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Block])
}

func (s *Store) CreateBlock(ctx context.Context, b Block) error {
	_, err := s.db.Exec(ctx, "INSERT INTO blocks (barber_id, starts_at, ends_at, reason) VALUES ($1, $2, $3, $4)",
		b.BarberID, b.StartsAt, b.EndsAt, b.Reason)
	return err
}

func (s *Store) DeleteBlock(ctx context.Context, id int64) error {
	_, err := s.db.Exec(ctx, "DELETE FROM blocks WHERE id = $1", id)
	return err
}

func (s *Store) BookingsPerWeek(ctx context.Context) ([]Stat, error) {
	rows, _ := s.db.Query(ctx, `SELECT to_char(date_trunc('week', starts_at), 'DD.MM'), count(*)::int
		FROM bookings
		WHERE status = 'active' AND starts_at >= date_trunc('week', now()) - interval '7 weeks'
		GROUP BY date_trunc('week', starts_at)
		ORDER BY date_trunc('week', starts_at)`)
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Stat])
}

func (s *Store) TopServices(ctx context.Context) ([]Stat, error) {
	rows, _ := s.db.Query(ctx, `SELECT s.name, count(*)::int
		FROM bookings b JOIN services s ON s.id = b.service_id
		WHERE b.status = 'active'
		GROUP BY s.name
		ORDER BY 2 DESC
		LIMIT 5`)
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Stat])
}
