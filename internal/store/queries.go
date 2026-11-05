package store

import (
	"context"

	"github.com/jackc/pgx/v5"
)

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
	return err
}
