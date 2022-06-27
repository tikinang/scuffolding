package model

import (
	"context"
	"github.com/pkg/errors"
	"icbaat/pkg/shared/skelet"
)

type Factory struct {
	orm *skelet.Orm
}

func NewFactory(
	runner *skelet.Runner,
	orm *skelet.Orm,
) (r *Factory) {
	defer func() { runner.Register(r) }()
	return &Factory{
		orm: orm,
	}
}

func (r *Factory) Before(ctx context.Context) error {
	if err := r.orm.GetDb().AutoMigrate(register()...); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (r *Factory) Begin(ctx context.Context) (*skelet.Tx, error) {
	return r.orm.Begin(ctx)
}

func register() []any {
	return []any{
		&Art{},
	}
}
