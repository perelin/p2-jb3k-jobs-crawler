package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	db "p2lab/recruitbot3000/pkg/db"
	"strconv"

	"github.com/gin-gonic/gin"
)

func printDir() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
	log.Println(dir)
}
func setupRouter() *gin.Engine {

	r := gin.Default()
	r.Use(gin.Logger())
	r.LoadHTMLGlob(os.Getenv("HEROKU_STATIC_PATH") + "templates/*.tmpl.html")

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
		c.HTML(http.StatusOK, "status.tmpl.html", gin.H{
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
	printDir()
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":" + port)
}
