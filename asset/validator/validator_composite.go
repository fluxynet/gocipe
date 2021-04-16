package validator

import (
	"context"

	"github.com/fluxynet/gocipe/asset"
)

type compositeValidator struct {
	validators []asset.Validator
}

func (c compositeValidator) Validate(ctx context.Context, args *asset.ValidateArgs) error {
	var err error

	for i := range c.validators {
		err = c.validators[i].Validate(ctx, args)
		if err != nil {
			return err
		}
	}

	return nil
}

// CompositeValidator allows multiple validators to be executed on a kind
func CompositeValidator(v ...asset.Validator) asset.Validator {
	return compositeValidator{validators: v}
}
