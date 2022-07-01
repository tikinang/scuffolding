package model

import (
	"context"
	"github.com/pkg/errors"
	"icbaat/pkg/shared/skelet"
)

type Register struct {
	Art    *skelet.Repository[ArtId, Art]
	Artist *skelet.Repository[ArtistId, Artist]
}

type Factory struct {
	orm *skelet.Orm

	register *Register
}

func NewFactory(
	runner *skelet.Runner,
	orm *skelet.Orm,
) (r *Factory) {
	defer func() { runner.Register(r) }()
	return &Factory{
		orm: orm,
		register: &Register{
			Art:    skelet.NewRepository[ArtId, Art](orm),
			Artist: skelet.NewRepository[ArtistId, Artist](orm),
		},
	}
}

func (r *Factory) Before(ctx context.Context) error {
	if err := r.orm.Db().AutoMigrate(register()...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *Factory) Begin(ctx context.Context) (*skelet.Tx, error) {
	return r.orm.Begin(ctx)
}

func (r *Factory) GetArt() *skelet.Repository[ArtId, Art] {
	return r.register.Art
}

func (r *Factory) GetArtist() *skelet.Repository[ArtistId, Artist] {
	return r.register.Artist
}

func register() []any {
	return []any{
		&Art{},
		&Artist{},
	}
}
