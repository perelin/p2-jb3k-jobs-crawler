package main

import (
	"encoding/json"
	"errors"
	"p2lab/recruitbot3000/pkg/db"
	"p2lab/recruitbot3000/pkg/helper"
	"p2lab/recruitbot3000/pkg/responses"
	"strconv"

	log "github.com/sirupsen/logrus"
)

func scrapeJobListingsFromJSON(query string, page int) bool {

	jobList, err := getJobListing(query, page)
	if err != nil {
		return false
	}

	for _, jobEntry := range jobList {

		if jobEntry.JobID == 0 {
			continue
		}

		jobsFoundForQuery++

		jobID := jobEntry.MusangKingID // because 'JobID' can return wrong IDs...

		jobAdCount := db.GetJobWithMonsterIDCount(jobID) // check if job already in DB
		if jobAdCount == 0 {
			log.WithFields(log.Fields{"job_id": jobID}).Debug("new job found")
			if jobEntry.JobViewURL != "" {
				scrapeJobAd(jobEntry, query)
			} else {
				log.Debug("URL is empty -> skipping ")
			}
		} else {
			//log.WithFields(log.Fields{"job_id": strconv.Itoa(jobEntry.JobID)}).Debug("job already in DB")
			db.TouchLastEncounter(jobID)
		}
	}

	log.Debug(jobsFoundForQuery, expectedResults)

	if len(jobList) == 0 || jobsFoundForQuery >= expectedResults {
		return false
	}

	return true
}

func getJobListing(query string, page int) ([]responses.MonsterJobAdListEntry, error) {
	helper.DelayForMonsterAPI()

	requestString := "https://www.monster.de/jobs/suche/pagination/?q=" + query + "&isDynamicPage=true&isMKPagination=true&page=" + strconv.Itoa(page)

	log.WithFields(log.Fields{"requestString": requestString}).Debug("getting new page")

	resp, _, errs := request.Get(requestString).End()

	if errs != nil {
		log.Error("Job Ad listing page couldn´t be loaded: ", errs)
	}
	if resp.StatusCode != 200 {
		log.Error("Job Ad listing page returns non 200, status code error: %d %s", resp.StatusCode, resp.Status)
		return nil, errors.New("sdfsdf")
	}

	var jobList []responses.MonsterJobAdListEntry

	err := json.NewDecoder(resp.Body).Decode(&jobList)
	if err != nil {
		log.Error("Couldn´t parse Job Ad listing page: ", err)
	}

	log.WithFields(log.Fields{"page": page, "result-count": len(jobList), "total-new-jobs": newJobs, "jobsFoundForQuery": jobsFoundForQuery}).Info("Received result list")

	return jobList, nil
}
