package crud

import (
	"errors"

	"github.com/gin-gonic/gin"
)

type CrudRestApiOptions[T any] struct {
	PrimaryKeyParser func(primaryKeyStringVal string) any
}

func AddCrudGinRestApi[T any](baseUrl string, ginEngine *gin.Engine, crudService CrudService[T], options *CrudRestApiOptions[T]) {
	r := ginEngine
	listEndPoint := func(c *gin.Context) {
		var filter DataFilter
		if c.Request.Method == "POST" {
			if err := c.ShouldBindJSON(&filter); err != nil {
				filter = *Paged(0, crudService.GetOptions().DefaultPageSize)
			}
		} else {
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
	r.POST(baseUrl, listEndPoint)

	r.GET(baseUrl+"/get/:publicId", func(c *gin.Context) {
		publicId := c.Param("publicId")
		result, err := crudService.FindOneByPublicId(publicId)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

	r.GET(baseUrl+"/count", func(c *gin.Context) {
		result, err := crudService.Count()
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

	r.POST(baseUrl+"/create", func(c *gin.Context) {
		var entity T
		if err := c.ShouldBindJSON(&entity); err != nil {
			c.AbortWithError(400, errors.New("invalid data to create"))
		}
		result, err := crudService.Create(&entity)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

	r.POST(baseUrl+"/create-all", func(c *gin.Context) {
		var entities []T
		if err := c.ShouldBindJSON(&entities); err != nil {
			c.AbortWithError(400, errors.New("invalid data to create"))
		}
		result, err := crudService.CreateAll(entities)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

	r.GET(baseUrl+"/delete/:publicId", func(c *gin.Context) {
		publicId := c.Param("publicId")
		result, err := crudService.DeleteByPublicId(publicId)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

	r.POST(baseUrl+"/delete-all", func(c *gin.Context) {
		var publicIdStrings []string
		if err := c.ShouldBindJSON(&publicIdStrings); err != nil {
			c.AbortWithError(400, errors.New("invalid list of public IDs"))
		}
		publicIds := make([]any, 0)
		for _, v := range publicIdStrings {
			publicIds = append(publicIds, v)
		}
		result, err := crudService.DeleteAll(publicIds)
		if err != nil {
			c.AbortWithError(500, err)
		} else {
			c.JSON(200, result)
		}
	})

}
