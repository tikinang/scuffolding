package api

import (
	"icbaat/pkg/api/model"
)

type Handler struct {
	repo *model.Factory
}

func New(
	repositoryFactory *model.Factory,
) (r *Handler) {
	return &Handler{
		repo: repositoryFactory,
	}
}
