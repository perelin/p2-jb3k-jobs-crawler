package main

import (
	"net/http"
	"strings"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func checkIfJobsAreAvailableForQuery(query string) bool {
	jobsAreAvailable := true
	jobListingURL := "https://www.monster.de/jobs/suche/?q=" + query
	res, err := http.Get(jobListingURL)
	if err != nil {
		log.Error("Couldnt check if jobs are available: ", err)
		return false
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("Couldnt check if jobs are available, returns non 200, status code error: %d %s", res.StatusCode, res.Status)
		return false
	}

	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("Couldnt check if jobs are available, couldn´t parse list page", err)
		return false
	}

	headerText := doc.Find("header.title h1.pivot").Text()

	negative := "Leider haben wir für diese Suche keinen passenden Job."

	actual := strings.TrimSpace(headerText)

	if negative == actual {
		jobsAreAvailable = false
	}

	log.WithFields(log.Fields{"query": query, "jobsAvailable": jobsAreAvailable}).Debug("Checked if jobs are available")

	return jobsAreAvailable
}
