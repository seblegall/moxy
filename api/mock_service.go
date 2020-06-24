package api

import "encoding/json"

type Store interface {
	List() ([]Mock, error)
	Get(string, string) (Mock, error)
	Add(Mock) error
}

type MockService struct {
	store Store
}

func NewMockService(store Store) *MockService {
	return &MockService{store:store}
}

func (m *MockService) List() ([]Mock, error) {
	return m.store.List()
}

func (m *MockService) Get(method, path string) (Mock, error) {
	return m.store.Get(method, path)
}


func (m *MockService) Add(method, path string, status int, body json.RawMessage) (Mock, error) {
	mock := NewMock(method, path, status, body)
	err := m.store.Add(mock)
	if err != nil {
		return Mock{}, err
	}

	return mock, nil

}

