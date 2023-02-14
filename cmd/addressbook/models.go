package main

type Contact struct {
	Id       int `json:"-"`
	FullName string
	PublicId string `gorm:"unique"`
	Email    string
	Phone    string
	Address  string
}

type Tag struct {
	Id       int
	PublicId string `gorm:"unique"`
	Name     string
}
