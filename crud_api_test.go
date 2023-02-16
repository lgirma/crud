package crud

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

var test_api_contacts_path string = "api/contacts"
var seed_data_size = 50

func setup_test_api() *gin.Engine {

	create_and_populate_test_db(seed_data_size)

	r := gin.Default()
	AddCrudGinRestApi(test_api_contacts_path, r, contactsService, &CrudRestApiOptions[TestContact, string]{})
	//r.Run(":55000")

	return r
}

func get_req(url string, r *gin.Engine, result ...any) (int, error) {
	w := httptest.NewRecorder()
	req, err := http.NewRequest("GET", path.Join("/", test_api_contacts_path, url), nil)
	r.ServeHTTP(w, req)

	if err != nil {
		return w.Code, err
	}

	bodyContent := w.Body.Bytes()
	if len(bodyContent) == 0 || len(result) == 0 {
		result = nil
		return w.Code, nil
	}
	err = json.Unmarshal(w.Body.Bytes(), result[0])

	return w.Code, err
}

func post_req(url string, r *gin.Engine, body any, result ...any) (int, error) {
	w := httptest.NewRecorder()
	body_json, _ := json.Marshal(body)
	req, err := http.NewRequest("POST", path.Join("/", test_api_contacts_path, url), bytes.NewReader(body_json))
	r.ServeHTTP(w, req)

	if err != nil {
		return w.Code, err
	}

	bodyContent := w.Body.Bytes()
	if len(bodyContent) == 0 || len(result) == 0 {
		return w.Code, nil
	}
	err = json.Unmarshal(w.Body.Bytes(), result[0])

	return w.Code, err
}

func TestDefaultListApi(t *testing.T) {
	r := setup_test_api()
	var result PagedList[TestContact]
	code, err := get_req("", r, &result)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Zero(t, result.CurrentPage)
	assert.Len(t, result.List, contactsService.GetOptions().DefaultPageSize)
}

func TestPagedListApi(t *testing.T) {
	r := setup_test_api()
	var result PagedList[TestContact]
	code, err := post_req("", r, Paged(1, 7), &result)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.CurrentPage)
	assert.Len(t, result.List, 7)
}

func TestGetApi(t *testing.T) {
	r := setup_test_api()
	var result TestContact
	code, err := get_req("get/"+crud_test_public_ids[0], r, &result)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, crud_test_public_ids[0], result.PublicId)
}

func TestGetNotFoundApi(t *testing.T) {
	r := setup_test_api()
	var result TestContact
	code, _ := get_req("get/non_existing_public_id", r, &result)

	assert.Equal(t, 404, code)
}

func TestCountApi(t *testing.T) {
	r := setup_test_api()
	var result int
	code, err := get_req("count", r, &result)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Equal(t, seed_data_size, result)
}

func TestCreateApi(t *testing.T) {
	r := setup_test_api()
	entity := &TestContact{
		FullName: "Mother Nature",
		Email:    "mona@gmail.com",
	}
	var rowAdded TestContact
	code, err := post_req("create", r, entity, &rowAdded)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.NotNil(t, rowAdded)
	assert.NotEmpty(t, rowAdded.PublicId)

	var totalCount int64
	crud_test_db.Model(&TestContact{}).Count(&totalCount)
	assert.Equal(t, int64(seed_data_size+1), totalCount)
}

func TestDeleteApi(t *testing.T) {
	r := setup_test_api()
	rowsDeleted := 0
	code, err := get_req("delete/"+crud_test_public_ids[0], r, &rowsDeleted)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Equal(t, 1, rowsDeleted)

	var totalCount int64
	crud_test_db.Model(&TestContact{}).Count(&totalCount)
	assert.Equal(t, int64(seed_data_size-1), totalCount)
}

func TestDeleteAllApi(t *testing.T) {
	r := setup_test_api()
	rowsDeleted := 0
	publicIds := []string{crud_test_public_ids[0], crud_test_public_ids[1]}
	code, err := post_req("delete-all", r, publicIds, &rowsDeleted)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Equal(t, len(publicIds), rowsDeleted)

	var totalCount int64
	crud_test_db.Model(&TestContact{}).Count(&totalCount)
	assert.Equal(t, int64(seed_data_size-len(publicIds)), totalCount)

	var count int64
	crud_test_db.Model(&TestContact{}).
		Where(contactsService.GetOptions().PublicIdColumnName+" in ?", publicIds).
		Count(&count)
	assert.Zero(t, count)
}

func TestDeleteNotFoundApi(t *testing.T) {
	r := setup_test_api()
	rowsDeleted := 0
	code, err := get_req("delete/non_existing_public_id", r, &rowsDeleted)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Equal(t, 0, rowsDeleted)

	var totalCount int64
	crud_test_db.Model(&TestContact{}).Count(&totalCount)
	assert.Equal(t, int64(seed_data_size), totalCount)
}

func TestUpdateApi(t *testing.T) {
	r := setup_test_api()
	entity := &TestContact{
		PublicId: crud_test_public_ids[0],
		FullName: "Mother Nature",
		Email:    "mona@gmail.com",
	}
	rowsAffected := 0
	code, err := post_req("update", r, entity, &rowsAffected)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Equal(t, 1, rowsAffected)

	var updatedEntity TestContact
	crud_test_db.Model(&TestContact{}).Where("public_id = ?", crud_test_public_ids[0]).First(&updatedEntity)

	assert.NotNil(t, updatedEntity)
	assert.Equal(t, entity.PublicId, updatedEntity.PublicId)
	assert.Equal(t, entity.FullName, updatedEntity.FullName)
	assert.Equal(t, entity.Email, updatedEntity.Email)
}

func TestUpdateNonExistingApi(t *testing.T) {
	r := setup_test_api()
	entity := &TestContact{
		PublicId: "invalid_public_id",
		FullName: "Mother Nature",
		Email:    "mona@gmail.com",
	}
	rowsAffected := 0
	code, err := post_req("update", r, entity, &rowsAffected)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Equal(t, 0, rowsAffected)
}

func TestUpdateAllApi(t *testing.T) {
	r := setup_test_api()
	entities := []TestContact{
		{PublicId: crud_test_public_ids[0], FullName: "Mother Nature", Email: "mona1@gmail.com"},
		{PublicId: crud_test_public_ids[1], FullName: "Father Nature", Email: "fana1@gmail.com"},
	}
	rowsAffected := 0
	code, err := post_req("update-all", r, entities, &rowsAffected)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Equal(t, 2, rowsAffected)

	var updatedEntities []TestContact
	crud_test_db.
		Model(&TestContact{}).
		Where("public_id in ?", []string{crud_test_public_ids[0], crud_test_public_ids[1]}).
		Find(&updatedEntities)

	assert.NotNil(t, updatedEntities)
	assert.Len(t, updatedEntities, 2)
	for _, entity := range entities {
		for _, updatedEntity := range updatedEntities {
			if entity.PublicId == updatedEntity.PublicId {
				assert.Equal(t, updatedEntity.FullName, entity.FullName)
				assert.Equal(t, updatedEntity.Email, entity.Email)
				assert.Equal(t, updatedEntity.FullName, entity.FullName)
				assert.Equal(t, updatedEntity.Email, entity.Email)
			}
		}
	}
}
