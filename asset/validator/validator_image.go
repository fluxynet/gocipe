package validator

import (
	"context"

	"github.com/fluxynet/gocipe/asset"
)

// ImageOpts represents options for the image validator
type ImageOpts struct {
	MinWidth  int
	MaxWidth  int
	MinHeight int
	MaxHeight int
}

type image struct {
	Basic
	ImageOpts
}

func (i image) Validate(ctx context.Context, args *asset.ValidateArgs) error {
	var err = i.Basic.Validate(ctx, args)
	if err != nil {
		return err
	}

	// todo: validate dimensions

	return nil
}

// Image validator accepting types jpeg, png and webp
func Image(opts ImageOpts, maxSize int) asset.Validator {
	return image{
		Basic: Basic{
			MaxSize: maxSize,
			AllowedMimes: []string{
				"image/jpeg",
				"image/png",
				"image/webp",
			},
		},
		ImageOpts: opts,
	}
}
