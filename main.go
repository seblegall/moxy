package main

import (
	"flag"
	"log"

	"github.com/seblegall/moxy/api"
	"github.com/seblegall/moxy/config"
	"github.com/seblegall/moxy/http"
	"github.com/seblegall/moxy/store"
)


func main() {

	//Read configuration from flag
	url := flag.String("backend", "", "the backend url")
	port := flag.Int("port", 8080, "the exposed port")
	flag.Parse()

	backend, err := config.NewBackend(*url)
	if err != nil {
		log.Fatal(err.Error())
	}

	s := http.NewServer(http.NewHandler(api.NewMockService(store.NewMap()), backend))
	s.Serve(*port)
}