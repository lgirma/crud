package crud

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
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

func http_req(method string, urlPath string, r *gin.Engine, body any, result ...any) (int, error) {
	w := httptest.NewRecorder()
	rootUrl, _ := url.Parse("/" + test_api_contacts_path)
	subUrl, _ := url.Parse(urlPath)
	rootUrl = rootUrl.JoinPath(subUrl.Path)
	rootUrl.RawQuery = subUrl.RawQuery
	body_json, _ := json.Marshal(body)

	req, err := http.NewRequest(method, rootUrl.String(), bytes.NewReader(body_json))
	r.ServeHTTP(w, req)

	if err != nil {
		return w.Code, err
	}

	bodyContent := w.Body.Bytes()
	if len(bodyContent) == 0 || len(result) == 0 {
		result = nil
		return w.Code, nil
	}
	fmt.Printf("Body: %s", string(bodyContent))
	err = json.Unmarshal(w.Body.Bytes(), result[0])

	return w.Code, err
}

func get_req(urlPath string, r *gin.Engine, result ...any) (int, error) {
	return http_req("GET", urlPath, r, nil, result...)
}

func post_req(urlPath string, r *gin.Engine, body any, result ...any) (int, error) {
	return http_req("POST", urlPath, r, body, result...)
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
	code, err := get_req("?page=1&limit=7", r, &result)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, 1, result.CurrentPage)
	assert.Len(t, result.List, 7)
}

func TestGetApi(t *testing.T) {
	r := setup_test_api()
	var result TestContact
	code, err := get_req(crud_test_public_ids[0], r, &result)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, crud_test_public_ids[0], result.PublicId)
}

func TestGetNotFoundApi(t *testing.T) {
	r := setup_test_api()
	var result TestContact
	code, _ := get_req("non_existing_public_id", r, &result)

	assert.Equal(t, 404, code)
}

func TestCreateBulkApi(t *testing.T) {
	r := setup_test_api()
	entities := []TestContact{
		{FullName: "Mother Nature", Email: "mona@gmail.com"},
		{FullName: "Father Nature", Email: "fana@gmail.com"},
	}
	var rowsAdded []TestContact
	code, err := post_req("", r, entities, &rowsAdded)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.NotNil(t, rowsAdded)
	assert.Len(t, rowsAdded, len(entities))

	var totalCount int64
	crud_test_db.Model(&TestContact{}).Count(&totalCount)
	assert.Equal(t, int64(seed_data_size+len(entities)), totalCount)
}

func TestDeleteApi(t *testing.T) {
	r := setup_test_api()
	rowsDeleted := 0
	code, err := http_req("DELETE", crud_test_public_ids[0], r, nil, &rowsDeleted)

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
	code, err := http_req("DELETE", strings.Join(publicIds, ","), r, nil, &rowsDeleted)

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
	code, err := http_req("DELETE", "non_existing_public_id", r, nil, &rowsDeleted)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Equal(t, 0, rowsDeleted)

	var totalCount int64
	crud_test_db.Model(&TestContact{}).Count(&totalCount)
	assert.Equal(t, int64(seed_data_size), totalCount)
}

func TestUpdateApi(t *testing.T) {
	r := setup_test_api()
	entities := []TestContact{
		{PublicId: crud_test_public_ids[0], FullName: "Mother Nature", Email: "mona1@gmail.com"},
		{PublicId: crud_test_public_ids[1], FullName: "Father Nature", Email: "fana1@gmail.com"},
	}
	rowsAffected := 0
	code, err := http_req("PUT", "", r, entities, &rowsAffected)

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

func TestUpdateNonExistingApi(t *testing.T) {
	r := setup_test_api()
	entities := []TestContact{
		{PublicId: "non_existing_public_id", FullName: "Mother Nature", Email: "mona1@gmail.com"},		
	}
	rowsAffected := 0
	code, err := http_req("PUT", "", r, entities, &rowsAffected)

	assert.Equal(t, 200, code)
	assert.Nil(t, err)
	assert.Zero(t, rowsAffected)	
}
