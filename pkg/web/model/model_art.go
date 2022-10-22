package model

import (
	"gorm.io/gorm/schema"
)

type ArtId string

type Art struct {
	Id       ArtId `gorm:"primarykey"`
	ArtistId ArtistId
	Artist   Artist `gorm:"foreignkey:ArtistId"`
	Hash     string
}

func (r *Art) TableName() string {
	return "art"
}

var _ schema.Tabler = (*Art)(nil)
