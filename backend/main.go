package main

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func getData(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, 4+5)
}

func main() {
	router := gin.Default()
	router.GET("/data", getData)

	router.Run("localhost:8080")
}
