package repository

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"errors"
	"fmt"
	"room-reservation/internal/domain/reservation"
	"room-reservation/internal/repository/postgres"
	"strings"

	"github.com/jackc/pgx/v5"
)

type ReservationRepository struct {
	db *postgres.DB
}

func NewReservationRepository(ctx context.Context, connString string) (*ReservationRepository, error) {
	db, err := postgres.New(ctx, connString)
	if err != nil {
		return nil, err
	}

	repo := &ReservationRepository{
		db: db,
	}

	return repo, nil
}

func (r *ReservationRepository) Close() {
	if r.db != nil {
		r.db.Close()
	}
}

func (r *ReservationRepository) Create(ctx context.Context, data reservation.Reservation) (string, error) {
	tx, err := r.db.Begin(ctx)
	if err != nil {
		return "", err
	}
	defer tx.Rollback(ctx)

	var existingID string
	checkOverlapQuery := `
		SELECT id 
		FROM reservation 
		WHERE room_id = @roomID 
		AND start_time < @endTime
		AND end_time > @startTime
	`

	err = tx.QueryRow(ctx, checkOverlapQuery, pgx.NamedArgs{
		"roomID":    data.RoomID,
		"startTime": data.StartTime,
		"endTime":   data.EndTime,
	}).Scan(&existingID)
	if err != nil && err != pgx.ErrNoRows {
		return "", err
	}

	if existingID != "" {
		return "", reservation.ErrorOverlaps
	}

	insertQuery := `
		INSERT INTO reservation (id, room_id, start_time, end_time)
		VALUES ($1, $2, $3, $4)
	`
	data.ID = r.generateID()
	args := []any{data.ID, data.RoomID, data.StartTime, data.EndTime}

	_, err = tx.Exec(ctx, insertQuery, args...)
	if err != nil {
		return "", err
	}

	if err = tx.Commit(ctx); err != nil {
		return "", err
	}

	return data.ID, nil
}

func (r *ReservationRepository) Get(ctx context.Context, ID string) (reservation.Reservation, error) {
	q := `
		SELECT id, room_id, start_time, end_time
		FROM reservation
		WHERE id = $1
	`

	res := reservation.Reservation{}

	err := r.db.QueryRow(ctx, q, ID).Scan(&res.ID, &res.RoomID, &res.StartTime, &res.EndTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return reservation.Reservation{}, reservation.ErrorNotFound
		}

		return reservation.Reservation{}, err
	}

	return res, nil
}

func (r *ReservationRepository) List(ctx context.Context, roomID string) ([]reservation.Reservation, error) {
	q := `
		SELECT id, room_id, start_time, end_time
		FROM reservation
		WHERE room_id = $1
	`

	rows, err := r.db.Query(ctx, q, roomID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reservations := []reservation.Reservation{}
	for rows.Next() {
		var res reservation.Reservation

		err := rows.Scan(&res.ID, &res.RoomID, &res.StartTime, &res.EndTime)
		if err != nil {
			return nil, err
		}

		reservations = append(reservations, res)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	if len(reservations) == 0 {
		return nil, reservation.ErrorNotFoundForRoom
	}

	return reservations, nil
}

func (r *ReservationRepository) Delete(ctx context.Context, ID string) error {
	q := `
		DELETE FROM reservation
		WHERE id = $1
	`

	args := []any{ID}

	result, err := r.db.Exec(ctx, q, args...)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return reservation.ErrorNotFound
	}

	return nil
}

func (r *ReservationRepository) Update(ctx context.Context, ID string, data reservation.Reservation) error {
	sets, args := r.prepareArgs(data)
	if len(sets) > 0 {
		args = append(args, ID)
		q := fmt.Sprintf("UPDATE reservation SET %s WHERE id = $%d", strings.Join(sets, ", "), len(args))

		result, err := r.db.Exec(ctx, q, args...)
		if err != nil {
			return err
		}

		if result.RowsAffected() == 0 {
			return reservation.ErrorNotFound
		}
	}

	return nil
}

func (r *ReservationRepository) prepareArgs(data reservation.Reservation) (sets []string, args []any) {
	if data.RoomID != "" {
		args = append(args, data.RoomID)
		sets = append(sets, fmt.Sprintf("room_id=$%d", len(args)))
	}

	if !data.StartTime.IsZero() {
		args = append(args, data.StartTime)
		sets = append(sets, fmt.Sprintf("start_time=$%d", len(args)))
	}

	if !data.EndTime.IsZero() {
		args = append(args, data.EndTime)
		sets = append(sets, fmt.Sprintf("end_time=$%d", len(args)))
	}

	return
}

func (r *ReservationRepository) generateID() string {
	bytes := make([]byte, 6)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
