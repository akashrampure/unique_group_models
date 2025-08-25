package service

import (
	"errors"
	"vds/config"
	"vds/model"

	"gorm.io/gorm/clause"
)

type ModelCount struct {
	GroupName string `gorm:"column:groupname" json:"groupname"`
	ModelName string `gorm:"column:modelname" json:"modelname"`
	Count     int    `gorm:"column:count" json:"count"`
}

// to create unique models from csv
func CreateUniqueModels(model []model.UniqueModel) error {
	if len(model) == 0 {
		return errors.New("model is empty")
	}
	batchSize := 1000
	for i := 0; i < len(model); i += batchSize {
		batch := model[i:min(i+batchSize, len(model))]
		res := config.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "groupid"}, {Name: "modelid"}},
			DoNothing: true,
		}).Omit("no")
		res.Create(&batch)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

// to get unique models from database
func GetUniqueModels() ([]model.UniqueModel, error) {
	var models []model.UniqueModel
	res := config.DB.Order("no asc").Find(&models)
	if res.Error != nil {
		return nil, res.Error
	}
	return models, nil
}

func CreateUniqueModelsCount(models []model.UniqueModelCount) error {
	if len(models) == 0 {
		return errors.New("models is empty")
	}
	batchSize := 1000

	for i := 0; i < len(models); i += batchSize {
		batch := models[i:min(i+batchSize, len(models))]
		res := config.DB.Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "groupid"}, {Name: "modelid"}},
			DoNothing: true,
		})
		res.Create(&batch)
		if res.Error != nil {
			return res.Error
		}
	}
	return nil
}

func GetUniqueModelsCount() ([]ModelCount, error) {

	query := `SELECT t1.groupname,t1.modelname,t2.count 
	FROM ams_meta.ams_unique_group_model t1
	JOIN ams_meta.ams_unique_model_count t2
	ON t1.groupid = t2.groupid
	AND t1.modelid = t2.modelid
	ORDER BY t2.count DESC
	`
	var models []ModelCount
	res := config.DB.Raw(query).Scan(&models)
	if res.Error != nil {
		return nil, res.Error
	}
	return models, nil
}
