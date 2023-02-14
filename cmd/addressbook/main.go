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
		func(e Contact) any { return e.PublicId },
		func(t *Contact, a any) { t.PublicId = a.(string) },
		&crud.CrudServiceOptions{},
	)

	tagsRepo := crud.NewCrudService(
		Db,
		func(e Tag) any { return e.PublicId },
		func(t *Tag, a any) { t.PublicId = a.(string) },
		&crud.CrudServiceOptions{},
	)

	crud.AddCrudGinRestApi[Contact]("api/contacts", r, contactsRepo, &crud.CrudRestApiOptions[Contact]{})
	crud.AddCrudGinRestApi[Tag]("api/tags", r, tagsRepo, &crud.CrudRestApiOptions[Tag]{})

	r.Run(":5050")
}
