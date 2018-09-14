package main

import (
	"fmt"
	"net/url"
	"os"
	"p2lab/recruitbot3000/pkg/db"
	"p2lab/recruitbot3000/pkg/models"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry

func init() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Error("Error loading local .env file")
	}

	logLevel, _ := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logLevel)
	logger = log.WithFields(log.Fields{"task": "cr-ads-stepstone"})
}

func collyDemo() {
	// Instantiate default collector
	c := colly.NewCollector(
		// Visit only domains: hackerspaces.org, wiki.hackerspaces.org
		colly.AllowedDomains("hackerspaces.org", "wiki.hackerspaces.org"),
	)

	// On every a element which has href attribute call callback
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		link := e.Attr("href")
		// Print link
		fmt.Printf("Link found: %q -> %s\n", e.Text, link)
		// Visit link found on page
		// Only those links are visited which are in AllowedDomains
		c.Visit(e.Request.AbsoluteURL(link))
	})

	// Before making a request print "Visiting ..."
	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL.String())
	})

	// Start scraping on https://hackerspaces.org
	c.Visit("https://hackerspaces.org/")
}

func getJobDetails(adURL string, query string) {

	var respBody []byte

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		respBody = r.Body
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {

		adModel := scrapeJobFromDetailPage(e)
		adModel.Query = query
		adModel.FullHTML = string(respBody)

		db.SaveJobAdToDB(adModel.JobSourceID, adModel)

		logger.WithFields(log.Fields{"ad title": adModel.Title}).Debug("found job ad")

	})

	c.Visit(adURL)
}

func scrapeJobFromDetailPage(e *colly.HTMLElement) models.MonsterJobAdModel {
	employer := e.ChildText("h6.listing__company-name a.at-listing-nav-company-name-link")

	title := e.ChildText("h1.listing__job-title")

	location := e.ChildText("li.at-listing__list-icons_location")

	industries, _ := e.DOM.Find("div.js-listing-container-right div.js-company-content-card").Attr("data-sectors")

	canonicalURL, _ := e.DOM.Find(`link[rel="canonical"]`).Attr("href")

	alternateURL, _ := e.DOM.Find(`link[rel="alternate"]`).Attr("href")
	dataJobID := getStepstoneIDFromAlternateURL(alternateURL)

	datePostedString, _ := e.DOM.Find("div.listing-header span.date-time-ago").Attr("data-date")

	datePosted, err := time.Parse("2006-01-02 15:04:05", datePostedString)
	if err != nil {
		logger.Debug(err)
	}

	// save to db
	adModel := models.MonsterJobAdModel{
		Title:          title,
		URL:            canonicalURL,
		Location:       location,
		Employer:       employer,
		Active:         true,
		Industry:       industries,
		DatePosted:     datePosted,
		FirstEncounter: time.Now(),
		LastEncounter:  time.Now(),
		JobSourceID:    dataJobID,
		JobSource:      "stepstone",
	}

	return adModel
}

func getStepstoneIDFromAlternateURL(url string) string {
	var ssID string
	splitString := strings.Split(url, "/")
	ssID = splitString[7]
	return ssID
}

func getResultList(listURL string, query string) int {

	resultCount := 0

	c := colly.NewCollector()
	//c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	c.OnHTML("div.job-element-row", func(e *colly.HTMLElement) {
		resultCount++
		getJobDetails(e.ChildAttr("div.job-element__body a", "href"), query)
	})
	c.OnResponse(func(r *colly.Response) {
		//fmt.Println("Visited", r.Body)
	})
	c.OnRequest(func(r *colly.Request) {
		logger.WithFields(log.Fields{"url": r.URL.String()}).Debug("crawling new result page")
	})

	c.OnScraped(func(r *colly.Response) {
		//fmt.Println("Finished", r.Request.URL)
	})

	c.Visit(listURL)

	return resultCount
}

func scanSingleJobName(query string) {

	logger.WithFields(log.Fields{"job query": query}).Info("starting new job name query")

	offset := 0
	queryEscaped := url.QueryEscape(query)
	resultCount := 1

	for resultCount != 0 {

		//url := "https://www.stepstone.de/5/ergebnisliste.html?ke=" + queryEscaped
		url := "https://www.stepstone.de/5/ergebnisliste.html?&rsearch=1&of=" + strconv.Itoa(offset) + "&ke=" + queryEscaped

		resultCount = getResultList(url, queryEscaped)

		offset = offset + resultCount

		logger.WithFields(log.Fields{"new offset": offset, "results found": resultCount}).Debug("finished result page")

		//fmt.Println(resultCount)

		//break
	}

	logger.WithFields(log.Fields{"job query": query, "total results": offset}).Info("finished job name query")
}

func scanOverJobNames() {
	jobNames := db.GetJobNames()

	for _, jobName := range jobNames {

		query := url.QueryEscape(jobName.Text)
		scanSingleJobName(query)
	}
}

// next up:
// /count hits and increase offset
// /break after 0 results
// get Stepstone Job Names (and other stuff)
// /adapt Databases
// save results
// check for delays

func main() {
	scanOverJobNames()
	//scanSingleJobName("Produktionsoptimierer/in")
	//scanSingleJobName("java")

}
