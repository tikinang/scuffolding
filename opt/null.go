package opt

import (
	"encoding/json"

	"github.com/pkg/errors"
)

type Null[T any] struct {
	filled bool
	value  T
}

func (r Null[T]) Some() T {
	return r.value
}

func (r Null[T]) Filled() bool {
	return r.filled
}

func (r Null[T]) Ask() (T, bool) {
	return r.value, r.filled
}

func New[T any](value T) Null[T] {
	return Null[T]{
		filled: true,
		value:  value,
	}
}

func Empty[T any]() Null[T] {
	return Null[T]{}
}

var _ json.Marshaler = (*Null[string])(nil)

func (r Null[T]) MarshalJSON() ([]byte, error) {
	if r.filled {
		b, err := json.Marshal(r.value)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return b, nil
	} else {
		return []byte("null"), nil
	}
}

var _ json.Unmarshaler = (*Null[string])(nil)

func (r *Null[T]) UnmarshalJSON(b []byte) error {
	if len(b) == 0 || string(b) == "null" {
		return nil
	}
	r.filled = true
	err := json.Unmarshal(b, &r.value)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
