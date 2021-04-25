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

func (c compositeValidator) Types() []string {
	var m = make(map[string]interface{})

	for i := range c.validators {
		var t, ok = c.validators[i].(asset.Types)

		if !ok {
			continue
		}

		var p = t.Types()

		for i := range p {
			m[p[i]] = nil
		}
	}

	var (
		ts = make([]string, len(m))
		i  = 0
	)

	for k := range m {
		ts[i] = k
	}

	return ts
}

// CompositeValidator allows multiple validators to be executed on a kind
func CompositeValidator(v ...asset.Validator) asset.Validator {
	return compositeValidator{validators: v}
}
