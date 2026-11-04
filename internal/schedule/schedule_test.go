package schedule

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var day = time.Date(2026, 11, 10, 0, 0, 0, 0, time.UTC)

func at(h, m int) time.Time {
	return time.Date(2026, 11, 10, h, m, 0, 0, time.UTC)
}

func TestSlots(t *testing.T) {
	slots := Slots(day)
	assert.Len(t, slots, 22)
	assert.Equal(t, at(10, 0), slots[0])
	assert.Equal(t, at(20, 30), slots[len(slots)-1])
}

func TestFreeSkipsBusy(t *testing.T) {
	busy := []Interval{{at(11, 0), at(12, 0)}}
	free := Free(day, time.Hour, busy, at(0, 0))
	assert.Contains(t, free, at(10, 0))
	assert.NotContains(t, free, at(10, 30))
	assert.NotContains(t, free, at(11, 0))
	assert.NotContains(t, free, at(11, 30))
	assert.Contains(t, free, at(12, 0))
}

func TestFreeFitsBeforeClose(t *testing.T) {
	free := Free(day, time.Hour, nil, at(0, 0))
	assert.Equal(t, at(20, 0), free[len(free)-1])
}

func TestFreeSkipsPast(t *testing.T) {
	free := Free(day, 30*time.Minute, nil, at(15, 10))
	assert.Equal(t, at(15, 30), free[0])
}
