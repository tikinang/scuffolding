package skelet

import (
	"context"
	"sync"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type OrmConfig struct {
}

func DefaultOrmConfig() OrmConfig {
	return OrmConfig{}
}

type Orm struct {
	config OrmConfig
	log    *Logger

	gorm *gorm.DB
}

func NewOrm(
	runner *Runner,
	config OrmConfig,
	log *Logger,
) *Orm {
	r := &Orm{
		config: config,
		log:    log,
	}
	runner.Register(r)
	return r
}

func (r *Orm) Before(ctx context.Context) error {
	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger:                   nil, // TODO(mpavlicek): own logger
		DisableNestedTransaction: true,
	}

	var err error
	r.gorm, err = gorm.Open(sqlite.Open("gorm.db"), cfg)
	if err != nil {
		return errors.Wrap(err, "open db")
	}

	return nil
}

func (r *Orm) Begin(ctx context.Context) (*Tx, error) {
	db := r.gorm.WithContext(ctx).Begin()
	if err := r.gorm.Error; err != nil {
		return nil, errors.WithStack(err)
	}
	return &Tx{
		log:  r.log,
		db:   db,
		once: &sync.Once{},
	}, nil
}

type DbProvider interface {
	Db() *gorm.DB
}

// TODO(mpavlicek): this could cause problems
func (r *Orm) Db() *gorm.DB {
	return r.gorm
}
