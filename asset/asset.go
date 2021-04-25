package asset

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/fluxynet/gocipe/util"
	"github.com/google/uuid"
)

var (
	// ErrorAssetNotFound returned if trying to access an asset by name and id; but it cannot be found in storage
	ErrorAssetNotFound = &StorageError{Code: http.StatusNotFound, Message: "asset not found"}
)

type Asset struct {
	ID         string
	Name       string
	Size       int
	Mime       string
	URI        string
	Resource   string
	ResourceID string
}

func (a Asset) String() string {
	return "" +
		"ID   = " + a.ID + "\n" +
		"Name = " + a.Name + "\n" +
		"Size = " + strconv.Itoa(a.Size) + "\n" +
		"Mime = " + a.Mime + "\n" +
		"URI  = " + a.URI
}

type StorageError struct {
	Code    int
	Message string
}

func (e StorageError) Error() string {
	return strconv.Itoa(e.Code) + `: ` + e.Message
}

// Manager for assets with http capabilities
type Manager struct {
	validators        map[string]Validator
	Prefix            string
	Storage           Storage
	PartitionResolver PartitionResolver
	UploadHandler     UploadHandler
	DeleteHandler     DeleteHandler
}

// Serve the upload handler
func (m Manager) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	case http.MethodDelete:
		m.Delete(w, r)
	case http.MethodPost:
		m.Create(w, r)
	}
}

