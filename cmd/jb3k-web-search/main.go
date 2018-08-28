package main

import (
	"log"
	"net/http"
	"os"
	db "p2lab/recruitbot3000/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	// Disable Console Color
	// gin.DisableConsoleColor()
	r := gin.Default()

	r.GET("/status", func(c *gin.Context) {
		count := db.GetJobAdCount()
		lastEntryTime := db.GetLastEntryDate()
		lastEntryTimeString := lastEntryTime.String()
		reply := "Total entries: " + strconv.Itoa(count) + "\nLast entry time: " + lastEntryTimeString
		c.String(http.StatusOK, reply)
	})

	// Ping test
	r.GET("/ping2", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	return r
}

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":" + port)
}
