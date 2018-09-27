package main

import (
	"net/url"
	"os"
	"p2lab/recruitbot3000/pkg/db"
	"p2lab/recruitbot3000/pkg/models"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
	"github.com/joho/godotenv"
	"github.com/satori/go.uuid"
	log "github.com/sirupsen/logrus"
)

var logger *log.Entry
var taskID string

type stepstone struct {
}

func (s stepstone) getJobAdListForQuery(query string) []models.MonsterJobAdModel {

	url := "https://www.stepstone.de/5/ergebnisliste.html?an=paging_next&rsearch=1&ke=" + url.QueryEscape(query)

	var jobAdURLs []models.MonsterJobAdModel

	offset := 0

	c := colly.NewCollector()

	//c := colly.NewCollector(colly.Debugger(&debug.LogDebugger{}))

	c.OnHTML("div.job-elements-list", func(e *colly.HTMLElement) {

		e.ForEach("div.job-element-row", func(i int, e2 *colly.HTMLElement) {

			detailsPageURL := e2.ChildAttr("div.job-element__body a", "href")

			title := e2.ChildText("h2.job-element__body__title")
			ssID := getStepstoneIDFromURL(detailsPageURL)

			jobAdURLs = append(jobAdURLs, models.MonsterJobAdModel{
				URL:         detailsPageURL,
				JobSource:   "stepstone",
				JobSourceID: ssID,
				Title:       title,
				Query:       query,
			})

			logger.WithFields(log.Fields{"id": ssID, "job query": query, "title": title}).Debug("found job ad")

		})

		offset = len(jobAdURLs)
		offsetURL := url + "&of=" + strconv.Itoa(offset)
		e.Request.Visit(offsetURL)
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

	c.Visit(url)

	return jobAdURLs
}

func (s stepstone) splitJobAdList(jobAds []models.MonsterJobAdModel) ([]models.MonsterJobAdModel, []models.MonsterJobAdModel) {

	var newJobAds []models.MonsterJobAdModel
	var existingJobAds []models.MonsterJobAdModel

	for _, jobAd := range jobAds {
		//sourceID := getStepstoneIDFromURL(jobAd.URL)
		//fmt.Println(jobAd.JobSourceID)
		//fmt.Println(db.IsJobInDB(jobAd.JobSourceID, jobAd.JobSource))

		if db.IsJobInDB(jobAd.JobSourceID, jobAd.JobSource) {
			existingJobAds = append(existingJobAds, jobAd)
		} else {
			newJobAds = append(newJobAds, jobAd)
		}
	}

	return newJobAds, existingJobAds
}

func (s stepstone) saveNewJobAds(jobAds []models.MonsterJobAdModel) {
	if len(jobAds) == 0 {
		return
	}
	for _, jobAd := range jobAds {
		getJobDetailsPage(jobAd.URL, jobAd.Query)
	}
	logger.WithFields(log.Fields{
		"count-total": len(jobAds),
		"query":       jobAds[1].Query}).Info(
		"saved new jobs to db")
}

func (s stepstone) updateExistingJobAds(jobAds []models.MonsterJobAdModel) {
	if len(jobAds) == 0 {
		return
	}

	rowsAffected := db.TouchLastEncounterBatch(jobAds, "stepstone")

	// for _, jobAd := range jobAds {
	// 	db.TouchLastEncounter(jobAd.JobSourceID, "stepstone")
	// }
	logger.WithFields(log.Fields{
		"count-total":  len(jobAds),
		"rows-updated": rowsAffected,
		"query":        jobAds[1].Query}).Info(
		"'last seen' timestamp of existing jobs was updated")
}

//-----------

func init() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Error("Error loading local .env file")
	}

	logLevel, _ := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logLevel)
	taskID = uuid.Must(uuid.NewV4()).String()
	logger = log.WithFields(log.Fields{"z-task": "cr-ads-stepstone", "z-task-id": taskID}) // z is used so the entry will appear at the end of the log line
}

func getJobDetailsPage(adURL string, query string) {

	var respBody []byte

	c := colly.NewCollector()

	c.OnResponse(func(r *colly.Response) {
		respBody = r.Body
	})

	c.OnHTML("html", func(e *colly.HTMLElement) {

		adModel := scrapeJobFromDetailPage(e)
		adModel.Query = query
		adModel.FullHTML = string(respBody)

		logger.WithFields(log.Fields{"ad title": adModel.Title, "query": query}).Debug("adding job to db")

		db.SaveJobAdToDB(adModel.JobSourceID, adModel)
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

func getStepstoneIDFromURL(detailsPageURL string) string {
	var ssID string
	parsedURL, _ := url.Parse(detailsPageURL)
	pathSplit := strings.Split(parsedURL.Path, "-")
	ssID = pathSplit[len(pathSplit)-2]
	return ssID
}

func scanOverJobNames() {

	ss := stepstone{}

	//jobNames := db.GetJobNames()
	jobNames := db.GetReducedJobNames()

	//jobNames := []models.MonsterJobListModel{models.MonsterJobListModel{Text: "Produktionsoptimierer/in"}}
	//jobNames := []models.MonsterJobListModel{models.MonsterJobListModel{Text: "ABAP"}}
	//jobNames := []models.MonsterJobListModel{models.MonsterJobListModel{Text: "PHP"}}

	for _, jobName := range jobNames {

		logger.WithFields(log.Fields{
			"query": jobName.Text}).Info(
			"starting job ad query crawl")

		jobAdURLs := ss.getJobAdListForQuery(jobName.Text)
		newJobAds, existingJobAds := ss.splitJobAdList(jobAdURLs)

		logger.WithFields(log.Fields{
			"count-total":    len(jobAdURLs),
			"count-new":      len(newJobAds),
			"count-existing": len(existingJobAds),
			"query":          jobName.Text}).Info(
			"finished job ad query crawl")

		db.SaveJobAdCrawlerEvent(models.JobAdCrawlerEventModel{
			JobSource:                "stepstone",
			EventTime:                time.Now(),
			JobAdResultCountTotal:    len(jobAdURLs),
			JobAdResultCountNew:      len(newJobAds),
			JobAdResultCountExisting: len(existingJobAds),
			Query:                    jobName.Text,
			TaskID:                   taskID,
		})

		ss.saveNewJobAds(newJobAds)
		ss.updateExistingJobAds(existingJobAds)
	}
}

// todo next:
// why does updating existing jobs takes so long -> try batch updates?
// why does splitting takes so long
// remove html from db

func main() {

	scanOverJobNames()

	//ss := stepstone{}
	//jobAdURLs := ss.getJobAdListForQuery("php")

	//fmt.Println(jobAdURLs)

	//spew.Dump(ss.getJobAdListForQuery("php"))

	//jobAdURLs := collyScraperTest("php")
	//fmt.Println(len(jobAdURLs))

	//scanSingleJobName("Produktionsoptimierer/in")
	//scanSingleJobName("java")

}
