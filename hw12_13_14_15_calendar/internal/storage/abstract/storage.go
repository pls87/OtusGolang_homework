package abstractstorage

import "context"

type Storage interface {
	Events() EventRepository
	Init(ctx context.Context) error
	Destroy() error
}
