package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	models "p2lab/recruitbot3000/pkg/models"

	"github.com/PuerkitoBio/goquery"
	u "github.com/alxrm/ugo"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
)

type monsterJobAdListResults []struct {
	JobID                 int    `json:"JobID"`
	Title                 string `json:"Title"`
	TitleLink             string `json:"TitleLink"`
	IsBolded              bool   `json:"IsBolded"`
	DatePostedText        string `json:"DatePostedText"`
	DatePosted            string `json:"DatePosted"`
	LocationText          string `json:"LocationText"`
	LocationLink          string `json:"LocationLink"`
	JobViewURL            string `json:"JobViewUrl"`
	ImpressionTracking    string `json:"ImpressionTracking"`
	HasLocationAddress    bool   `json:"HasLocationAddress"`
	IsSavedJob            bool   `json:"IsSavedJob"`
	IsAppliedJob          bool   `json:"IsAppliedJob"`
	IsNewJob              bool   `json:"IsNewJob"`
	HasAdapt              bool   `json:"HasAdapt"`
	HasProDiversity       bool   `json:"HasProDiversity"`
	HasSpecialCommitments bool   `json:"HasSpecialCommitments"`
	Company               struct {
		Name              string `json:"Name"`
		CompanyLink       string `json:"CompanyLink"`
		HasCompanyAddress bool   `json:"HasCompanyAddress"`
	} `json:"Company"`
	Text                       string      `json:"Text"`
	LocationClickJsFunction    string      `json:"LocationClickJsFunction"`
	CompanyClickJsFunction     string      `json:"CompanyClickJsFunction"`
	JobTitleClickJsFunction    string      `json:"JobTitleClickJsFunction"`
	JobDescription             string      `json:"JobDescription"`
	ApplyMethod                int         `json:"ApplyMethod"`
	ApplyType                  string      `json:"ApplyType"`
	IsAggregated               string      `json:"IsAggregated"`
	CityText                   string      `json:"CityText"`
	StateText                  string      `json:"StateText"`
	JobDescriptionMeta         string      `json:"JobDescriptionMeta"`
	EmploymentTypeMeta         string      `json:"EmploymentTypeMeta"`
	IndustryTypeMeta           string      `json:"IndustryTypeMeta"`
	JobViewURLMeta             string      `json:"JobViewUrlMeta"`
	IsFastApply                bool        `json:"IsFastApply"`
	Target                     interface{} `json:"Target"`
	IsSecondaryJob             bool        `json:"IsSecondaryJob"`
	JobIDCloud                 int         `json:"JobIdCloud"`
	MusangKingID               string      `json:"MusangKingId"`
	IsSecondarySearchResultJob bool        `json:"IsSecondarySearchResultJob"`
	InlineAdIndex              int         `json:"InlineAdIndex"`
	ShowCompanyAsLink          bool        `json:"ShowCompanyAsLink"`
	ShowLocationAsLink         bool        `json:"ShowLocationAsLink"`
	HideCompanyLogo            bool        `json:"HideCompanyLogo"`
	ShowMultilocHover          bool        `json:"ShowMultilocHover"`
	MultilocHoverTitle         interface{} `json:"MultilocHoverTitle"`
	MultilocHover              interface{} `json:"MultilocHover"`
}

var request = gorequest.New()

