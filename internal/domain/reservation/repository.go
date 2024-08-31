package reservation

import "context"

type Repository interface {
	Create(context.Context, Reservation) (ID string, err error)
	Get(ctx context.Context, ID string) (Reservation, error)
	List(ctx context.Context, roomID string) ([]Reservation, error)
	Delete(ctx context.Context, ID string) error
	Update(ctx context.Context, ID string, data Reservation) error
}
