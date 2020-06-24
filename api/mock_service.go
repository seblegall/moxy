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

func (ms *MockService) List() ([]Mock, error) {
	return ms.store.List()
}

func (ms *MockService) Get(method, path string) (Mock, error) {
	p := strings.TrimLeft(pa.Clean(path), "/")
	return ms.store.Get(strings.ToUpper(method), p)
}


func (ms *MockService) Add(method, path string, status int, body json.RawMessage) (Mock, error) {

	meth := strings.ToUpper(method)
	if !isMethodValid(meth) {
		return Mock{}, fmt.Errorf(InvalidMethod)
	}

	if path == "" {
		return Mock{}, fmt.Errorf("path cannot be empty")
	}

	mock := NewMock(meth, strings.TrimLeft(pa.Clean(path), "/"), status, body)
	err := ms.store.Add(mock)
	if err != nil {
		return Mock{}, err
	}

	return mock, nil
}

func (ms *MockService) Load(rawMocks []byte) error {

	var mocks []Mock
	if err := json.Unmarshal(rawMocks, &mocks); err != nil {
		return err
	}

	for i, mock := range mocks {
		_, err := ms.Add(mock.Method, mock.Path, mock.StatusCode, mock.Body)
		if err != nil {
			return fmt.Errorf("error when importing mock #%d : %s", i, err.Error())
		}
	}

	return nil
}

func isMethodValid(method string) bool {
	switch method {
	case http.MethodGet, http.MethodPost, http.MethodPatch, http.MethodPut, http.MethodDelete:
		return true
	default:
		return false
	
	}
}

