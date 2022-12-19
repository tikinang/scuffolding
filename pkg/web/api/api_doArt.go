package api

import (
	"context"
	"scuffolding/pkg/web/model"
)

type DoArtIn struct {
	Id   model.ArtId `json:"id" binding:"required"`
	Hash string      `json:"hash" binding:"required"`
}

type DoArtOut struct {
	OldHash *string `json:"oldHash"`
}

func (r *Handler) DoArt(ctx context.Context, in DoArtIn) (DoArtOut, error) {
	repo := r.repo.GetArt()

	entity, found, err := repo.GetById(in.Id)
	if err != nil {
		return DoArtOut{}, err
	}

	tx, err := r.repo.Begin(ctx)
	if err != nil {
		return DoArtOut{}, err
	}
	defer tx.Rollback()

	activeRepo := repo.Activate(tx)

	var oldHash *string
	if !found {
		entity.Id = in.Id
		entity.Hash = in.Hash
		if err := activeRepo.Create(entity); err != nil {
			return DoArtOut{}, err
		}
	} else {
		oldHash = pointer(entity.Hash)
		entity.Hash = in.Hash
		if err := activeRepo.Update(entity); err != nil {
			return DoArtOut{}, err
		}
	}

	if err := tx.Commit(); err != nil {
		return DoArtOut{}, err
	}

	return DoArtOut{
		OldHash: oldHash,
	}, nil
}

func pointer[T any](t T) *T {
	return &t
}
