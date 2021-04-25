package validator

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/fluxynet/gocipe/asset"
)

const (
	MimeAll = "*"
)

var (
	// ErrMaxSizeExceeded means the size of the file is bigger than allowed
	ErrMaxSizeExceeded = errors.New("maximum file size exceeded")

	// ErrInvalidMimeType means the mimetype is not in the allowed list
	ErrInvalidMimeType = errors.New("invalid mime type")
)

// Basic Validator meant to have common validation rules. Can be used for composing more specific validators
type Basic struct {
	MaxSize      int
	AllowedMimes []string
}

func (v Basic) Validate(ctx context.Context, args *asset.ValidateArgs) error {
	if s := len(args.Data); v.MaxSize > 0 && s >= v.MaxSize {
		return fmt.Errorf(
			"%w. filesize = %d",
			ErrMaxSizeExceeded,
			s,
		)
	}

	for i := 0; i < len(v.AllowedMimes); i++ {
		switch v.AllowedMimes[i] {
		case args.Asset.Mime, MimeAll:
			return nil
		}
	}

	return fmt.Errorf(
		"%w. obtained: [%s] allowed types: %s",
		ErrInvalidMimeType,
		args.Asset.Mime,
		strings.Join(v.AllowedMimes, ", "),
	)
}

func (v Basic) Types() []string {
	return v.AllowedMimes
}
