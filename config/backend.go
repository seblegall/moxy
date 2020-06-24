package config

import (
	"fmt"
	"net/url"
)

type Backend struct {
	url *url.URL
}

func NewBackend(URL string) (Backend, error) {
	if URL == "" {
		return Backend{}, fmt.Errorf("a backend url must be provided. None given")
	}

	u, err := url.Parse(URL)
	if err != nil {
		return Backend{}, fmt.Errorf("unable to parse given backend url")
	}

	return Backend{url:u}, nil
}

func (b Backend) URL() *url.URL {
	return b.url
}