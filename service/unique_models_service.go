package service

import (
	"errors"
	"vds/config"
	"vds/model"

	"gorm.io/gorm/clause"
)

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
