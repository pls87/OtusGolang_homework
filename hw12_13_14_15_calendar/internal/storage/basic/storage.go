package basic

import (
	"context"
	"errors"
)

var ErrDoesNotExist = errors.New("entity doesn't exist")

type Storage interface {
	Events() EventRepository
	Init(ctx context.Context) error
	Dispose() error
}
