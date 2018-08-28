package main

import (
	"fmt"
	db "p2lab/recruitbot3000/pkg/db"
)

func main() {
	lastEntryTime := db.GetLastEntryDate()
	fmt.Println(lastEntryTime)
	jobAds := db.GetAllJobs()
	fmt.Println(jobAds)

	// walk over ever job
	// - check ifjob page returns 404
	// -- set active to false
}
