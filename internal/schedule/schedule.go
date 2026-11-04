package schedule

import "time"

const (
	OpenHour  = 10
	CloseHour = 21
	Step      = 30 * time.Minute
)

type Interval struct {
	Start, End time.Time
}

func (i Interval) Overlaps(o Interval) bool {
	return i.Start.Before(o.End) && o.Start.Before(i.End)
}

func Slots(day time.Time) []time.Time {
	y, m, d := day.Date()
	t := time.Date(y, m, d, OpenHour, 0, 0, 0, day.Location())
	end := time.Date(y, m, d, CloseHour, 0, 0, 0, day.Location())
	var slots []time.Time
	for ; t.Before(end); t = t.Add(Step) {
		slots = append(slots, t)
	}
	return slots
}

func Free(day time.Time, duration time.Duration, busy []Interval, now time.Time) []time.Time {
	y, m, d := day.Date()
	end := time.Date(y, m, d, CloseHour, 0, 0, 0, day.Location())
	var free []time.Time
	for _, t := range Slots(day) {
		slot := Interval{t, t.Add(duration)}
		if t.Before(now) || slot.End.After(end) || overlapsAny(slot, busy) {
			continue
		}
		free = append(free, t)
	}
	return free
}

func overlapsAny(slot Interval, busy []Interval) bool {
	for _, b := range busy {
		if slot.Overlaps(b) {
			return true
		}
	}
	return false
}
