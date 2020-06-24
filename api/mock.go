package api

import (
	"encoding/json"

	uuid "github.com/satori/go.uuid"
)

type Mock struct {
	ID string `json:"id"`
	Method string `json:"method"`
	Path string `json:"path"`
	StatusCode int `json:"status_code"`
	Body json.RawMessage `json:"body"`
}


func NewMock(method, path string, status int, body json.RawMessage) Mock {
	id := uuid.NewV4()
	
	return Mock{
		ID:         id.String(),
		Method: 	method,
		Path:       path,
		StatusCode: status,
		Body:       body,
	}
}