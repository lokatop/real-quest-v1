package models

import (
	"github.com/jinzhu/gorm"
	u "real-quest-v1/utils"
)

type Like struct {
	gorm.Model
	QuestId      int `json:"QuestId"`
	AccountRefer uint
}

func (like *Like) Create() (map[string]interface{}) {
	GetDB().Create(like)
	
	resp:= u.Message(true,"success")
	resp["like"] = like
	return resp
}

func getLikes(user uint) ([]*Like)  {
	likes := make([]*Like,0)
	err := GetDB().Table("likes").Where("")
}