package disk

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fluxynet/gocipe/storage"
	"github.com/fluxynet/gocipe/tenant"
)

// Must returns a disk storage or
func Must(path string) storage.Storage {
	var s, err = Disk(path)
	if err != nil {
		panic("failed to initialize Disk storage: " + err.Error())
	}

	return s
}

// Disk storage mechanism
func Disk(path string) (storage.Storage, error) {
	var err error
	path, err = filepath.Abs(path)

	if err == nil {
		err = os.MkdirAll(path, 0744)
	}

	return disk{
		Path: path,
	}, err
}

type disk struct {
	Path string
}

func (d disk) Store(ctx context.Context, filename string, contntts []byte) error {
	var (
		tnt = tenant.Get(ctx)
		dir = filepath.Join(d.Path, tnt, filepath.Dir(filename))
		err = os.MkdirAll(dir, 0744)
	)

	if err == nil {
		err = os.WriteFile(filename, contntts, 0644)
	}

	return err
}

func (d disk) Delete(ctx context.Context, path string) error {
	var (
		tnt    = tenant.Get(ctx)
		loc    = filepath.Join(d.Path, tnt, path)
		_, err = os.Stat(loc)
	)

	if os.IsNotExist(err) {
		return nil
	}

	return os.RemoveAll(loc)
}
