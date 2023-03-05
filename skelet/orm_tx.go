package skelet

import (
	"sync"

	"gorm.io/gorm"
)

type Tx struct {
	log  *Logger
	db   *gorm.DB
	once *sync.Once
}

func (r *Tx) Commit() error {
	var err error
	r.once.Do(func() {
		err = r.db.Commit().Error
	})
	return err
}

func (r *Tx) Db() *gorm.DB {
	return r.db
}

func (r *Tx) Rollback() {
	r.once.Do(func() {
		if err := r.db.Rollback().Error; err != nil {
			r.log.WithError(err).Error("rollback error")
		}
	})
}
