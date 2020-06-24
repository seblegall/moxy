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
)


func main() {

	//Read configuration from flag
	url := flag.String("backend", "", "the backend url")
	port := flag.Int("port", 8080, "the exposed port")
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

	s := http.NewServer(http.NewHandler(mock, backend))
	s.Serve(*port)
}