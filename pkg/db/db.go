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

func GetLastEntryDate(source string) time.Time {
	db := initDB()
	defer db.Close()

	var lastJobAd models.MonsterJobAdModel

	db.Where("job_source = ?", source).Last(&lastJobAd)

	//fmt.Println(lastJobAd.CreatedAt)

	return lastJobAd.CreatedAt
}

func GetJobAdCount(source string) int {
	db := initDB()
	defer db.Close()

	var count int

	db.Model(&models.MonsterJobAdModel{}).Where("job_source = ?", source).Count(&count)

	return count
}

func GetAllJobs(active bool) []models.MonsterJobAdModel {
	db := initDB()
	defer db.Close()

	var jobAds []models.MonsterJobAdModel

	if os.Getenv("ENVIRONMENT") == "production" {

		if active {
			db.Select("id, title, url, job_source_id, first_encounter, last_encounter, active").Where("active = ?", active).Find(&jobAds)
		} else {
			db.Select("id, title, url, job_source_id, first_encounter, last_encounter, active").Find(&jobAds)
		}
	} else {
		if active {
			db.Select("id, title, url, job_source_id, first_encounter, last_encounter, active").Where("active = ?", active).Limit(5).Find(&jobAds)
		} else {
			db.Select("id, title, url, job_source_id, first_encounter, last_encounter, active").Limit(5).Find(&jobAds)
		}
		//db.Select("id, title, url, job_source_id, first_encounter, last_encounter, active").Where("id = ?", 534).Limit(30).Find(&jobAds)
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

	db.Where("job_source_id = ?", monsterID).Find(&jobAd)

	return jobAd
}

func GetJobWithMonsterIDCount(monsterID string) int {
	db := initDB()
	defer db.Close()

	var jobAd models.MonsterJobAdModel
	var count int
	db.Where("job_source_id = ?", monsterID).Find(&jobAd).Count(&count)

	return count
}

func TouchLastEncounter(dataJobID string, source string) {
	db := initDB()
	defer db.Close()

	var jobAd models.MonsterJobAdModel

	//fmt.Println(time.Now())

	db.Model(&jobAd).Where("job_source_id = ? AND job_source = ?", dataJobID, source).Update("last_encounter", time.Now())
}

func TestTouchLastEncounter(dataJobID string, source string, runs int) {
	db := initDB()
	defer db.Close()

	var jobAd models.MonsterJobAdModel

	//fmt.Println(time.Now())
	for i := 0; i < runs; i++ {
		db.Model(&jobAd).Where("job_source_id = ? AND job_source = ?", dataJobID, source).Update("last_encounter", time.Now())
	}
}

func TouchLastEncounterBatch(jobModels []models.MonsterJobAdModel, source string) int64 {
	db := initDB()
	defer db.Close()

	var jobAd models.MonsterJobAdModel
	var dataJobIDs []string

	for _, jobAd := range jobModels {
		dataJobIDs = append(dataJobIDs, jobAd.JobSourceID)
	}

	//fmt.Println(dataJobIDs)
	//fmt.Println(source)

	//db.Model(&jobAd).Where("job_source_id IN (?) AND job_source = ?", dataJobIDs, source).Update("last_encounter", time.Now())
	rowsAffected := db.Model(&jobAd).Where("job_source_id IN (?) AND job_source = ?", dataJobIDs, source).Updates(map[string]interface{}{"last_encounter": time.Now()}).RowsAffected
	//fmt.Println(rowsAffected)

	return rowsAffected
}

func GetJobNames() []models.MonsterJobListModel {

	db := initDB()
	defer db.Close()

	var jobNames []models.MonsterJobListModel

	db.Find(&jobNames)

	//fmt.Println(jobNames)

	return jobNames
	// query := url.QueryEscape("Analyst Beschaffung")
	// fmt.Println(query)
}

func GetReducedJobNames() []models.MonsterJobListModel {

	var jobNames []models.MonsterJobListModel
	tiobeList := []string{
		"Java",
		"C",
		"Python",
		"C++",
		"Visual Basic .NET",
		"C#",
		"PHP",
		"JavaScript",
		"SQL",
		"Objective-C",
		"Delphi/Object Pascal",
		"Ruby",
		"MATLAB",
		"Assembly language",
		"Swift",
		"Go",
		"Golang",
		"Perl",
		"R",
		"PL/SQL",
		"Visual Basic",
	}

	for _, jobName := range tiobeList {
		jobNames = append(jobNames, models.MonsterJobListModel{Text: jobName})
	}

	//fmt.Println(jobNames)

	return jobNames
	// query := url.QueryEscape("Analyst Beschaffung")
	// fmt.Println(query)
}

func SaveJobAdToDB(dataJobID string, jobModel models.MonsterJobAdModel) {
	db := initDB()
	defer db.Close()

	db.Where(models.MonsterJobAdModel{
		JobSourceID: dataJobID,
	}).FirstOrCreate(&jobModel)

}

func IsJobInDB(dataJobID string, source string) bool {
	isJobInDB := false

	db := initDB()
	defer db.Close()

	var jobAd models.MonsterJobAdModel

	db.Where("job_source_id = ? AND job_source = ?", dataJobID, source).Find(&jobAd)

	if (jobAd != models.MonsterJobAdModel{}) {
		isJobInDB = true
	}

	return isJobInDB
}

func SaveJobAdCrawlerEvent(jobAdCrawlerEvent models.JobAdCrawlerEventModel) {
	db := initDB()
	defer db.Close()

	db.Create(&jobAdCrawlerEvent)

}
