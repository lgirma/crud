package crud

import (
	"errors"
	"strings"

	"gorm.io/gorm"
)

type CrudService[T any] interface {
	GetAll(filter ...*DataFilter) (*PagedList[T], error)
	FindAll(criteria *T, filter ...*DataFilter) (*PagedList[T], error)
	FindAllWhere(query string, paramValuesAndFilter ...any) (*PagedList[T], error)
	Lookup(searchKey string, filter ...*DataFilter) (*PagedList[T], error)

	FindOne(criteria ...*T) (*T, error)
	FindOneByPublicId(publicId any) (*T, error)
	FindOneWhere(query string, paramValues ...any) (*T, error)

	Count(criteria ...*T) (int, error)
	CountWhere(query string, paramValues ...any) (int, error)

	CreateAll(entities []T) ([]T, error)
	Create(entity *T) (*T, error)

	Delete(criteria *T) (int, error)
	DeleteByPublicId(publicId any) (int, error)
	DeleteAll(publicIds []any) (int, error)
	DeleteWhere(query string, paramValues ...any) (int, error)

	Update(entity *T) (int, error)
	UpdateAll(entities []T) (int, error)
	UpdateWhere(entity *T, query string, paramValues ...any) (int, error)

	GetOptions() CrudServiceOptions
}

type CrudServiceOptions struct {
	DefaultPageSize         int
	PublicIdColumnName      string
	IdGenerator             IdGenerator
	DisableAutoIdGeneration bool
	LookupQuery             string
}

func GetDefaultCrudServiceOptions() *CrudServiceOptions {
	return &CrudServiceOptions{
		DefaultPageSize:         5,
		PublicIdColumnName:      "public_id",
		IdGenerator:             CreateNewIdGenerator(),
		DisableAutoIdGeneration: false,
		LookupQuery:             "",
	}
}

type CrudServiceImpl[T any] struct {
	_db         *gorm.DB
	SetPublicId func(*T, any)
	GetPublicId func(T) any
	_options    *CrudServiceOptions
}

func NewCrudService[T any](db *gorm.DB, getPublicId func(T) any, setPublicId func(*T, any), options *CrudServiceOptions) *CrudServiceImpl[T] {
	defaultOptions := GetDefaultCrudServiceOptions()
	if options == nil {
		options = defaultOptions
	} else {
		if options.DefaultPageSize < 1 {
			options.DefaultPageSize = defaultOptions.DefaultPageSize
		}
		if len(options.PublicIdColumnName) == 0 {
			options.PublicIdColumnName = defaultOptions.PublicIdColumnName
		}
		if options.IdGenerator == nil {
			options.IdGenerator = defaultOptions.IdGenerator
		}
		if len(options.LookupQuery) == 0 {
			options.LookupQuery = defaultOptions.LookupQuery
		}
	}
	return &CrudServiceImpl[T]{
		_db:         db,
		SetPublicId: setPublicId,
		GetPublicId: getPublicId,
		_options:    options,
	}
}

func (service *CrudServiceImpl[T]) FindAll(criteria *T, filterParam ...*DataFilter) (*PagedList[T], error) {
	var resultList []T
	var filter *DataFilter
	if len(filterParam) > 0 {
		filter = filterParam[0]
	}
	filter = normalizeFilter(filter, service._options.DefaultPageSize)
	db_result := service._db.Model(new(T)).
		Where(&criteria).
		Limit(filter.ItemsPerPage).
		Offset(filter.CurrentPage * filter.ItemsPerPage).
		Find(&resultList)
	totalCount, err := service.Count(criteria)
	if db_result.Error != nil {
		return nil, db_result.Error
	}
	if err != nil {
		return nil, err
	}
	result := NewPagedList(resultList, totalCount, filter)
	return result, nil
}

func (service *CrudServiceImpl[T]) FindAllWhere(query string, paramValuesAndFilter ...any) (*PagedList[T], error) {
	var resultList []T
	var filter *DataFilter
	var paramValues []any
	paramValuesCount := strings.Count(query, "?")
	if paramValuesCount > 0 {
		paramValues = paramValuesAndFilter[:len(paramValuesAndFilter)-1]
	}
	if len(paramValuesAndFilter) == len(paramValues)+1 {
		filter = paramValuesAndFilter[len(paramValuesAndFilter)-1].(*DataFilter)
	}
	filter = normalizeFilter(filter, service._options.DefaultPageSize)
	db_result := service._db.Model(new(T)).
		Where(query, paramValues...).
		Limit(filter.ItemsPerPage).
		Offset(filter.CurrentPage * filter.ItemsPerPage).
		Find(&resultList)
	totalCount, err := service.CountWhere(query, paramValues...)
	if db_result.Error != nil {
		return nil, db_result.Error
	}
	if err != nil {
		return nil, err
	}
	result := NewPagedList(resultList, totalCount, filter)
	return result, nil
}

