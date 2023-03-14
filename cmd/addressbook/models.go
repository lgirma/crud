package main

type Contact struct {
	Id       int `json:"-"`
	FullName string `json:"full_name"`
	PublicId string `gorm:"index:idx_contacts_public_id,unique" json:"public_id"`
	Email    string `json:"email"`
	Phone    string `json:"phone"`
	Address  string `json:"address"`
}

type Tag struct {
	Id       int `json:"-"`
	PublicId string `gorm:"index:idx_tags_public_id,unique"`
	Name     string
}
