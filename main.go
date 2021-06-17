package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
)

// type request struct {
// 	name, surname string
// }

func main() {
	router := gin.Default()

	router.POST("/new", func (c *gin.Context) {
		fmt.Println(c.Request.Body)
	})

	router.Run(":1337")
}