func (service *CrudServiceImpl[T]) Lookup(searchKey string, filter ...*DataFilter) (*PagedList[T], error) {
	if len(service._options.LookupQuery) == 0 {
		return nil, errors.New("lookup query should be provided when using NewCrudService options")
	}
	params := make([]any, 0)
	for i := 0; i < strings.Count(service._options.LookupQuery, "?"); i++ {
		params = append(params, searchKey)
	}
	if len(filter) > 0 {
		params = append(params, filter[0])
	}
	return service.FindAllWhere(service._options.LookupQuery, "%"+searchKey+"%", params)
}

func (service *CrudServiceImpl[T]) GetAll(filters ...*DataFilter) (*PagedList[T], error) {
	var filter *DataFilter
	if len(filters) > 0 {
		filter = filters[0]
	}
	return service.FindAllWhere("1=1", filter)
}

func (service *CrudServiceImpl[T]) FindOne(criteria ...*T) (*T, error) {
	var result *PagedList[T]
	var err error
	if len(criteria) > 0 {
		result, err = service.FindAll(criteria[0], Paged(0, 1))
	} else {
		result, err = service.GetAll(Paged(0, 1))
	}

	if err != nil {
		return nil, err
	}
	if result.TotalCount == 0 {
		return nil, nil
	}
	return &result.List[0], nil
}

func (service *CrudServiceImpl[T]) FindOneWhere(query string, paramValues ...any) (*T, error) {
	params := append(paramValues, Paged(0, 1))
	result, err := service.FindAllWhere(query, params...)
	if err != nil {
		return nil, err
	}
	if result.TotalCount == 0 {
		return nil, nil
	}
	return &result.List[0], nil
}

func (service *CrudServiceImpl[T]) FindOneByPublicId(publicId any) (*T, error) {
	return service.FindOneWhere(service._options.PublicIdColumnName+" = ?", publicId)
}

func (service *CrudServiceImpl[T]) CountWhere(query string, paramValues ...any) (int, error) {
	var result int64
	var model T
	db_result := service._db.Model(&model).Where(query, paramValues...).Count(&result)
	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(result), nil
}

func (service *CrudServiceImpl[T]) Count(criteriaParam ...*T) (int, error) {
	var result int64
	var db_result *gorm.DB
	if len(criteriaParam) > 0 {
		db_result = service._db.Model(criteriaParam[0]).Where(criteriaParam[0]).Count(&result)
	} else {
		db_result = service._db.Model(new(T)).Count(&result)
	}

	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(result), nil
}

func (service *CrudServiceImpl[T]) CreateAll(entities []T) ([]T, error) {
	for i := range entities {
		service.SetPublicId(&entities[i], service._options.IdGenerator.GetNewId())
	}
	db_result := service._db.Create(&entities)
	return entities, db_result.Error
}

func (service *CrudServiceImpl[T]) Create(entity *T) (*T, error) {
	if entity == nil {
		return nil, errors.New("cannot create nil entity")
	}
	result, err := service.CreateAll([]T{*entity})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("entity was not created")
	}
	return &result[0], nil
}

func (service *CrudServiceImpl[T]) Delete(criteria *T) (int, error) {
	db_result := service._db.Where(criteria).Delete(new(T))
	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(db_result.RowsAffected), nil
}

func (service *CrudServiceImpl[T]) DeleteByPublicId(publicId any) (int, error) {
	entity := new(T)
	service.SetPublicId(entity, publicId)
	return service.Delete(entity)
}

func (service *CrudServiceImpl[T]) DeleteWhere(query string, paramValues ...any) (int, error) {
	db_result := service._db.Where(query, paramValues...).Delete(new(T))
	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(db_result.RowsAffected), nil
}

func (service *CrudServiceImpl[T]) DeleteAll(publicIds []any) (int, error) {
	return service.DeleteWhere(service._options.PublicIdColumnName+" in ?", publicIds)
}

func (service *CrudServiceImpl[T]) UpdateAll(entities []T) (int, error) {
	publicIds := make([]any, 0)
	for i := range entities {
		publicIds = append(publicIds, service.GetPublicId(entities[i]))
	}
	db_result := service._db.Where(service._options.PublicIdColumnName+" in ?", publicIds).Save(entities)
	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(db_result.RowsAffected), nil
}

func (service *CrudServiceImpl[T]) Update(entity *T) (int, error) {
	return service.UpdateAll([]T{*entity})
}

func (service *CrudServiceImpl[T]) UpdateWhere(entity *T, query string, paramValues ...any) (int, error) {
	db_result := service._db.Model(new(T)).Where(query, paramValues...).Updates(entity)
	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(db_result.RowsAffected), nil
}

func (service *CrudServiceImpl[T]) GetOptions() CrudServiceOptions {
	return *service._options
}
