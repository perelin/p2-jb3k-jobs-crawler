package main

import (
	"fmt"
	"os"
	db "p2lab/recruitbot3000/pkg/db"
	"strconv"
	"time"

	"github.com/joho/godotenv"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
)

var request = gorequest.New()

func init() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Error("Error loading local .env file")
	}

	logLevel, _ := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logLevel)

}

func delayForMonsterAPI() {
	//log.Debug("Waiting delay for api throtteling...")
	delay, _ := strconv.Atoi(os.Getenv("DELAY"))
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func main() {
	lastEntryTime := db.GetLastEntryDate()
	fmt.Println(lastEntryTime)
	jobAds := db.GetAllJobs(true)

	// walk over ever job
	// - check ifjob page returns 404
	// -- set active to false

	for _, jobAd := range jobAds {
		delayForMonsterAPI()
		resp, _, errs := request.Get(jobAd.URL).End()

		if errs != nil {
			log.Error("Job Ad page couldn´t be loaded: ", errs)
		}
		if resp.StatusCode == 404 {
			log.WithFields(log.Fields{"url": jobAd.URL}).Debug("job ad page returns 404, job ad might no longer be active")
			db.UpdateJobActiveStatus(int(jobAd.ID), false)
		} else if resp.StatusCode == 200 {
			log.WithFields(log.Fields{"url": jobAd.URL}).Debug("job ad seems to be alive")
		}
	}
}
