package crud

import (
	"errors"

	"gorm.io/gorm"
)

type CrudService[T any] interface {
	FindWhere(criteria *T, filter ...*DataFilter) (*PagedList[T], error)
	FindByQuery(filter *DataFilter, query string, paramValues ...any) (*PagedList[T], error)
	GetAll(filter ...*DataFilter) (*PagedList[T], error)

	FindOne(criteria ...*T) (*T, error)
	FindOneWhere(query string, paramValues ...any) (*T, error)

	Count(criteria ...*T) (int, error)
	CountWhere(query string, paramValues ...any) (int, error)

	CreateAll(entities []T) ([]T, error)
	Create(entity *T) (*T, error)

	Delete(criteria *T) (int, error)
	DeleteWhere(query string, paramValues ...any) (int, error)

	Update(entity *T) (int, error)
	UpdateAll(entities []T) (int, error)
	UpdateWhere(entity *T, query string, paramValues ...any) (int, error)
}

type CrudServiceOptions struct {
	DefaultPageSize    int
	PublicIdColumnName string
	IdGenerator        IdGenerator
	DontGenerateIds    bool
}

func GetDefaultCrudServiceOptions() *CrudServiceOptions {
	return &CrudServiceOptions{
		DefaultPageSize:    5,
		PublicIdColumnName: "public_id",
		IdGenerator:        CreateNewIdGenerator(),
		DontGenerateIds:    false,
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
	}
	return &CrudServiceImpl[T]{
		_db:         db,
		SetPublicId: setPublicId,
		GetPublicId: getPublicId,
		_options:    options,
	}
}

func (service *CrudServiceImpl[T]) FindWhere(criteria *T, filters ...*DataFilter) (*PagedList[T], error) {
	var resultList []T
	var filter *DataFilter
	if len(filters) > 0 {
		filter = filters[0]
	}
	normalizeFilter(filter, service._options.DefaultPageSize)
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

func (service *CrudServiceImpl[T]) FindByQuery(filter *DataFilter, query string, paramValues ...any) (*PagedList[T], error) {
	var resultList []T
	normalizeFilter(filter, service._options.DefaultPageSize)
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

func (service *CrudServiceImpl[T]) GetAll(filters ...*DataFilter) (*PagedList[T], error) {
	var filter *DataFilter
	if len(filters) > 0 {
		filter = filters[0]
	}
	return service.FindByQuery(filter, "1=1")
}

func (service *CrudServiceImpl[T]) FindOne(criteria ...*T) (*T, error) {
	var result *PagedList[T]
	var err error
	if len(criteria) > 0 {
		result, err = service.FindWhere(criteria[0], Paged(0, 1))
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
	result, err := service.FindByQuery(Paged(0, 1), query, paramValues...)
	if err != nil {
		return nil, err
	}
	if result.TotalCount == 0 {
		return nil, nil
	}
	return &result.List[0], nil
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

func (service *CrudServiceImpl[T]) DeleteWhere(query string, paramValues ...any) (int, error) {
	db_result := service._db.Where(query, paramValues...).Delete(new(T))
	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(db_result.RowsAffected), nil
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
