package crud

import (
	"fmt"
	"strconv"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type TestContact struct {
	Id       int
	FullName string
	PublicId string
	Code     int
	Email    string
	Phone    string
}

var crud_test_db *gorm.DB
var contactsService CrudService[TestContact, string]

func create_and_populate_test_db(seedDataLength int) {
	dbName := GetRandomStr(5)
	Db, err := gorm.Open(sqlite.Open(fmt.Sprintf("file:db_%s?mode=memory&cache=shared", dbName)), &gorm.Config{})
	if err != nil {
		panic("Db connect failed: " + err.Error())
	}
	Db.AutoMigrate(&TestContact{})
	contacts := make([]TestContact, 0)
	for i := 0; i < seedDataLength; i++ {
		istr := strconv.Itoa(i)
		c := TestContact{
			FullName: "Cont-" + istr,
			Email:    "c_" + istr + "@gmail.com",
			PublicId: uuid.NewString(),
		}
		contacts = append(contacts, c)
	}
	Db.Create(&contacts)
	crud_test_db = Db
	contactsService = NewCrudService(crud_test_db,
		func(t TestContact) string { return t.PublicId },
		func(t *TestContact, s string) { t.PublicId = s },
		&CrudServiceOptions[TestContact, string]{
			LookupQuery: "full_name like ? or email like ?",
		},
	)
}

func TestCount(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.Count()
	assert.Nil(t, err)
	assert.Equal(t, 30, result)
}

func TestCountWhere(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.Count(&TestContact{FullName: "Cont-1"})
	assert.Equal(t, 1, result)
	assert.Nil(t, err)
}

func TestCountQuery(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.CountWhere("full_name LIKE ?", "Cont-1%")
	assert.Equal(t, 11, result)
	assert.Nil(t, err)
}

func TestFindAll(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.FindAll(&TestContact{FullName: "Cont-1"}, Paged(0, 10))
	assert.Equal(t, 1, result.TotalCount)
	assert.Equal(t, 1, result.TotalPages)
	assert.Equal(t, 0, result.CurrentPage)
	assert.Equal(t, false, result.HasNext)
	assert.Equal(t, false, result.HasPrevious)
	assert.Len(t, result.List, 1)
	assert.Equal(t, "Cont-1", result.List[0].FullName)
	assert.Equal(t, "c_1@gmail.com", result.List[0].Email)
	assert.Nil(t, err)

	result, err = contactsService.FindAll(&TestContact{Code: 0}, Paged(1, 5))
	assert.Equal(t, 30, result.TotalCount)
	assert.Equal(t, 6, result.TotalPages)
	assert.Equal(t, 1, result.CurrentPage)
	assert.Equal(t, true, result.HasNext)
	assert.Equal(t, true, result.HasPrevious)
	assert.Len(t, result.List, 5)
	assert.Equal(t, "Cont-5", result.List[0].FullName)
	assert.Equal(t, "c_5@gmail.com", result.List[0].Email)
	assert.Nil(t, err)
}

func TestFindAllWhere(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.FindAllWhere("full_name like ?", "Cont-%", Paged(0, 10))
	assert.Equal(t, 30, result.TotalCount)
	assert.Equal(t, 3, result.TotalPages)
	assert.Equal(t, 0, result.CurrentPage)
	assert.Equal(t, true, result.HasNext)
	assert.Equal(t, false, result.HasPrevious)
	assert.Len(t, result.List, 10)
	assert.Equal(t, "Cont-0", result.List[0].FullName)
	assert.Equal(t, "c_0@gmail.com", result.List[0].Email)
	assert.Nil(t, err)

	result, err = contactsService.FindAllWhere("full_name like ?", "Cont-%", Paged(1, 5))
	assert.Equal(t, 30, result.TotalCount)
	assert.Equal(t, 6, result.TotalPages)
	assert.Equal(t, 1, result.CurrentPage)
	assert.Equal(t, true, result.HasNext)
	assert.Equal(t, true, result.HasPrevious)
	assert.Len(t, result.List, 5)
	assert.Equal(t, "Cont-5", result.List[0].FullName)
	assert.Equal(t, "c_5@gmail.com", result.List[0].Email)
	assert.Nil(t, err)

	result, err = contactsService.FindAllWhere("invalid_column like ?", "Cont-%", Paged(1, 5))
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestGetAll(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.GetAll(Paged(0, 10))
	assert.Equal(t, 30, result.TotalCount)
	assert.Equal(t, 3, result.TotalPages)
	assert.Equal(t, 0, result.CurrentPage)
	assert.Equal(t, true, result.HasNext)
	assert.Equal(t, false, result.HasPrevious)
	assert.Len(t, result.List, 10)
	assert.Equal(t, "Cont-0", result.List[0].FullName)
	assert.Equal(t, "c_0@gmail.com", result.List[0].Email)
	assert.Nil(t, err)
}

func TestLookup(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.Lookup("-2")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Len(t, result.List, contactsService.GetOptions().DefaultPageSize)

	result2, err := contactsService.Lookup("-2", Paged(0, 8))
	assert.Nil(t, err)
	assert.NotNil(t, result2)
	assert.Len(t, result2.List, 8)
}

func TestFindOne(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.FindOne(&TestContact{FullName: "Cont-1"})
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Cont-1", result.FullName)
	assert.Equal(t, "c_1@gmail.com", result.Email)

	result, err = contactsService.FindOne(&TestContact{FullName: "258-888"})
	assert.Nil(t, err)
	assert.Nil(t, result)
}

func TestFindOneWhere(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.FindOneWhere("full_name = ?", "Cont-1")
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Cont-1", result.FullName)
	assert.Equal(t, "c_1@gmail.com", result.Email)

	result, err = contactsService.FindOneWhere("full_name = ?", "58-1")
	assert.Nil(t, err)
	assert.Nil(t, result)

	result, err = contactsService.FindOneWhere("non_existing_col = ?", "58-1")
	assert.NotNil(t, err)
	assert.Nil(t, result)
}

func TestFindOneByPublicId(t *testing.T) {
	create_and_populate_test_db(30)

	var entity TestContact
	crud_test_db.Model(&TestContact{}).First(&entity)

	result, err := contactsService.FindOneByPublicId(entity.PublicId)
	assert.Nil(t, err)
	assert.NotNil(t, result)
}

func TestCreateAll(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.CreateAll([]TestContact{
		{FullName: "NewCont", Email: "test@mail.com"},
		{FullName: "NewCont2", Email: "test2@mail.com"},
	})

	var count int64
	crud_test_db.Model(&TestContact{}).Count(&count)

	assert.Nil(t, err)
	assert.Len(t, result, 2)
	assert.Equal(t, 32, int(count))
	assert.NotEmpty(t, result[0].PublicId)
}

func TestCreate(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.Create(&TestContact{FullName: "NewCont", Email: "test@mail.com"})

	var count int64
	crud_test_db.Model(&TestContact{}).Count(&count)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 31, int(count))
	assert.NotEmpty(t, result.PublicId)
}

func TestDelete(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.Delete(&TestContact{FullName: "Cont-1"})

	assert.Nil(t, err)
	assert.Equal(t, 1, result)

	var count int64
	crud_test_db.Model(&TestContact{}).Count(&count)

	assert.Equal(t, 29, int(count))
}

func TestDeleteByPublicId(t *testing.T) {
	create_and_populate_test_db(30)

	var entity TestContact
	crud_test_db.Model(&TestContact{}).First(&entity)

	rowsAffected, err := contactsService.DeleteByPublicId(entity.PublicId)
	assert.Nil(t, err)
	assert.Equal(t, 1, rowsAffected)
}

func TestDeleteAll(t *testing.T) {
	create_and_populate_test_db(30)

	var entities []TestContact
	crud_test_db.Model(&TestContact{}).Limit(2).Find(&entities)

	publicIds := make([]string, 0)
	for _, v := range entities {
		publicIds = append(publicIds, v.PublicId)
	}

	rowsAffected, err := contactsService.DeleteAll(publicIds)
	assert.Nil(t, err)
	assert.Equal(t, 2, rowsAffected)

	crud_test_db.Model(&TestContact{}).Where("public_id in ?", publicIds).Find(&entities)
	assert.Len(t, entities, 0)
}

func TestDeleteWhere(t *testing.T) {
	create_and_populate_test_db(30)
	result, err := contactsService.DeleteWhere("full_name like ?", "Cont-1%")

	assert.Nil(t, err)
	assert.Equal(t, 11, result)

	var count int64
	crud_test_db.Model(&TestContact{}).Count(&count)

	assert.Equal(t, 19, int(count))
}

func TestUpdateAll(t *testing.T) {
	create_and_populate_test_db(30)
	var items []TestContact
	crud_test_db.Model(&TestContact{}).Where("full_name like ?", "Cont-1%").Find(&items)

	for i := range items {
		items[i].FullName += "_updated_" + strconv.Itoa(i)
	}

	count, err := contactsService.UpdateAll(items)

	assert.Nil(t, err)
	assert.Equal(t, 11, count)
}

func TestUpdate(t *testing.T) {
	create_and_populate_test_db(30)
	var entity TestContact
	crud_test_db.Model(&TestContact{}).Where("full_name = ?", "Cont-10").First(&entity)

	entity.FullName += "_updated"

	count, err := contactsService.Update(&entity)

	assert.Nil(t, err)
	assert.Equal(t, 1, count)

	var entityAfterUpdate TestContact
	crud_test_db.Model(&TestContact{}).Where("public_id = ?", entity.PublicId).First(&entityAfterUpdate)

	assert.Equal(t, "Cont-10_updated", entityAfterUpdate.FullName)
}

func TestUpdateWhere(t *testing.T) {
	create_and_populate_test_db(30)
	rowsAffected, err := contactsService.UpdateWhere(&TestContact{Code: 5}, "full_name like ?", "Cont-1%")

	assert.Nil(t, err)
	assert.Equal(t, 11, rowsAffected)

	var items []TestContact
	crud_test_db.Model(&TestContact{}).Where("full_name like ?", "Cont-1%").Find(&items)
	assert.Equal(t, 5, items[0].Code)
}
