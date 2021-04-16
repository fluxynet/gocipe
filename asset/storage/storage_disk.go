package storage

import (
	"context"
	"os"
	"path/filepath"

	"github.com/fluxynet/gocipe/asset"
)

func MustDisk(path string) asset.Storage {
	var s, err = Disk(path)
	if err != nil {
		panic("failed to initialize Disk asset storage: " + err.Error())
	}

	return s
}

// Disk storage mechanism
func Disk(path string) (asset.Storage, error) {
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
	BaseURL string
	Path    string
}

func (d disk) Store(ctx context.Context, args asset.StoreArgs) (asset.Asset, error) {
	var (
		dir = filepath.Join(d.Path, args.Partition, args.Asset.ID)
		err = os.MkdirAll(dir, 0744)
	)

	if err == nil {
		var name = filepath.Join(dir, args.Asset.Name)
		err = os.WriteFile(name, args.Data, 0644)
	}

	var prefix string
	if d.BaseURL == "" {
		prefix = "/"
	} else {
		prefix = d.BaseURL + "/"
	}

	args.Asset.URI = prefix + args.Partition + "/" + args.Asset.ID + "/" + args.Asset.Name

	return args.Asset, err
}

func (d disk) Delete(ctx context.Context, args asset.DeleteArgs) error {
	var (
		path      = filepath.Join(d.Path, args.Partition, args.ID)
		info, err = os.Stat(path)
	)

	if os.IsNotExist(err) {
		err = nil
	}

	if err != nil {
		return err
	}

	if info.IsDir() {
		err = os.RemoveAll(path)
	}

	return err
}
