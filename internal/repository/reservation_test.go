package repository

import (
	"context"
	"room-reservation/internal/domain/reservation"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestReservationRepository(t *testing.T) {
	t.Helper()

	ctx := context.Background()

	repo := &ReservationRepository{
		db: db,
	}

	t.Run("Create reservation", func(t *testing.T) {
		testCreateReservation(ctx, repo, t)
	})

	t.Run("Create reservation with overlapping", func(t *testing.T) {
		testCreateReservationShouldOverlap(ctx, repo, t)
	})

	t.Run("Create reservation with no overlapping", func(t *testing.T) {
		testCreateReservationShouldNotOverlap(ctx, repo, t)
	})

	t.Run("Get reservation", func(t *testing.T) {
		testGetReservation(ctx, repo, t)
	})

	t.Run("List reservations", func(t *testing.T) {
		testListReservation(ctx, repo, t)
	})

	t.Run("Update reservation", func(t *testing.T) {
		testUpdateReservation(ctx, repo, t)
	})

	t.Run("Delete reservation", func(t *testing.T) {
		testDeleteReservation(ctx, repo, t)
	})
}

var testData = reservation.Reservation{
	RoomID:    "1",
	StartTime: time.Date(2024, 8, 29, 13, 0, 0, 0, time.UTC),
	EndTime:   time.Date(2024, 8, 29, 14, 0, 0, 0, time.UTC),
}

func testCreateReservation(ctx context.Context, repo *ReservationRepository, t *testing.T) {
	ID, err := repo.Create(ctx, testData)
	require.NoError(t, err, "could not create reservation")

	require.NotEmpty(t, ID, "expected a non-empty ID")

	testData.ID = ID // for other tests
}

func testCreateReservationShouldOverlap(ctx context.Context, repo *ReservationRepository, t *testing.T) {
	overlapping := reservation.Reservation{
		RoomID:    testData.RoomID,
		StartTime: testData.StartTime.Add(30 * time.Minute),
		EndTime:   testData.EndTime,
	}

	_, err := repo.Create(ctx, overlapping)
	require.ErrorIs(t, err, reservation.ErrorOverlaps, "expected overlap")
}

func testCreateReservationShouldNotOverlap(ctx context.Context, repo *ReservationRepository, t *testing.T) {
	nonOverlapping := reservation.Reservation{
		RoomID:    testData.RoomID,
		StartTime: testData.EndTime,
		EndTime:   testData.EndTime.Add(30 * time.Minute),
	}

	_, err := repo.Create(ctx, nonOverlapping)
	require.NoError(t, err, "expected no overlap")
}

func testGetReservation(ctx context.Context, repo *ReservationRepository, t *testing.T) {
	res, err := repo.Get(ctx, testData.ID)
	require.NoError(t, err, "failed to get reservation")

	require.Equalf(t, testData, res, "expected %v, got %v", testData, res)
}

func testListReservation(ctx context.Context, repo *ReservationRepository, t *testing.T) {
	reservations, err := repo.List(ctx, testData.RoomID)
	require.NoError(t, err, "failed to list reservations")

	require.NotEmpty(t, reservations, "expected at least one reservation")

	require.Equalf(t, testData.RoomID, reservations[0].RoomID, "expected room ID %s, got %s", testData.RoomID, reservations[0].RoomID)
}

func testUpdateReservation(ctx context.Context, repo *ReservationRepository, t *testing.T) {
	updatedEndTime := testData.EndTime.Add(time.Hour)
	toUpdate := reservation.Reservation{
		EndTime: updatedEndTime,
	}

	err := repo.Update(ctx, testData.ID, toUpdate)
	require.NoError(t, err, "failed to update reservation")

	updated, err := repo.Get(ctx, testData.ID)
	require.NoError(t, err, "failed to get updated reservation")

	require.Equalf(t, updatedEndTime, updated.EndTime, "expected end time %v, got %v", updatedEndTime, updated.EndTime)
}

func testDeleteReservation(ctx context.Context, repo *ReservationRepository, t *testing.T) {
	err := repo.Delete(ctx, testData.ID)
	require.NoError(t, err, "failed to delete reservation")

	_, err = repo.Get(ctx, testData.ID)
	require.Error(t, err, "expected an error when getting deleted reservation")
}
