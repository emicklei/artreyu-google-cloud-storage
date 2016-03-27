package main

import (
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/artreyu/model"
	"github.com/emicklei/artreyu/transport"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/storage/v1"
)

const (
	// This scope allows the application full control over resources in Google Cloud Storage
	scope = storage.DevstorageFullControlScope
)

type Repository struct {
	settings *model.Settings
	config   model.RepositoryConfig
}

func NewRepository(config model.RepositoryConfig, settings *model.Settings) Repository {
	return Repository{settings, config}
}

func (r Repository) ID() string { return "gcs" }

func (r Repository) Store(a model.Artifact, source string) error {
	service, _, err := r.gcsService()
	if err != nil {
		return err
	}
	repo := "releases"
	if a.IsSnapshot() {
		repo = "snapshots"
	}
	destination := filepath.Join(strings.TrimLeft(r.config.Path, "/"), repo, a.StorageLocation(r.settings.OS, a.AnyOS))
	object := &storage.Object{Name: destination}
	file, err := os.Open(source)
	if err != nil {
		model.Printf("error opening %q: %v", source, err)
		return err
	}
	model.Printf("uploading gcr.io/%s%s to %s\n", r.config.Bucket, source, destination)
	if _, err := service.Objects.Insert(r.config.Bucket, object).Media(file).Do(); err != nil {
		model.Printf("unable to create file %v at location %v\n\n", source, destination)
	}
	return nil
}

func (r Repository) Fetch(a model.Artifact, destination string) error {
	service, client, err := r.gcsService()
	if err != nil {
		return err
	}
	repo := "releases"
	if a.IsSnapshot() {
		repo = "snapshots"
	}
	source := r.config.URL + filepath.Join(strings.TrimLeft(r.config.Path, "/"), repo, a.StorageLocation(r.settings.OS, a.AnyOS))
	model.Printf("downloading gcr.io/%s%s to %s\n", r.config.Bucket, source, destination)

	res, err := service.Objects.Get(r.config.Bucket, source).Do()
	if err != nil {
		model.Printf("failed to get %s/%s: %s.", r.config.Bucket, source, err)
		return err
	}

	return transport.HttpGetFile(client, res.MediaLink, destination)
}

func (r Repository) Exists(a model.Artifact) bool {
	service, _, err := r.gcsService()
	if err != nil {
		return false
	}
	repo := "releases"
	if a.IsSnapshot() {
		repo = "snapshots"
	}
	source := r.config.URL + filepath.Join(strings.TrimLeft(r.config.Path, "/"), repo, a.StorageLocation(r.settings.OS, a.AnyOS))

	if _, err := service.Objects.Get(r.config.Bucket, source).Do(); err != nil {
		return false
	}
	return true
}

func (r Repository) gcsService() (*storage.Service, *http.Client, error) {
	// Authentication is provided by the gcloud tool when running locally, and
	// by the associated service account when running on Compute Engine.
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		model.Printf("unable to get default client: %v", err)
		return nil, nil, err
	}
	service, err := storage.New(client)
	if err != nil {
		model.Printf("unable to create storage service: %v", err)
		return nil, nil, err
	}
	if _, err := service.Buckets.Get(r.config.Bucket).Do(); err != nil {
		model.Printf("bucket %s does not exist", r.config.Bucket)
		return nil, nil, err
	}
	return service, client, nil
}