func (m Manager) Create(w http.ResponseWriter, r *http.Request) {
	var (
		kind          = strings.TrimPrefix(r.URL.Path, m.Prefix)
		ctx           = r.Context()
		validator, ok = m.validators[kind]
		partition     string
		storageError  *StorageError
		err           error
		asset         Asset
		response      []byte
		resCode       int
	)

	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	defer func() {
		if storageError == nil {
			response = []byte(`{"uri": "` + asset.URI + `"}`)
			resCode = http.StatusOK
		} else {
			var msg = strings.Replace(storageError.Message, `"`, "", -1)
			response = []byte(`{"error": "` + msg + `"}`)
			resCode = storageError.Code
			fmt.Println(">>> ", resCode)
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
		w.WriteHeader(resCode)
	}()

	if m.PartitionResolver != nil {
		partition = m.PartitionResolver.GetPartition(ctx, r)
	}

	var body = util.ReadAll(r.Body)

	if body == nil {
		storageError = &StorageError{
			Code: http.StatusRequestEntityTooLarge,
		}
		return
	}

	var query = r.URL.Query()
	asset.Name = query.Get("name")
	asset.Resource = query.Get("resource")
	asset.ResourceID = query.Get("resource_id")

	if asset.Name == "" {
		storageError = &StorageError{
			Code:    http.StatusBadRequest,
			Message: "missing query parameter: name",
		}
	} else if asset.Resource == "" {
		storageError = &StorageError{
			Code:    http.StatusBadRequest,
			Message: "missing query parameter: resource",
		}
	} else if asset.ResourceID == "" {
		storageError = &StorageError{
			Code:    http.StatusBadRequest,
			Message: "missing query parameter: resource_id",
		}
	}

	if storageError != nil {
		return
	}

	asset.ID = uuid.NewString()
	asset.Size = len(body)
	asset.Mime = http.DetectContentType(body)

	if i := strings.IndexRune(asset.Mime, ';'); i != -1 {
		asset.Mime = asset.Mime[:i]
	}

	if e := validator.Validate(
		ctx,
		&ValidateArgs{
			Request:   r,
			Partition: partition,
			Kind:      kind,
			Asset:     asset,
			Data:      body,
		}); e != nil {
		storageError = &StorageError{
			Code:    http.StatusBadRequest,
			Message: e.Error(),
		}

		return
	}

	if asset, err = m.Storage.Store(ctx, StoreArgs{
		Partition: partition,
		Kind:      kind,
		Asset:     asset,
		Data:      body,
	}); err != nil {
		storageError = &StorageError{
			Code:    http.StatusInternalServerError,
			Message: "failed to store file: " + err.Error(),
		}

		return
	}

	if m.UploadHandler != nil {
		asset, storageError = m.UploadHandler.OnUpload(
			ctx,
			&UploadHandlerArgs{
				Request:   r,
				Storage:   m.Storage,
				Partition: partition,
				Kind:      kind,
				Asset:     asset,
			},
		)

		return
	}
}

func (m Manager) Delete(w http.ResponseWriter, r *http.Request) {
	var (
		path            = strings.TrimPrefix(r.URL.Path, m.Prefix)
		asset           Asset
		kind, partition string
		err             *StorageError
		ctx             = r.Context()
	)

	if i := strings.IndexRune(path, '/'); i != -1 {
		kind, asset.ID = path[:i], path[i+1:]
	}

	if _, ok := m.validators[kind]; !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	if m.PartitionResolver != nil {
		partition = m.PartitionResolver.GetPartition(ctx, r)
	}

	if e := m.Storage.Delete(ctx, DeleteArgs{
		Partition: partition,
		ID:        asset.ID,
	}); e != nil {
		err = &StorageError{
			Code:    http.StatusInternalServerError,
			Message: "failed to delete file: " + e.Error(),
		}
	}

	if err != nil {
		var msg = strings.Replace(err.Message, `"`, "", -1)
		w.Write([]byte(`{"error": "` + msg + `"}`))
		w.WriteHeader(err.Code)
		return
	}

	if r.Method == http.MethodDelete && m.DeleteHandler != nil {
		err = m.DeleteHandler.OnDelete(
			ctx,
			&DeleteHandlerArgs{
				Request:   r,
				Storage:   m.Storage,
				Partition: partition,
				Kind:      kind,
				ID:        asset.ID,
			},
		)
	}

	w.WriteHeader(http.StatusOK)
}

// Register a validator
func (m *Manager) Register(kind string, validator Validator) *Manager {
	if m.validators == nil {
		m.validators = make(map[string]Validator)
	}

	m.validators[kind] = validator
	return m
}

func (m *Manager) Uploads() []Upload {
	var (
		uploads []Upload
		prefix  = strings.TrimSuffix(m.Prefix, "/")
	)

	if i := strings.LastIndexByte(prefix, '/'); i == -1 {
		prefix = "/"
	} else {
		prefix = prefix[i:]
	}

	for kind := range m.validators {
		if t, ok := m.validators[kind].(Types); ok {
			var u = upload{
				types: t.Types(),
				path:  prefix + "/" + kind,
			}

			uploads = append(uploads, u)
		}
	}

	return uploads
}

type ValidateArgs struct {
	Request   *http.Request
	Partition string
	Kind      string
	Asset     Asset
	Data      []byte
}

// Validator for upload
type Validator interface {
	Validate(ctx context.Context, args *ValidateArgs) error
}

// Types reports which media types are supported by a validator if the validator is validates media types
type Types interface {
	Types() []string
}

// Upload defines media uploading paths compatible with openapi interface
type Upload interface {
	Path() string
	Types() []string
}

type upload struct {
	types []string
	path  string
}

func (u upload) Types() []string {
	return u.types
}

func (u upload) Path() string {
	return u.path
}

// StoreArgs represents arguments for storing in storage
type StoreArgs struct {
	Partition string
	Kind      string
	Asset     Asset
	Data      []byte
}

// DeleteArgs represents arguments for deleting from storage
type DeleteArgs struct {
	Partition string
	ID        string
}

// Storage engine
type Storage interface {
	Store(ctx context.Context, args StoreArgs) (Asset, error)
	Delete(ctx context.Context, args DeleteArgs) error
}

// UploadHandlerArgs represents arguments to the upload handler
type UploadHandlerArgs struct {
	Request   *http.Request
	Storage   Storage
	Partition string
	Kind      string
	Asset     Asset
}

// UploadHandler allows a callback to a new upload
type UploadHandler interface {
	OnUpload(ctx context.Context, args *UploadHandlerArgs) (Asset, *StorageError)
}

type compositeUploadHandler struct {
	Handlers []UploadHandler
}

func (c compositeUploadHandler) OnUpload(ctx context.Context, args *UploadHandlerArgs) (Asset, *StorageError) {
	var err *StorageError
	for i := range c.Handlers {
		args.Asset, err = c.Handlers[i].OnUpload(ctx, args)

		if err != nil {
			return args.Asset, err
		}
	}

	return args.Asset, nil
}

// CompositeUploadHandler allows multiple upload handlers to be executed
func CompositeUploadHandler(handlers ...UploadHandler) UploadHandler {
	return compositeUploadHandler{Handlers: handlers}
}

// DeleteHandlerArgs represents arguments to the delete handler
type DeleteHandlerArgs struct {
	Request   *http.Request
	Storage   Storage
	Partition string
	Kind      string
	ID        string
}

type compositeDeleteHandler struct {
	Handlers []DeleteHandler
}

func (c compositeDeleteHandler) OnDelete(ctx context.Context, args *DeleteHandlerArgs) *StorageError {
	var err *StorageError
	for i := range c.Handlers {
		err = c.Handlers[i].OnDelete(ctx, args)

		if err != nil {
			return err
		}
	}

	return nil
}

// CompositeDeleteHandler allows multiple delete handlers to be executed
func CompositeDeleteHandler(handlers ...DeleteHandler) DeleteHandler {
	return compositeDeleteHandler{Handlers: handlers}
}

// DeleteHandler allows a callback when an asset has been deleted
type DeleteHandler interface {
	OnDelete(ctx context.Context, args *DeleteHandlerArgs) *StorageError
}

// PartitionResolver determines which partition to store an asset
// must return the partition given an *http.Request
type PartitionResolver interface {
	GetPartition(ctx context.Context, r *http.Request) string
}

type staticPartitionResolver struct {
	partition string
}

func (s staticPartitionResolver) GetPartition(ctx context.Context, r *http.Request) string {
	return s.partition
}

// StaticPartitionResolver always resolves a specified string
func StaticPartitionResolver(partition string) PartitionResolver {
	return staticPartitionResolver{partition: partition}
}
