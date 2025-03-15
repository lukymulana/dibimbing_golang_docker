package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
}

func NewHandler() Handler {
	return Handler{}
}

func (h *Handler) HelloWorld(c *gin.Context) {
	c.JSON(http.StatusOK, []byte(`{"hello": "world}`))
}