func delayForMonsterAPI() {
	log.Debug("Waiting delay for api throtteling...")
	delay, _ := strconv.Atoi(os.Getenv("DELAY"))
	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func scrapeJobListingsFromJSON(query string, page int) bool {

	delayForMonsterAPI()

	requestString := "https://www.monster.de/jobs/suche/pagination/?q=" + query + "&isDynamicPage=true&isMKPagination=true&page=" + strconv.Itoa(page)

	log.WithField("url", requestString).Info("Request URL for Job Ad listing")

	resp, _, errs := request.Get(requestString).End()

	if errs != nil {
		log.Error("Job Ad listing page couldn´t be loaded: ", errs)
	}
	if resp.StatusCode != 200 {
		log.Error("Job Ad listing page returns non 200, status code error: %d %s", resp.StatusCode, resp.Status)
		return false
	}

	var jobList monsterJobAdListResults

	err := json.NewDecoder(resp.Body).Decode(&jobList)
	if err != nil {
		log.Error("Couldn´t parse Job Ad listing page: ", err)
	}

	log.WithFields(log.Fields{"page": page, "result count": len(jobList)}).Info("Received result list")

	for _, jobEntry := range jobList {
		if jobEntry.JobViewURL != "" {
			scrapeJobAd(jobEntry.JobViewURL, query)
		}
	}

	if len(jobList) == 0 || len(jobList) < 26 {
		return false
	}

	return true
}

func scrapeJobAd(linkURL string, query string) {

	delayForMonsterAPI()

	log.Debug(linkURL)

	res, err := http.Get(linkURL)
	if err != nil {
		log.Error("Couldnt load Job Ad page: ", err)
		return
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		log.Errorf("Job Ad page returns non 200, status code error: %d %s", res.StatusCode, res.Status)
		return
	}

	bodyBytes, _ := ioutil.ReadAll(res.Body)
	res.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	bodyText := string(bodyBytes)

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(res.Body)
	if err != nil {
		log.Error("Couldn´t parse Job Ad page", err)
		return
	}

	header := doc.Find("div#JobViewHeader")
	title := header.Find("h1").Text()
	subtitle := header.Find("h2").Text()
	//fmt.Println(title)
	//fmt.Println(subtitle)

	companyBox := doc.Find("div#AboutCompany")
	employer := companyBox.Find("h3.name").Text()
	//fmt.Println(employer)

	trackingDiv := doc.Find("div#trackingIdentification")
	dataJobID, _ := trackingDiv.Attr("data-job-id")

	adModel := models.MonsterJobAdModel{
		Title:          title,
		URL:            linkURL,
		Location:       subtitle,
		Employer:       employer,
		Query:          query,
		FullHTML:       bodyText,
		Active:         true,
		FirstEncounter: time.Now(),
		LastEncounter:  time.Now(),
		MonsterJobID:   dataJobID,
	}
	saveJobAdToDB(dataJobID, adModel)
}

func saveJobAdToDB(dataJobID string, jobModel models.MonsterJobAdModel) {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("failed to connect database", err)
		panic("failed to connect database")
	}
	defer db.Close()

	db.Where(models.MonsterJobAdModel{
		MonsterJobID: dataJobID,
	}).FirstOrCreate(&jobModel)
}

func init() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Error("Error loading local .env file")
	}

	// Log as JSON instead of the default ASCII formatter.
	//log.SetFormatter(&log.JSONFormatter{})

	// Output to stdout instead of the default stderr
	// Can be any io.Writer, see below for File example
	//log.SetOutput(os.Stdout)

	// Only log the warning severity or above.
	//log.SetLevel(log.DebugLevel)
	logLevel, _ := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logLevel)

}

func getJobNames() []models.MonsterJobListModel {

	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Println("failed to connect database", err)
		panic("failed to connect database")
	}
	defer db.Close()

	var jobNames []models.MonsterJobListModel

	db.Find(&jobNames)

	//fmt.Println(jobNames)

	return jobNames
	// query := url.QueryEscape("Analyst Beschaffung")
	// fmt.Println(query)
}

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

func getJobAdsForJobNames() {
	jobNames := getJobNames()
	jobNamesSeq := u.From(jobNames, len(jobNames)) // shuffling names, this is a bit weird. should refactor
	jobNamesShuffled := u.Shuffle(jobNamesSeq)
	for _, jobName := range jobNamesShuffled {
		jobNameAsserted := jobName.(models.MonsterJobListModel) // casting back to original type (from shuffling type)
		query := url.QueryEscape(jobNameAsserted.Text)
		log.WithField("Job Query", query).Info("Starting new job name query") // should be in a test if this still works
		jobsAvailable := checkIfJobsAreAvailableForQuery(query)
		if !jobsAvailable {
			continue
		}
		continueToNextPage := true
		for i := 1; continueToNextPage; i++ {
			continueToNextPage = scrapeJobListingsFromJSON(query, i)
			log.Debug("Continues to next page? ", continueToNextPage)
		}
	}
}

func main() {

	log.Info("Job collector starting")

	getJobAdsForJobNames()

}
