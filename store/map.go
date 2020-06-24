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

	ms.mocks = append(ms.mocks, mock)

	return nil
}