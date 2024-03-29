package skelet

import (
	"context"
	"reflect"
	"sync"

	"github.com/pkg/errors"
)

type Runner struct {
	log        *Logger
	components []any
}

func NewRunner(
	log *Logger,
) *Runner {
	return &Runner{
		log: log,
	}
}

type Before interface {
	Before(ctx context.Context) error
}

// TODO(mpavlicek): just return error, no need to use chan; cancel context on error happening
type Run interface {
	Run(ctx context.Context, done chan<- error)
}

type After interface {
	After() error
}

func (r *Runner) Register(component any) {
	r.log.WithField("component", reflect.TypeOf(component)).Debug("registering component")
	r.components = append(r.components, component)
}

func (r *Runner) RunBefore(ctx context.Context) error {
	for _, component := range r.components {
		before, is := component.(Before)
		if !is {
			continue
		}

		log := r.log.WithField("component", reflect.TypeOf(before))
		log.Debug("before")

		if err := before.Before(ctx); err != nil {
			log.WithError(err).Error("before error")
			return errors.WithStack(err)
		}
	}
	return nil
}

func (r *Runner) RunAfter() {
	for i := len(r.components); i > 0; i-- {
		after, is := r.components[i-1].(After)
		if !is {
			continue
		}

		log := r.log.WithField("component", reflect.TypeOf(after))
		log.Debug("before")

		if err := after.After(); err != nil {
			log.WithError(err).Error("after error")
		}
	}
}

func (r *Runner) Run(ctx context.Context) error {
	if err := r.RunBefore(ctx); err != nil {
		return err
	}

	var wg sync.WaitGroup
	for _, component := range r.components {
		run, is := component.(Run)
		if !is {
			continue
		}
		wg.Add(1)
		go func(run Run) {
			defer wg.Done()

			log := r.log.WithField("component", reflect.TypeOf(run))
			log.Debug("run")

			done := make(chan error)
			go run.Run(ctx, done)
			if err := <-done; err != nil {
				log.WithError(err).Error("run error")
			}
		}(run)
	}
	wg.Wait()

	r.RunAfter()

	return nil
}
