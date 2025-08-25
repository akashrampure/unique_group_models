package model

type UniqueModel struct {
	GroupNo   int    `gorm:"column:groupno;" json:"groupno"`
	GroupId   string `gorm:"column:groupid" json:"groupid"`
	ModelId   string `gorm:"column:modelid" json:"modelid"`
	GroupName string `gorm:"column:groupname" json:"groupname"`
	ModelName string `gorm:"column:modelname" json:"modelname"`
}

func (UniqueModel) TableName() string {
	return "ams_meta.ams_unique_group_model"
}
