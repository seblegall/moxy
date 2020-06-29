package http

import "fmt"


type handler interface {
	Run(addr ...string) (err error)
}

type server struct {
	handler handler
}

func NewServer(h handler) *server {
	return &server{
		handler: h,
	}
}

// Serve launch the webserver
func (s *server) Serve(port int) error {
	if err := s.handler.Run(fmt.Sprintf(":%d", port)); err != nil {
		return err
	}

	return nil
}
