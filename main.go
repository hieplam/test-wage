package main

import (
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	fmt.Println("init server")
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	err := r.Run()
	if err != nil {
		log.Fatal("err when init server", err)
	}
}
