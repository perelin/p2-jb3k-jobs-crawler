package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MonsterJobListModel struct {
	Text      string
	MonsterID string
	Query     string
	gorm.Model
}

type MonsterJobAdModel struct {
	Title          string
	URL            string
	Location       string
	Employer       string
	MonsterJobID   string
	Query          string
	FullHTML       string
	Industry       string
	DatePosted     time.Time
	FirstEncounter time.Time
	LastEncounter  time.Time
	Active         bool
	gorm.Model
}
