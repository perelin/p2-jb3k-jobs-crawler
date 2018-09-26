package base

import (
	"p2lab/recruitbot3000/pkg/models"
)

type JobAdPlatform interface {
	getJobAdListForQuery(query string) []models.MonsterJobAdModel

	splitJobAdList([]models.MonsterJobAdModel) ([]models.MonsterJobAdModel, []models.MonsterJobAdModel)

	saveNewJobAds([]models.MonsterJobAdModel)

	updateExistingJobAds([]models.MonsterJobAdModel)
}
