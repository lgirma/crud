package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/lgirma/crud"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var Db *gorm.DB

func Init() {
	_db, err := gorm.Open(sqlite.Open("addressbook.db"), &gorm.Config{})
	Db = _db
	if err != nil {
		log.Fatal("Failed to connect to database")
	}

	Db.AutoMigrate(
		&Tag{},
		&Contact{},
	)
}

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	Init()

	r := gin.Default()
	r.Use(CORSMiddleware())

	contactsRepo := crud.NewCrudService(
		Db,
		func(e Contact) string { return e.PublicId },
		func(t *Contact, a string) { t.PublicId = a },
		&crud.CrudServiceOptions[Contact, string]{
			LookupQuery: "full_name like ? or email like ?",
		},
	)

	tagsRepo := crud.NewCrudService(
		Db,
		func(e Tag) string { return e.PublicId },
		func(t *Tag, a string) { t.PublicId = a },
		&crud.CrudServiceOptions[Tag, string]{
			LookupQuery: "name like ?",
		},
	)

	crud.AddCrudGinRestApi[Contact, string]("api/contacts", r, contactsRepo, &crud.CrudRestApiOptions[Contact, string]{})
	crud.AddCrudGinRestApi[Tag, string]("api/tags", r, tagsRepo, &crud.CrudRestApiOptions[Tag, string]{})

	r.Run(":5050")
}
