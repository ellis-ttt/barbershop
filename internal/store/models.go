package store

import "time"

type Service struct {
	ID          int64
	Name        string
	Price       int
	DurationMin int
	Description string
}

func (s Service) Duration() time.Duration {
	return time.Duration(s.DurationMin) * time.Minute
}

type Barber struct {
	ID    int64
	Name  string
	Title string
	Photo string
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

type Block struct {
	ID       int64
	BarberID int64
	StartsAt time.Time
	EndsAt   time.Time
	Reason   string
}

type Stat struct {
	Label string
	Count int
}
