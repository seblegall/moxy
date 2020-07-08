package http

import (
	"bytes"
	"encoding/json"
	"html/template"
	"net/http"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/packr/v2"
	"github.com/seblegall/moxy/api"
	"github.com/sirupsen/logrus"
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

	t, err := loadTemplate()
	if err != nil {
		logrus.Fatal(err.Error())
	}
	h.engine.SetHTMLTemplate(t)
	h.engine.Use(cors.New(corsConfig))

	h.engine.GET("/moxy/api/mocks", h.listMock)
	h.engine.POST("/moxy/api/mocks", h.addMock)
	h.engine.Any("/", func(c *gin.Context) {
		c.Redirect(301, "/moxy/dashboard")
	})
	h.engine.GET("/moxy/dashboard", h.renderDashboard)
	h.engine.GET("/moxy/dashboard/add", h.renderAddMock)

	return h
}

func loadTemplate() (*template.Template, error) {
	t := template.New("")
	box := packr.New("myBox", "../templates")
	for _, file := range box.List() {
		h, err := box.FindString(file)
		if err != nil {
			return nil, err
		}
		t, err = t.New(file).Parse(h)
		if err != nil {
			return nil, err
		}
	}
	return t, nil
}

func (h *mockerHandler) renderDashboard(c *gin.Context) {
		mocks, err := h.mockService.List()
		if err != nil {
			c.AbortWithError(http.StatusInternalServerError, err)
			return
		}

		c.HTML(200, "dashboard.gohtml", gin.H{
			"mocks": mocks,
		})
}

func (h *mockerHandler) renderAddMock(c *gin.Context) {
	c.HTML(200, "add_mock.gohtml", gin.H{})
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

	var out bytes.Buffer
	var body json.RawMessage
	err := json.Indent(&out, query.Body, "", "     ")
	if err != nil {
		body = query.Body
	} else {
		body = out.Bytes()
	}

	mock, err := h.mockService.Add(query.Method, query.Path, query.StatusCode,body)
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