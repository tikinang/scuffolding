package model

import (
	"gorm.io/gorm/schema"
	"icbaat/pkg/shared/skelet"
)

func (r *Factory) NewArtRepository() *skelet.Repository[ArtId, Art] {
	return skelet.NewRepository[ArtId, Art](r.orm.GetDb())
}

type ArtId string

type Art struct {
	Id   ArtId `gorm:"primarykey"`
	Hash string
}

func (r *Art) TableName() string {
	return "art"
}

var _ schema.Tabler = (*Art)(nil)
