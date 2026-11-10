package store

import "time"

type Service struct {
	ID          int64
	Name        string
	Price       int
	DurationMin int
}

func (s Service) Duration() time.Duration {
	return time.Duration(s.DurationMin) * time.Minute
}

type Barber struct {
	ID   int64
	Name string
}

type Booking struct {
	ID          int64
	ServiceID   int64
	BarberID    int64
	StartsAt    time.Time
	EndsAt      time.Time
	ClientName  string
	ClientPhone string
	ServiceName string
	BarberName  string
}
