package reservation

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Request struct {
	RoomID    string   `json:"room_id" example:"1"`
	StartTime DateTime `json:"start_time" example:"29-08-2024 13:00" swaggertype:"primitive,string"`
	EndTime   DateTime `json:"end_time" example:"29-08-2024 14:00" swaggertype:"primitive,string"`
}

type DateTime struct {
	time.Time
}

const dateTimeLayout = "02-01-2006 15:04"

func (dt *DateTime) UnmarshalJSON(b []byte) (err error) {
	s := strings.Trim(string(b), "\"")
	if s == "null" || s == "" {
		dt.Time = time.Time{}
		return
	}

	dt.Time, err = time.Parse(dateTimeLayout, s)
	if err != nil {
		return fmt.Errorf("invalid time format %v", err)
	}
	return
}

func (dt DateTime) MarshalJSON() ([]byte, error) {
	formatted := dt.Format(dateTimeLayout)
	return []byte(`"` + formatted + `"`), nil
}

func (r *Request) Validate() error {
	if r.RoomID == "" {
		return errors.New("room_id is required")
	}

	if r.StartTime.IsZero() {
		return errors.New("invalid start_time")
	}

	if r.EndTime.IsZero() {
		return errors.New("invalid end_time")
	}

	if r.StartTime.After(r.EndTime.Time) {
		return errors.New("start_time must be before end_time")
	}

	return nil
}

type UpdateRequest struct {
	RoomID    string   `json:"room_id" example:"1"`
	StartTime DateTime `json:"start_time" example:"29-08-2024 13:00" swaggertype:"primitive,string"`
	EndTime   DateTime `json:"end_time" example:"29-08-2024 14:00" swaggertype:"primitive,string"`
}

func (r *UpdateRequest) Validate() error {
	if r.RoomID == "" && r.StartTime.IsZero() && r.EndTime.IsZero() {
		return errors.New("no fields to update")
	}
	return nil
}

type Response struct {
	ID        string   `json:"id"`
	RoomID    string   `json:"room_id"`
	StartTime DateTime `json:"start_time"`
	EndTime   DateTime `json:"end_time"`
}

func ToResponse(data Reservation) Response {
	return Response{
		ID:        data.ID,
		RoomID:    data.RoomID,
		StartTime: DateTime{data.StartTime},
		EndTime:   DateTime{data.EndTime},
	}
}

func ToResponseSlice(data []Reservation) []Response {
	res := make([]Response, 0)

	for _, r := range data {
		res = append(res, ToResponse(r))
	}

	return res
}
