package memory

import (
	"context"
	"sync"

	"movix/metadata/internal/repository"
	"movix/metadata/pkg/model"
)

type Repository struct {
	mut      *sync.RWMutex
	metadata map[string]*model.Metadata
}

func New() *Repository {
	return &Repository{
		mut:      &sync.RWMutex{},
		metadata: make(map[string]*model.Metadata),
	}
}

func (r *Repository) Get(_ context.Context, id string) (*model.Metadata, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()
	m, ok := r.metadata[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	return m, nil
}

func (r *Repository) Put(_ context.Context, id string, metadata *model.Metadata) error {
	r.mut.Lock()
	r.metadata[id] = metadata
	r.mut.Unlock()
	return nil
}
