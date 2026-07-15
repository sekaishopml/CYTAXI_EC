package repository

import "context"

type Repository interface {
	FindByID(ctx context.Context, id string) (any, error)
	Save(ctx context.Context, entity any) error
}
