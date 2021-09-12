package routes

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Return the home page
func Home(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}
