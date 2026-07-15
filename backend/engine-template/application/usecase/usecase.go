package usecase

import (
	"context"
	"github.com/sekaishopml/cytaxi/backend/engine-template/domain/entity"
	"github.com/sekaishopml/cytaxi/backend/engine-template/domain/repository"
)

type UseCase struct {
	repo repository.Repository
}

func New(repo repository.Repository) *UseCase {
	return &UseCase{repo: repo}
}

func (uc *UseCase) Execute(ctx context.Context, id string) (*entity.Entity, error) {
	return uc.repo.FindByID(ctx, id)
}
