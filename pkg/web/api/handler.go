package api

import (
	"icbaat/pkg/web/model"
)

type Handler struct {
	repo *model.Factory
}

func New(
	repositoryFactory *model.Factory,
) *Handler {
	return &Handler{
		repo: repositoryFactory,
	}
}
