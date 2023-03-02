package crud

import (
	"errors"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

type CrudRestApiOptions[T any, TPublicId any] struct {
}

func AddCrudGinRestApi[T any, TPublicId any](baseUrl string, ginEngine *gin.Engine, crudService CrudService[T, TPublicId], options *CrudRestApiOptions[T, TPublicId]) {
	r := ginEngine
	listEndPoint := func(c *gin.Context) {
		var filter DataFilter
		if err := c.ShouldBind(&filter); err != nil {
			filter = *Paged(0, crudService.GetOptions().DefaultPageSize)
		}
		result, err := crudService.GetAll(&filter)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	}

	r.GET(baseUrl, listEndPoint)

	getOneEndPoint := func(c *gin.Context) {
		publicId := c.Param("publicId")
		result, err := crudService.FindOneByPublicId(Parse[TPublicId](publicId))
		if err != nil {
			c.AbortWithError(500, err)
		} else if result == nil {
			c.AbortWithError(404, errors.New("not found"))
		} else {
			c.JSON(200, result)
		}
	}

	r.GET(baseUrl+"/:publicId", getOneEndPoint)

	r.POST(baseUrl, func(c *gin.Context) {
		var entities []T
		if err := c.ShouldBindJSON(&entities); err != nil {
			c.AbortWithError(400, errors.New("invalid data to create"))
			log.Printf("failed to bind create data: %v", err)
			return
		}
		result, err := crudService.CreateAll(entities)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

	r.PUT(baseUrl, func(c *gin.Context) {
		var entities []T
		if err := c.ShouldBindJSON(&entities); err != nil {
			c.AbortWithError(400, errors.New("invalid data to create"))
		}
		result, err := crudService.UpdateAll(entities)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

	r.DELETE(baseUrl+"/:publicIds", func(c *gin.Context) {
		publicIdSrc := c.Param("publicIds")
		publicIdStrings := strings.Split(publicIdSrc, ",")
		publicIds := make([]TPublicId, 0)
		for _, v := range publicIdStrings {
			publicIds = append(publicIds, Parse[TPublicId](v))
		}
		result, err := crudService.DeleteAll(publicIds)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

}
