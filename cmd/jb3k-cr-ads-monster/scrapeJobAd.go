package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"p2lab/recruitbot3000/pkg/db"
	"p2lab/recruitbot3000/pkg/helper"
	"p2lab/recruitbot3000/pkg/models"
	"p2lab/recruitbot3000/pkg/responses"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func scrapeJobAd(jobAdEntry responses.MonsterJobAdListEntry, query string) {

	helper.DelayForMonsterAPI()

	log.Debug(jobAdEntry.JobViewURL)

	// load job ad page
	res, err := http.Get(jobAdEntry.JobViewURL)
	if err != nil {
		log.Error("Couldnt load Job Ad page: ", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("Job Ad page returns non 200, status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	// extract raw document
	bodyBytes, _ := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyText := string(bodyBytes)

	// extract parsed document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("Couldn´t parse Job Ad page", err)
		return
	}

	// extract job data
	jobJSONText := doc.Find("[type='application/ld+json']").Text()
	var jobJSON responses.JobPosting
	err = json.NewDecoder(strings.NewReader(jobJSONText)).Decode(&jobJSON)
	if err != nil {
		log.Error("Couldn´t find/parse Job data JSON: ", err)
	}

	// Data collection
	var industryList responses.IndustryList
	industryList.Industry = jobJSON.Industry
	industryString := responses.IndustryToJSON(industryList)
	header := doc.Find("div#JobViewHeader")
	title := header.Find("h1").Text()
	subtitle := header.Find("h2").Text()
	companyBox := doc.Find("div#AboutCompany")
	employer := companyBox.Find("h3.name").Text()
	trackingDiv := doc.Find("div#trackingIdentification")
	dataJobID, _ := trackingDiv.Attr("data-job-id")
	datePosted, err := time.Parse("2006-01-02T15:04", jobAdEntry.DatePosted)
	if err != nil {
		log.Debug(err)
	}

	// save to db
	adModel := models.MonsterJobAdModel{
		Title:          title,
		URL:            jobAdEntry.JobViewURL,
		Location:       subtitle,
		Employer:       employer,
		Query:          query,
		FullHTML:       bodyText,
		Active:         true,
		Industry:       industryString,
		DatePosted:     datePosted,
		FirstEncounter: time.Now(),
		LastEncounter:  time.Now(),
		MonsterJobID:   dataJobID,
	}
	db.SaveJobAdToDB(dataJobID, adModel)
	newJobs++
	newJobsForQuery++
}
