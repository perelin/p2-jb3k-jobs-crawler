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

func GetAllJobs(active bool) []models.MonsterJobAdModel {
	db := initDB()
	defer db.Close()

	var jobAds []models.MonsterJobAdModel

	if os.Getenv("ENVIRONMENT") == "production" {

		if active {
			db.Select("id, title, url, monster_job_id, first_encounter, last_encounter, active").Where("active = ?", active).Find(&jobAds)
		} else {
			db.Select("id, title, url, monster_job_id, first_encounter, last_encounter, active").Find(&jobAds)
		}
	} else {
		if active {
			db.Select("id, title, url, monster_job_id, first_encounter, last_encounter, active").Where("active = ?", active).Limit(5).Find(&jobAds)
		} else {
			db.Select("id, title, url, monster_job_id, first_encounter, last_encounter, active").Limit(5).Find(&jobAds)
		}
		//db.Select("id, title, url, monster_job_id, first_encounter, last_encounter, active").Where("id = ?", 534).Limit(30).Find(&jobAds)
	}

	return jobAds
}

func GetAllJobsFull() []models.MonsterJobAdModel {
	db := initDB()
	defer db.Close()
	var jobAds []models.MonsterJobAdModel
	if os.Getenv("ENVIRONMENT") == "production" {
		db.Find(&jobAds)
	} else {
		db.Limit(5).Find(&jobAds)
	}
	return jobAds
}

func UpdateJobActiveStatus(jobID int, active bool) {
	db := initDB()
	defer db.Close()
	//db.Model(&models.MonsterJobAdModel{}).Where("id = ?", jobID).Updates(map[string]interface{}{"active": false, "last_encounter": time.Now()})
	db.Model(&models.MonsterJobAdModel{}).Where("id = ?", jobID).Updates(map[string]interface{}{"active": false})
}

func GetJobWithMonsterID(monsterID string) models.MonsterJobAdModel {
	db := initDB()
	defer db.Close()

	var jobAd models.MonsterJobAdModel

	db.Where("monster_job_id = ?", monsterID).Find(&jobAd)

	return jobAd
}

func GetJobWithMonsterIDCount(monsterID string) int {
	db := initDB()
	defer db.Close()

	var jobAd models.MonsterJobAdModel
	var count int
	db.Where("monster_job_id = ?", monsterID).Find(&jobAd).Count(&count)

	return count
}

func TouchLastEncounter(monsterID string) {
	db := initDB()
	defer db.Close()

	var jobAd models.MonsterJobAdModel

	db.Model(&jobAd).Where("monster_job_id = ?", monsterID).Update("last_encounter", time.Now)
}
