package main

import (
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/emicklei/artreyu/model"
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
	// Authentication is provided by the gcloud tool when running locally, and
	// by the associated service account when running on Compute Engine.
	client, err := google.DefaultClient(context.Background(), scope)
	if err != nil {
		log.Printf("Unable to get default client: %v", err)
		return err
	}
	service, err := storage.New(client)
	if err != nil {
		log.Printf("Unable to create storage service: %v", err)
		return err
	}
	if _, err := service.Buckets.Get(r.config.Bucket).Do(); err != nil {
		log.Printf("Bucket %s does not exist", r.config.Bucket)
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
		log.Printf("Error opening %q: %v", source, err)
		return err
	}
	model.Printf("uploading %s to %s\n", source, destination)
	if _, err := service.Objects.Insert(r.config.Bucket, object).Media(file).Do(); err != nil {
		log.Printf("Unable to create file %v at location %v\n\n", source, destination)
	}
	return nil
}

func (r Repository) Fetch(a model.Artifact, destination string) error { return nil }

func (r Repository) Exists(a model.Artifact) bool { return false }
