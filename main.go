package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

var (
	backendUrl string
	backend *url.URL
)


func proxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(backend)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = backend.Host
		req.URL.Scheme = backend.Scheme
		req.URL.Host = backend.Host
		req.URL.Path = fmt.Sprintf("%s%s", backend.Path, c.Param("proxyPath"))
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

func main() {

	flag.StringVar(&backendUrl, "backend", "", "the backend url")
	port := flag.Int("port", 8080, "the exposed port")
	flag.Parse()

	if backendUrl == "" {
		log.Fatal("a backend url must be provided. None given.")
	}

	if u, err := url.Parse(backendUrl); err != nil {
		log.Fatal("Unable to parse given backend url")
	} else {
		backend = u
	}

	r := gin.Default()

	//CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization", "Remote-User")
	r.Use(cors.New(corsConfig))

	r.Any("/*proxyPath", proxy)

	r.Run(fmt.Sprintf(":%d", *port))
}
