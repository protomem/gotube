package access

import (
	"context"
	"errors"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
	mongodbadapter "github.com/casbin/mongodb-adapter/v3"
	"github.com/protomem/gotube/pkg/logging"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var _ Manager = (*Casbin)(nil)

type CasbinOptions struct {
	MongoURI  string
	ModelPath string
	Lazy      bool
}

type Casbin struct {
	logger logging.Logger

	opts CasbinOptions

	once     sync.Once
	enforcer *casbin.Enforcer
}

func NewCasbin(_ context.Context, logger logging.Logger, opts CasbinOptions) (*Casbin, error) {
	c := &Casbin{
		logger: logger.With("system", "accessManager", "systemType", "casbin"),
		opts:   opts,
	}

	if !opts.Lazy {
		_, err := c.lazyEnforcer()
		if err != nil {
			return nil, fmt.Errorf("access:Casbin.New: %w", err)
		}
	}

	return c, nil
}

func (c *Casbin) Close(_ context.Context) error {
	return nil
}

func (c *Casbin) lazyEnforcer() (*casbin.Enforcer, error) {
	const op = "lazyEnforcer"
	var err error

	c.once.Do(func() {
		opts := options.Client().ApplyURI(c.opts.MongoURI)

		var a persist.BatchAdapter
		a, err = mongodbadapter.NewAdapterWithClientOption(opts, "gotube_casbin")
		if err != nil {
			err = fmt.Errorf("mongo adapter: %w", err)
			return
		}

		var m model.Model
		m, err = model.NewModelFromFile(c.opts.ModelPath)
		if err != nil {
			err = fmt.Errorf("model: %w", err)
			return
		}

		var e *casbin.Enforcer
		e, err = casbin.NewEnforcer(m, a)
		if err != nil {
			err = fmt.Errorf("enforcer: %w", err)
			return
		}

		c.enforcer = e
	})

	if err != nil {
		return nil, fmt.Errorf("%s: %w", op, err)
	}

	if c.enforcer == nil {
		return nil, fmt.Errorf("%s: %w", op, errors.New("enforcer is nil"))
	}

	return c.enforcer, nil
}
