package skelet

import (
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Repository[Id any, Entity any] struct {
	DbProvider
}

func NewRepository[Id any, Entity any](
	db DbProvider,
) *Repository[Id, Entity] {
	return &Repository[Id, Entity]{
		DbProvider: db,
	}
}

func (r *Repository[Id, Entity]) Activate(tx *Tx) *ActivatedRepository[Id, Entity] {
	return &ActivatedRepository[Id, Entity]{
		Repository: NewRepository[Id, Entity](tx),
	}
}

func (r *Repository[Id, Entity]) GetById(id Id) (Entity, bool, error) {
	var entity Entity
	if err := r.Db().Where(id).Take(&entity).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return entity, false, nil
		}
		return entity, false, errors.WithStack(err)
	}
	return entity, true, nil
}

type ActivatedRepository[Id any, Entity any] struct {
	*Repository[Id, Entity]
}

func (r *ActivatedRepository[Id, Entity]) Create(entity Entity) error {
	return r.Db().Create(entity).Error
}

func (r *ActivatedRepository[Id, Entity]) Update(entity Entity) error {
	return r.Db().Save(entity).Error
}

func (r *ActivatedRepository[Id, Entity]) Delete(entity Entity) error {
	return r.Db().Delete(entity).Error
}
