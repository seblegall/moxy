package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/seblegall/moxy/api"
	"github.com/seblegall/moxy/config"
	"github.com/seblegall/moxy/http"
	"github.com/seblegall/moxy/store"
	"golang.org/x/sync/errgroup"
)


func main() {

	//Read configuration from flag
	url := flag.String("backend", "", "the backend url")
	proxyPort := flag.Int("proxy-port", 8080, "the exposed port for the proxy")
	adminPort := flag.Int("admin-port", 8181, "the exposed port for the admin dashboard")
	file := flag.String("mock-file", "", "a json file containing mocks to load when starting the server")
	flag.Parse()

	backend, err := config.NewBackend(*url)
	if err != nil {
		log.Fatal(err.Error())
	}

	mock := api.NewMockService(store.NewMap())

	//Load mocks from file
	if *file != "" {
		f, err := os.Open(*file)
		if err != nil {
			log.Fatal(err.Error())
		}

		b, err := ioutil.ReadAll(f)
		if err != nil {
			log.Fatal(err.Error())
		}

		if err := mock.Load(b); err != nil {
			log.Fatal(err.Error())
		}
	}

	var g errgroup.Group

	g.Go(func() error {
		s := http.NewServer(http.NewProxyHandler(mock, backend))
		return s.Serve(*proxyPort)
	})

	g.Go(func() error {
		s := http.NewServer(http.NewMockerHandler(mock))
		return s.Serve(*adminPort)
	})

	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

}