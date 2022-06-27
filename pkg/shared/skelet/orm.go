package skelet

import (
	"context"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"sync"
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
) (r *Orm) {
	defer func() { runner.Register(r) }()
	return &Orm{
		config: config,
		log:    log,
	}
}

func (r *Orm) Before(ctx context.Context) error {

	cfg := &gorm.Config{
		SkipDefaultTransaction: true,
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
		Logger:                   nil, // FIXME: mpavlicek - fill own logger
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

// FIXME: mpavlicek - this could cause problems
func (r *Orm) GetDb() *gorm.DB {
	return r.gorm
}
