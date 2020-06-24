package http

import "fmt"

type server struct {
	handler *handler
}

func NewServer(h *handler) *server {
	return &server{
		handler: h,
	}
}

// Serve launch the webserver
func (s *server) Serve(port int) {
	if err := s.handler.engine.Run(fmt.Sprintf(":%d", port)); err != nil {
		panic(err)
	}
}
