package helper

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Error("Error loading local .env file")
	}

	logLevel, _ := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logLevel)

}

func DelayForMonsterAPI() {
	//log.Debug("Waiting delay for api throtteling...")
	delay, _ := strconv.Atoi(os.Getenv("DELAY"))
	time.Sleep(time.Duration(delay) * time.Millisecond)
}
