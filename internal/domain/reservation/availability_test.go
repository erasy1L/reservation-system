package reservation

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestAvailability(t *testing.T) {
	tests := map[string]func(t *testing.T){
		"Overlap":         testOverlap,
		"No Overlap":      testNoOverlap,
		"Exact End-Start": testExactEndStart,
	}

	for name, test := range tests {
		t.Run(name, test)
	}
}

func testOverlap(t *testing.T) {
	room1 := Reservation{
		RoomID:    "1",
		StartTime: time.Date(2024, 8, 29, 13, 0, 0, 0, time.Local),
		EndTime:   time.Date(2024, 8, 29, 14, 0, 0, 0, time.Local),
	}
	room2 := Reservation{
		RoomID:    "1",
		StartTime: time.Date(2024, 8, 29, 13, 30, 0, 0, time.Local),
		EndTime:   time.Date(2024, 8, 29, 14, 30, 0, 0, time.Local),
	}

	ok := room1.Overlaps(room2)
	assert.Truef(t, ok, "expected overlap between reservations %v and %v", room1, room2)
}

func testNoOverlap(t *testing.T) {
	room1 := Reservation{
		RoomID:    "1",
		StartTime: time.Date(2024, 8, 29, 13, 0, 0, 0, time.Local),
		EndTime:   time.Date(2024, 8, 29, 14, 0, 0, 0, time.Local),
	}
	room2 := Reservation{
		RoomID:    "1",
		StartTime: time.Date(2024, 8, 29, 14, 0, 0, 0, time.Local),
		EndTime:   time.Date(2024, 8, 29, 15, 0, 0, 0, time.Local),
	}

	ok := room1.Overlaps(room2)
	assert.Falsef(t, ok, "expected no overlap between reservations %v and %v", room1, room2)
}

func testExactEndStart(t *testing.T) {
	room1 := Reservation{
		RoomID:    "1",
		StartTime: time.Date(2024, 8, 29, 13, 0, 0, 0, time.Local),
		EndTime:   time.Date(2024, 8, 29, 14, 0, 0, 0, time.Local),
	}
	room2 := Reservation{
		RoomID:    "1",
		StartTime: time.Date(2024, 8, 29, 14, 0, 0, 0, time.Local),
		EndTime:   time.Date(2024, 8, 29, 15, 0, 0, 0, time.Local),
	}

	ok := room1.Overlaps(room2)
	assert.Falsef(t, ok, "expected no overlap between reservations %v and %v", room1, room2)
}
