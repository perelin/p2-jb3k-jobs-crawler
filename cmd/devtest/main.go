package main

import (
	"fmt"
	"log"
	"os"
	"p2lab/recruitbot3000/pkg/db"
	"p2lab/recruitbot3000/pkg/models"
)

func printDir() {
	dir, err := os.Getwd()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(dir)
}

func profileTouchLastEncounter() {
	for i := 0; i < 30; i++ {
		db.TouchLastEncounter("200554833", "monster")
	}
}

func profileTouchLastEncounter2() {
	db.TestTouchLastEncounter("200554833", "monster", 30)
}

func profileTouchLastEncounter3() {

	var jobAds []models.MonsterJobAdModel

	for i := 0; i < 30; i++ {

		jobAds = append(jobAds, models.MonsterJobAdModel{JobSourceID: "200554833"})

		//db.TouchLastEncounter("200554833", "monster")
	}

	db.TouchLastEncounterBatch(jobAds, "monster")
}

func main() {
	profileTouchLastEncounter3()
	//db.GetLastEntryDate()
}
