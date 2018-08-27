package models

import (
	"time"

	"github.com/jinzhu/gorm"
)

type MonsterJobAdModel struct {
	Title          string
	URL            string
	Location       string
	Employer       string
	MonsterJobID   string
	FullHTML       string
	FirstEncounter time.Time
	LastEncounter  time.Time
	Active         bool
	gorm.Model
}
