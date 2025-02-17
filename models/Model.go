package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Fullname string `json:"fullname"`;
	Email string `json:"email" gorm:"unique"`;
	Password string `json:"password"`;
	Books []Book `gorm:"foreignkey:AuthorID"`;
}

type Book struct {
	gorm.Model
	Title string `json:"title"`;
	BookURL string `json:"bookurl"`;
	BookSlug string `json:"bookslug" gorm:"unique"`;
	Description string `json:"description"`;
	Picture string `json:"picture"`;
	AuthorID uint `gorm:"index"`;
	Author User `gorm:"foreignKey:AuthorID;constraint:OnDelete:CASCADE;"` 
}