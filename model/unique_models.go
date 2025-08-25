package model

type UniqueModel struct {
	No        int    `gorm:"column:no" json:"no"`
	GroupId   string `gorm:"column:groupid" json:"groupid"`
	ModelId   string `gorm:"column:modelid" json:"modelid"`
	GroupName string `gorm:"column:groupname" json:"groupname"`
	ModelName string `gorm:"column:modelname" json:"modelname"`
}

func (UniqueModel) TableName() string {
	return "ams_meta.ams_unique_group_model"
}
