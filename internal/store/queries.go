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

func slotErr(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == "23P01" {
		return ErrSlotTaken
	}
	return err
}

func (s *Store) Bookings(ctx context.Context, from, to time.Time) ([]Booking, error) {
	rows, _ := s.db.Query(ctx, `SELECT id, service_id, barber_id, starts_at, ends_at, client_name, client_phone
		FROM bookings WHERE starts_at < $2 AND ends_at > $1 ORDER BY starts_at`, from, to)
	return pgx.CollectRows(rows, pgx.RowToStructByPos[Booking])
}
