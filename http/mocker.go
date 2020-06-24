package http

import (
	"encoding/json"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (h *handler) handleMocker(c *gin.Context) {

	if c.Request.URL.Path == "/moxy/mocks" && c.Request.Method == http.MethodPost {
		h.addMock(c)
	}

	if c.Request.URL.Path == "/moxy/mocks" && c.Request.Method == http.MethodGet {
		h.listMock(c)
	}

}

func (h *handler) addMock(c *gin.Context) {
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

func (h *handler) listMock(c *gin.Context) {
	mocks, err := h.mockService.List()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		c.Abort()
		return
	}

	c.JSON(http.StatusOK, mocks)
	c.Abort()
}