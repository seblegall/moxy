package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	pa "path"
	"strings"
)

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
	p := strings.TrimLeft(pa.Clean(path), "/")
	return m.store.Get(strings.ToUpper(method), p)
}


func (m *MockService) Add(method, path string, status int, body json.RawMessage) (Mock, error) {

	meth := strings.ToUpper(method)

	if !isMethodValid(meth) {
		return Mock{}, fmt.Errorf(InvalidMethod)
	}

	mock := NewMock(meth, strings.TrimLeft(pa.Clean(path), "/"), status, body)
	err := m.store.Add(mock)
	if err != nil {
		return Mock{}, err
	}

	return mock, nil
}

func isMethodValid(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	
	}
}

