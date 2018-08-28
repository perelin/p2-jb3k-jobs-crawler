package db

import (
	"os"
	"p2lab/recruitbot3000/pkg/models"
	"time"

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

func initDB() *gorm.DB {
	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	//db.LogMode(true)
	if err != nil {
		log.Println("failed to connect database", err)
		panic("failed to connect database")
	}
	return db
}

func GetLastEntryDate() time.Time {
	db := initDB()
	defer db.Close()

	var lastJobAd models.MonsterJobAdModel

	db.Last(&lastJobAd)

	//fmt.Println(lastJobAd.CreatedAt)

	return lastJobAd.CreatedAt
}

func GetJobAdCount() int {
	db := initDB()
	defer db.Close()

	var count int

	db.Model(&models.MonsterJobAdModel{}).Count(&count)

	return count
}

func GetAllJobs() []models.MonsterJobAdModel {
	db := initDB()
	defer db.Close()

	var jobAds []models.MonsterJobAdModel

	db.Select("title, url, monster_job_id, first_encounter, last_encounter, active").Limit(3).Find(&jobAds)

	return jobAds
}
