package models

type User struct {
	Id           uint   `json:"id"`
	Name         string `json:"name" gorm:"not null;unique"`
	PasswordHash string `json:"passwordHash"`
}
