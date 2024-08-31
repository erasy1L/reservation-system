package reservation

import (
	"errors"
	"time"
)

type Reservation struct {
	ID        string    `db:"id"`
	RoomID    string    `db:"room_id"`
	StartTime time.Time `db:"start_time"`
	EndTime   time.Time `db:"end_time"`
}

var ErrorNotFound error = errors.New("reservation not found")
var ErrorNotFoundForRoom error = errors.New("reservations not found for room")
var ErrorOverlaps error = errors.New("reservation overlaps with another")
