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

	r := gin.Default()
	r.LoadHTMLGlob("templates/*")

	r.GET("/status", func(c *gin.Context) {
		count := db.GetJobAdCount("monster")
		lastEntryTime := db.GetLastEntryDate("monster")
		lastEntryTimeString := lastEntryTime.String()
		reply := "Total entries: " + strconv.Itoa(count) + "\nLast entry time: " + lastEntryTimeString
		c.String(http.StatusOK, reply)
	})

	r.GET("/", func(c *gin.Context) {

		monsterEntries := db.GetJobAdCount("monster")
		monsterLastEntryTime := db.GetLastEntryDate("monster")
		monsterLastEntryString := monsterLastEntryTime.String()
		ssEntries := db.GetJobAdCount("stepstone")
		ssLastEntryTime := db.GetLastEntryDate("stepstone")
		ssLastEntryString := ssLastEntryTime.String()
		c.HTML(http.StatusOK, "status.html", gin.H{
			"monsterEntries":   monsterEntries,
			"monsterLastEntry": monsterLastEntryString,
			"ssEntries":        ssEntries,
			"ssLastEntry":      ssLastEntryString,
		})
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
