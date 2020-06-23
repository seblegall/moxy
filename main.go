package main

import (
	"encoding/json"
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
	mockList []mock
)

type mock struct {
	Path string `json:"path"`
	StatusCode int `json:"status_code"`
	Body json.RawMessage `json:"body"`
}

func main() {

	//Read configuration from flag
	flag.StringVar(&backendUrl, "backend", "", "the backend url")
	port := flag.Int("port", 8080, "the exposed port")
	flag.Parse()

	//Parse backend url
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

	r.Any("/*proxyPath", handleMoxy)

	r.Run(fmt.Sprintf(":%d", *port))
}


func handleMoxy(c *gin.Context) {

	//catch mock management routes
	if handleMock(c) {
		return
	}

	//if path match mock, return the mock
	for _, m := range mockList {
		if m.Path == c.Request.URL.Path {
			c.Data(m.StatusCode, "text/json", m.Body)
			return
		}
	}

	//Else proxify the request to the backend
	proxy(c)
}

//proxy proxify a request to the backend
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

func handleMock(c *gin.Context) bool {

	if c.Request.URL.Path == "/mocks" && c.Request.Method == http.MethodPost {
		addMock(c)
		return true
	}

	if c.Request.URL.Path == "/mocks" && c.Request.Method == http.MethodGet {
		listMock(c)
		return true
	}

	return false
}

func listMock(c *gin.Context) {
	c.JSON(http.StatusOK, mockList)
}


func addMock(c *gin.Context) {

	type mockRequest struct {
		StatusCode int `json:"status_code"`
		Body json.RawMessage `json:"body"`
	}

	var query mockRequest
	if err := c.BindJSON(&query); err != nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	mockList = append(mockList, mock{
		Path: c.Request.URL.Path,
		StatusCode: query.StatusCode,
		Body: query.Body,
	})

	c.Status(http.StatusCreated)
}