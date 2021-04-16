package validator

import (
	"context"
	"errors"
	"fmt"

	"github.com/fluxynet/gocipe/asset"
	"github.com/fluxynet/gocipe/repository"
	"github.com/fluxynet/gocipe/types/fields/entity"
)

var (
	// ErrResourceNotKnown when an unknown resource is passed as kind
	ErrResourceNotKnown = errors.New("resource not known")
)

type resource struct {
	repo  repository.Repositorium
	names []repository.Named
}

func (r *resource) Validate(ctx context.Context, args *asset.ValidateArgs) error {
	var found bool

	for i := range r.names {
		if args.Asset.Resource == r.names[i].Name() {
			found = true
			break
		}
	}

	if !found {
		return fmt.Errorf("%w: %s", ErrResourceNotKnown, args.Asset.Resource)
	}

	var _, err = r.repo.Get(ctx, entity.ID(args.Asset.Resource), args.Asset.ResourceID)
	return err
}

// Resource validator checks if a resource name is valid and the id exists
func Resource(repo repository.Repositorium, names ...repository.Named) asset.Validator {
	return &resource{
		repo:  repo,
		names: names,
	}
}
