package main

import (
	"github.com/gin-gonic/gin"
	"sihamais/searchingo/internal/routes"
)

func main() {

	r := gin.Default()

	r.LoadHTMLGlob("../templates/*")

	r.GET("/", routes.Home)
	r.GET("/search", routes.Search)

	r.Run(":8080")
}
