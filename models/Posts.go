package models

import "github.com/jinzhu/gorm"

type Post struct {
	gorm.Model
	Title    string `json:"title"`
	Category string `json:"category"`
	Likes    int    `json:"likes"`
	Tasks    string `json:"tasks"`
}
