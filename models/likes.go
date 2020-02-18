package models

import (
	"github.com/jinzhu/gorm"
	u "real-quest-v1/utils"
)

type Like struct {
	gorm.Model
	PostId uint `json:"post_Id"`
	UserId uint `json:"user_id"` //The user that this contact belongs to
}

func (like *Like) Create() map[string]interface{} {
	GetDB().Create(like)

	resp := u.Message(true, "success")
	resp["like"] = like
	return resp
}
