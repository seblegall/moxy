package http

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/seblegall/moxy/api"
	"github.com/seblegall/moxy/config"
)



type handler struct {
	engine *gin.Engine
	backend config.Backend
	mockService mockService
}

type mockService interface {
	List() ([]api.Mock, error)
	Get(string, string) (api.Mock, error)
	Add(string, string, int, json.RawMessage) (api.Mock, error)
}

func NewHandler(mock mockService, backend config.Backend) *handler {
	h := &handler{
		engine:      gin.Default(),
		backend: backend,
		mockService: mock,
	}

	//CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization", "Remote-User")
	h.engine.Use(cors.New(corsConfig))
	h.engine.Use()

	h.engine.Any("/*proxyPath", h.handleMocker, h.handleMock, h.handleProxy)

	return h
}

func (h *handler) handleMock(c *gin.Context) {
	mock, err := h.mockService.Get(c.Request.Method, c.Request.URL.Path)
	if err == nil {
		c.Data(mock.StatusCode, "text/json", mock.Body)
		c.Abort()
	}
}

//proxy proxify a request to the backend
func (h *handler) handleProxy(c *gin.Context) {
	proxy := httputil.NewSingleHostReverseProxy(h.backend.URL())
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = h.backend.URL().Host
		req.URL.Scheme = h.backend.URL().Scheme
		req.URL.Host = h.backend.URL().Host
		req.URL.Path = fmt.Sprintf("%s%s", h.backend.URL().Path, c.Param("proxyPath"))
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}
