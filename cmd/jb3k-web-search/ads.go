package main

import (
	"fmt"
	"os"

	"p2lab/recruitbot3000/pkg/models"

	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"
)

func init() {

	err := godotenv.Load("../../.env")
	if err != nil {
		log.Error("Error loading local .env file")
	}

	logLevel, _ := log.ParseLevel(os.Getenv("LOG_LEVEL"))
	log.SetLevel(logLevel)

}

func GetJobAdCount() int {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	//db.LogMode(true)
	if err != nil {
		log.Println("failed to connect database", err)
		panic("failed to connect database")
	}
	defer db.Close()

	var jobAds []models.MonsterJobAdModel

	db.Select([]string{"title", "url", "monster_job_id"}).Find(&jobAds)

	fmt.Println(len(jobAds))

	return len(jobAds)
}
