package main

type Contact struct {
	Id       int `json:"-"`
	FullName string
	PublicId string `gorm:"index:idx_contacts_public_id,unique"`
	Email    string
	Phone    string
	Address  string
}

type Tag struct {
	Id       int `json:"-"`
	PublicId string `gorm:"index:idx_tags_public_id,unique"`
	Name     string
}
