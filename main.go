package main

import (
	"test-case/models"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	models.ConnectedDB()

	routes(r)
	r.Run(":8080")

}
