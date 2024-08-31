package reservation

func (r *Reservation) Overlaps(other Reservation) bool {
	return r.RoomID == other.RoomID && r.StartTime.Before(other.EndTime) && r.EndTime.After(other.StartTime)
}
