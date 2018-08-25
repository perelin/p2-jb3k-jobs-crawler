// You can edit this code!
// Click here and start typing.
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	"github.com/parnurzeal/gorequest"
)

type monsterJobListResults struct {
	D struct {
		Type   string `json:"__type"`
		Result struct {
			Index       interface{} `json:"Index"`
			TooManyData bool        `json:"TooManyData"`
			Items       []struct {
				Text string      `json:"Text"`
				ID   int         `json:"ID"`
				Data interface{} `json:"Data"`
			} `json:"Items"`
		} `json:"Result"`
	} `json:"d"`
}

type monsterJobListModel struct {
	Text      string
	MonsterID string
	Query     string
	gorm.Model
}

var request = gorequest.New()

func buildQueryURLArbeitgeber(queryString string) string {

	query := queryString
	maxResults := "1000"
	searchType := "132" // 132
	searchFlags := "1"  // 419

	fullQueryURL := "https://arbeitgeber.monster.de/SharedUI/Services/AutoComplete.asmx/GetCompletionList?request=%7B%22Query%22%3A%22" + query + "%22%2C%22MaxResults%22%3A" + maxResults + "%2C%22SearchType%22%3A" + searchType + "%2C%22SearchFlags%22%3A" + searchFlags + "%7D"

	return fullQueryURL

}

func buildQueryURLStartpage(queryString string) string {

	query := queryString
	maxResults := "1000"
	searchType := "132" // 132
	channelID := "419"  // 419

	fullQueryURL := "https://www.monster.de/Services/MAutoComplete.asmx/GetCompletionList?request=%7B%22Query%22%3A%22" + query + "%22%2C%22MaxResults%22%3A" + maxResults + "%2C%22SearchType%22%3A" + searchType + "%2C%22ChannelId%22%3A" + channelID + "%7D"

	return fullQueryURL

}

func runjobListingQuery(queryString string) {
	//request := gorequest.New()

	getRequestString := buildQueryURLArbeitgeber(queryString)

	//fmt.Println(query)
	//fmt.Println(getRequestString)

	resp, _, errs := request.Get(getRequestString).Set("Content-Type", "application/json").End()
	if errs != nil {
		fmt.Println(getRequestString)
		//fmt.Sprintf("%+v", errs)
		fmt.Println("%+v", errs)
	}

	// fmt.Println("body")
	//fmt.Println(resp)
	//fmt.Println(body)

	var jobLists monsterJobListResults

	// if err := json.Unmarshal([]byte(body), &jobLists);
	err := json.NewDecoder(resp.Body).Decode(&jobLists)
	if err != nil {
		// fmt.Println("err ")
		// fmt.Println(body)
		fmt.Println(err)
	}

	//fmt.Println(jobLists.D.Result.Items)
	//fmt.Println(len(jobLists.D.Result.Items))

	//fmt.Println("whatup!")
	//fmt.Println(os.Getenv("DB_URL"))

	db, err := gorm.Open("postgres", os.Getenv("DB_URL"))
	if err != nil {
		log.Println(err)
		panic("failed to connect database")
	}
	defer db.Close()

	db.AutoMigrate(&monsterJobListModel{})

	for _, job := range jobLists.D.Result.Items {

		//fmt.Println(job.Text)

		jobModel := monsterJobListModel{
			Text:      job.Text,
			MonsterID: fmt.Sprintf("%v", job.ID),
			Query:     queryString,
		}

		db.Where(monsterJobListModel{
			MonsterID: fmt.Sprintf("%v", job.ID),
		}).FirstOrCreate(&jobModel)

		//db.Create(&jobModel)
	}

	fmt.Println(len(jobLists.D.Result.Items))
}

func iterativeQuery() {

	var delay time.Duration = 100

	sem := make(chan struct{}, 32)
	var wg sync.WaitGroup

	for i := 65; i <= 90; i++ {
		time.Sleep(delay * time.Millisecond)
		char := byte(i)
		//fmt.Println(char)
		//fmt.Println(string(char))
		wg.Add(1)
		go func(queryString string) {
			sem <- struct{}{}
			defer func() { <-sem }()
			defer wg.Done()
			runjobListingQuery(queryString)
		}(string(char))

		for j := 65; j <= 90; j++ {
			time.Sleep(delay * time.Millisecond)
			char2 := byte(j)
			fmt.Println(string(char) + string(char2))
			wg.Add(1)
			go func(queryString2 string) {
				sem <- struct{}{}
				defer func() { <-sem }()
				defer wg.Done()
				runjobListingQuery(queryString2)
			}(string(char) + string(char2))

			for k := 65; k <= 90; k++ {
				time.Sleep(delay * time.Millisecond)
				char3 := byte(k)
				fmt.Println(string(char) + string(char2) + string(char3))
				wg.Add(1)
				go func(queryString3 string) {
					sem <- struct{}{}
					defer func() { <-sem }()
					defer wg.Done()
					runjobListingQuery(queryString3)
				}(string(char) + string(char2) + string(char3))

				//go runjobListingQuery(string(char) + string(char2))

			}

		}

	}
	wg.Wait()
}

func linearQuery(digits int, delayMiliseconds time.Duration) {

	for i := 65; i <= 90; i++ {
		time.Sleep(delayMiliseconds * time.Millisecond)
		char := byte(i)
		fmt.Println(string(char))
		runjobListingQuery(string(char))

		if digits > 1 {
			for j := 65; j <= 90; j++ {
				time.Sleep(delayMiliseconds * time.Millisecond)
				char2 := byte(j)
				fmt.Println(string(char) + string(char2))
				runjobListingQuery(string(char) + string(char2))

				if digits > 2 {
					for k := 65; k <= 90; k++ {
						time.Sleep(delayMiliseconds * time.Millisecond)
						char3 := byte(k)
						fmt.Println(string(char) + string(char2) + string(char3))
						runjobListingQuery(string(char) + string(char2) + string(char3))

						//go runjobListingQuery(string(char) + string(char2))

					}
				}

			}
		}

	}
}

func main() {

	err := godotenv.Load()
	if err != nil {
		fmt.Println("Error loading .env file")
	}

	//iterativeQuery()
	linearQuery(3, 500)
	//runjobListingQuery("A")

}
