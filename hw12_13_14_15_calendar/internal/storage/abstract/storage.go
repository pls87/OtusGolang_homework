package abstractstorage

import "context"

type Storage interface {
	Events() EventRepository
	Connect(ctx context.Context) error
	Close() error
}
