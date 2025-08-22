package model

type UniqueModel struct {
	GroupId   string `gorm:"column:groupid" json:"groupid"`
	ModelId   string `gorm:"column:modelid" json:"modelid"`
	GroupName string `gorm:"column:groupname" json:"groupname"`
	ModelName string `gorm:"column:modelname" json:"modelname"`
}

func TableName() string {
	return "dmt.unique_models"
}
