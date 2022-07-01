package model

import (
	"gorm.io/gorm/schema"
)

type ArtistId string

type Artist struct {
	Id   ArtistId `gorm:"primarykey"`
	Name string
}

func (r *Artist) TableName() string {
	return "artist"
}

var _ schema.Tabler = (*Artist)(nil)
