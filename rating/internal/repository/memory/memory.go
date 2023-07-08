package memory

import (
	"context"
	"movix/rating/internal/repository"
	"movix/rating/pkg/model"
	"sync"
)

type Memory struct {
	mut  *sync.RWMutex
	data map[model.RecordType]map[model.RecordID][]model.Rating
}

func New() *Memory {
	return &Memory{
		mut:  &sync.RWMutex{},
		data: make(map[model.RecordType]map[model.RecordID][]model.Rating),
	}
}

func (m *Memory) Get(_ context.Context, recordType model.RecordType, id model.RecordID) ([]model.Rating, error) {
	m.mut.RLock()
	defer m.mut.RUnlock()

	records, ok := m.data[recordType]
	if !ok {
		return nil, repository.ErrRecordTypeNotFound
	}

	if record, ok := records[id]; ok {
		return record, nil
	}

	return nil, repository.ErrRecordNotFound
}

func (m *Memory) Put(_ context.Context, recordType model.RecordType, id model.RecordID, rating model.Rating) error {
	m.mut.Lock()
	defer m.mut.Unlock()

	if _, ok := m.data[recordType]; !ok {
		m.data[recordType] = make(map[model.RecordID][]model.Rating)
	}

	m.data[recordType][id] = append(m.data[recordType][id], rating)
	return nil
}
