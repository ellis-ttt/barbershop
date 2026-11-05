package store

import (
	"context"
	"io/fs"

	"barbershop/migrations"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Store struct {
	db *pgxpool.Pool
}

func New(ctx context.Context, dsn string) (*Store, error) {
	db, err := pgxpool.New(ctx, dsn)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(ctx); err != nil {
		return nil, err
	}
	return &Store{db: db}, nil
}

func (s *Store) Close() {
	s.db.Close()
}

func (s *Store) Migrate(ctx context.Context) error {
	files, err := fs.Glob(migrations.FS, "*.sql")
	if err != nil {
		return err
	}
	for _, f := range files {
		sql, err := migrations.FS.ReadFile(f)
		if err != nil {
			return err
		}
		if _, err := s.db.Exec(ctx, string(sql)); err != nil {
			return err
		}
	}
	return nil
}
