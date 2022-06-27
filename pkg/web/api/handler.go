package api

import (
	"icbaat/pkg/web/model"
)

type Handler struct {
	repositoryFactory *model.Factory
}

func New(
	repositoryFactory *model.Factory,
) (r *Handler) {
	return &Handler{
		repositoryFactory: repositoryFactory,
	}
}
