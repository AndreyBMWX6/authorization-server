package storage

import (
	"context"
	"sync"
)

// Storage - for simplicity we will use inmem storage
type Storage struct {
	kv map[string]interface{}
	mu *sync.RWMutex
}

func New() *Storage {
	return &Storage{
		kv: make(map[string]interface{}),
		mu: &sync.RWMutex{},
	}
}

func (s *Storage) Insert(ctx context.Context, key string, value interface{}) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, ok := s.kv[key]
	if ok {
		return ErrAlreadyExists
	}

	s.kv[key] = value
	return nil
}
