package main

import (
	"net/url"
	"os"

	"p2lab/recruitbot3000/pkg/db"

	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
)

var request = gorequest.New()

var newJobs int
var newJobsForQuery int
var jobsFoundForQuery int
var expectedResults int

func init() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Error("Error loading local .env file")
	}

	logLevel, _ := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logLevel)

}

func getJobAdsForJobNames() {
	jobNames := db.GetJobNames()

	for _, jobName := range jobNames {
		query := url.QueryEscape(jobName.Text)
		//expectedResults = GetTotalResultsForQuery(query)
		expectedResults = GetTotalResultsForQueryFromHTML(query)
		log.WithFields(log.Fields{"Job Query": query, "expectedResults": expectedResults}).Info("Starting new job name query")
		jobsAvailable := checkIfJobsAreAvailableForQuery(query)
		if !jobsAvailable {
			continue
		}
		continueToNextPage := true
		jobsFoundForQuery = 0
		newJobsForQuery = 0
		for i := 1; continueToNextPage; i++ {
			continueToNextPage = scrapeJobListingsFromJSON(query, i)
			log.Debug("Continues to next page? ", continueToNextPage)
		}

		log.WithFields(log.Fields{"query": query, "newJobsForQuery": newJobsForQuery}).Info("Finished query")

	}
}

// func getJobAdsForJobNamesShuffle() {
// 	jobNames := db.GetJobNames()
// 	jobNamesSeq := u.From(jobNames, len(jobNames)) // shuffling names, this is a bit weird. should refactor
// 	jobNamesShuffled := u.Shuffle(jobNamesSeq)
// 	for _, jobName := range jobNamesShuffled {
// 		jobNameAsserted := jobName.(models.MonsterJobListModel) // casting back to original type (from shuffling type)
// 		query := url.QueryEscape(jobNameAsserted.Text)
// 		log.WithField("Job Query", query).Info("Starting new job name query") // should be in a test if this still works
// 		jobsAvailable := checkIfJobsAreAvailableForQuery(query)
// 		if !jobsAvailable {
// 			continue
// 		}
// 		continueToNextPage := true
// 		for i := 1; continueToNextPage; i++ {
// 			continueToNextPage = scrapeJobListingsFromJSON(query, i)
// 			log.Debug("Continues to next page? ", continueToNextPage)
// 		}
// 	}
// }

func main() {

	log.Info("Job collector starting")

	getJobAdsForJobNames()

}
