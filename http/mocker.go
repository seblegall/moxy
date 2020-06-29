package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/seblegall/moxy/api"
)

type mockerHandler struct {
	engine *gin.Engine
	mockService *api.MockService
}

func NewMockerHandler(mock *api.MockService) *mockerHandler {
	h := &mockerHandler{
		engine:      gin.Default(),
		mockService: mock,
	}

	//CORS
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowAllOrigins = true
	corsConfig.AddAllowHeaders("Authorization", "Remote-User")
	h.engine.Use(cors.New(corsConfig))
	h.engine.Use()

	h.engine.GET("/moxy/api/mocks", h.listMock)
	h.engine.POST("/moxy/api/mocks", h.addMock)

	return h
}

func (h *mockerHandler) addMock(c *gin.Context) {
	type mockRequest struct {
		Path string `json:"path"  binding:"required"`
		Method string `json:"method"  binding:"required"`
		StatusCode int `json:"status_code"  binding:"required"`
		Body json.RawMessage `json:"body"  binding:"required"`
	}

	var query mockRequest
	if err := c.ShouldBindJSON(&query); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	mock, err := h.mockService.Add(query.Method, query.Path, query.StatusCode,query.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusCreated, mock)
	c.Abort()
}

func (h *mockerHandler) listMock(c *gin.Context) {
	mocks, err := h.mockService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, mocks)
	c.Abort()
}

func (h *mockerHandler) Run(addr ...string) (err error) {
	return h.engine.Run(addr...)
}