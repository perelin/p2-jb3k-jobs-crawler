package main

import (
	"net/http"
	"p2lab/recruitbot3000/pkg/helper"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	log "github.com/sirupsen/logrus"
)

func GetTotalResultsForQuery(query string) int {

	helper.DelayForMonsterAPI()

	res, err := http.Get("https://www.monster.de/jobs/suche/?q=" + query)
	if err != nil {
		log.Error("Couldnt load job ad result listing html page: ", err)
		return 0
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("job ad result listing html page returns non 200, status code error: %d %s", res.StatusCode, res.Status)
		return 0
	}

	// extract parsed document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("Couldn´t parse job ad result listing html page", err)
		return 0
	}

	resultDiv := doc.Find("div.mux-search-results")
	resultcountString, _ := resultDiv.Attr("data-results-total")
	resultcount, _ := strconv.Atoi(resultcountString)
	return resultcount
}
