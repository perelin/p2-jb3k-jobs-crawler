package main

import (
	"fmt"
	"os"
	db "p2lab/recruitbot3000/pkg/db"
	"p2lab/recruitbot3000/pkg/models"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/davecgh/go-spew/spew"

	"github.com/joho/godotenv"
	"github.com/parnurzeal/gorequest"
	log "github.com/sirupsen/logrus"
)

type JobPosting struct {
	Type        string   `json:"@type"`
	Context     string   `json:"@context"`
	Title       string   `json:"title"`
	DatePosted  string   `json:"datePosted"`
	Description string   `json:"description"`
	Industry    []string `json:"industry"`
	JobLocation struct {
		Type string `json:"@type"`
		Geo  struct {
			Type      string `json:"@type"`
			Latitude  string `json:"latitude"`
			Longitude string `json:"longitude"`
		} `json:"geo"`
		Address struct {
			Type            string `json:"@type"`
			AddressLocality string `json:"addressLocality"`
			AddressRegion   string `json:"addressRegion"`
			PostalCode      string `json:"postalCode"`
			AddressCountry  string `json:"addressCountry"`
		} `json:"address"`
	} `json:"jobLocation"`
	URL                string `json:"url"`
	HiringOrganization struct {
		Type string `json:"@type"`
		Name string `json:"name"`
		Logo string `json:"logo"`
	} `json:"hiringOrganization"`
	ValidThrough   string `json:"validThrough"`
	SalaryCurrency string `json:"salaryCurrency"`
	Identifier     struct {
		Type  string `json:"@type"`
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"identifier"`
}

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

func extractDatePosted(jobAd models.MonsterJobAdModel) {
	//spew.Dump(jobAd.FullHTML)
	// err := ioutil.WriteFile("tmp/"+strconv.Itoa(int(jobAd.ID))+".html", []byte(jobAd.FullHTML), 0644)
	// if err != nil {
	// 	log.Debug("couldn´t save doc: " + strconv.Itoa(int(jobAd.ID)))
	// }
	htmlReader := strings.NewReader(jobAd.FullHTML)
	doc, err := goquery.NewDocumentFromReader(htmlReader)
	if err != nil {
		log.Debug("couldn´t convert html to doc: " + strconv.Itoa(int(jobAd.ID)))
	}

	title := doc.Find("title").Text()

	//spew.Dump(reflect.TypeOf(doc))
	spew.Dump(title)
}

func main() {
	lastEntryTime := db.GetLastEntryDate()
	fmt.Println(lastEntryTime)
	jobAds := db.GetAllJobsFull()

	log.WithFields(log.Fields{"count": len(jobAds)}).Info("scanning jobs")

	jobs404 := 0
	jobsProcessed := 0

	for i, jobAd := range jobAds {
		delayForMonsterAPI()

		log.WithFields(log.Fields{"monsterID": jobAd.JobSourceID, "running no": strconv.Itoa(i)}).Debug("processing new job")

		extractDatePosted(jobAd)

		// resp, _, errs := request.Get(jobAd.URL).End()

		// if errs != nil {
		// 	log.Error("Job Ad page couldn´t be loaded: ", errs)
		// 	continue
		// }

		// jobsProcessed++

		// if resp.StatusCode == 404 {
		// 	log.WithFields(log.Fields{"url": jobAd.URL}).Debug("job ad page returns 404, job ad might no longer be active")
		// 	//db.UpdateJobActiveStatus(int(jobAd.ID), false)
		// 	jobs404++
		// } else if resp.StatusCode == 200 {
		// 	log.WithFields(log.Fields{"url": jobAd.URL}).Debug("job ad seems to be alive")
		// 	extractDatePosted(jobAd)
		// }
	}

	log.WithFields(log.Fields{"total-new-inactive": jobs404, "total-proccesed": jobsProcessed}).Info("finished scan")
}
