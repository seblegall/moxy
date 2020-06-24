package store

import (
	"fmt"
	"sync"

	"github.com/seblegall/moxy/api"
)

type MapStore struct {
	mu sync.Mutex
	mocks []api.Mock
}

func NewMap() *MapStore {
	var mocks []api.Mock
	return &MapStore{
		mocks: mocks,
	}
}

func (ms *MapStore) List() ([]api.Mock, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.mocks, nil
}

func (ms *MapStore) Get(method, path string) (api.Mock, error) {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	return ms.get(method, path)
}

func (ms *MapStore) get(method, path string) (api.Mock, error) {
	for _, m := range ms.mocks {
		if m.Method == method && m.Path == path {
			return m, nil
		}
	}
	return api.Mock{}, fmt.Errorf(MockNotFound)
}

func (ms *MapStore) Add(mock api.Mock) error {
	ms.mu.Lock()
	defer ms.mu.Unlock()

	_, err := ms.get(mock.Method, mock.Path)
	if err == nil {
		return fmt.Errorf(MockAlreadyExists)
	}

	ms.mocks = append(ms.mocks, mock)
	return nil
}