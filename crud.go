package crud

import (
	"errors"

	"gorm.io/gorm"
)

type CrudService[T any] interface {
	FindWhere(criteria *T, filter *DataFilter) (*PagedList[T], error)
	FindByQuery(filter *DataFilter, query string, paramValues ...any) (*PagedList[T], error)

	FindOneWhere(criteria *T) (*T, error)
	FindOneByQuery(query string, paramValues ...any) (*T, error)

	Count() int
	CountByQuery(query string, paramValues ...any) (int, error)
	CountWhere(criteria *T) (int, error)

	CreateAll(entities []T) ([]T, error)
	Create(entity *T) (*T, error)

	DeleteWhere(criteria *T) (int, error)
	DeleteByQuery(query string, paramValues ...any) (int, error)

	UpdateAll(entities []T) (int, error)
	Update(entity *T) (int, error)
}

type CrudServiceOptions struct {
	DefaultPageSize    int
	PublicIdColumnName string
	IdGenerator        IdGenerator
}

func GetDefaultCrudServiceOptions() *CrudServiceOptions {
	return &CrudServiceOptions{
		DefaultPageSize:    5,
		PublicIdColumnName: "public_id",
		IdGenerator:        CreateNewIdGenerator(),
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

func (service *CrudServiceImpl[T]) FindWhere(criteria *T, filter *DataFilter) (*PagedList[T], error) {
	var resultList []T
	normalizeFilter(filter, service._options.DefaultPageSize)
	db_result := service._db.Model(new(T)).
		Where(&criteria).
		Limit(filter.ItemsPerPage).
		Offset(filter.CurrentPage * filter.ItemsPerPage).
		Find(&resultList)
	totalCount, err := service.CountWhere(criteria)
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
	totalCount, err := service.CountByQuery(query, paramValues...)
	if db_result.Error != nil {
		return nil, db_result.Error
	}
	if err != nil {
		return nil, err
	}
	result := NewPagedList(resultList, totalCount, filter)
	return result, nil
}

func (service *CrudServiceImpl[T]) FindOneWhere(criteria *T) (*T, error) {
	result, err := service.FindWhere(criteria, Paged(0, 1))
	if err != nil {
		return nil, err
	}
	if result.TotalCount == 0 {
		return nil, nil
	}
	return &result.List[0], nil
}

func (service *CrudServiceImpl[T]) FindOneByQuery(query string, paramValues ...any) (*T, error) {
	result, err := service.FindByQuery(Paged(0, 1), query, paramValues...)
	if err != nil {
		return nil, err
	}
	if result.TotalCount == 0 {
		return nil, nil
	}
	return &result.List[0], nil
}

func (service *CrudServiceImpl[T]) Count() int {
	var result int64
	var model T
	service._db.Model(&model).Count(&result)
	return int(result)
}

func (service *CrudServiceImpl[T]) CountByQuery(query string, paramValues ...any) (int, error) {
	var result int64
	var model T
	db_result := service._db.Model(&model).Where(query, paramValues...).Count(&result)
	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(result), nil
}

func (service *CrudServiceImpl[T]) CountWhere(criteria *T) (int, error) {
	var result int64
	db_result := service._db.Model(criteria).Where(criteria).Count(&result)
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
	result, err := service.CreateAll([]T{*entity})
	if err != nil {
		return nil, err
	}
	if len(result) == 0 {
		return nil, errors.New("entity was not created")
	}
	return &result[0], nil
}

func (service *CrudServiceImpl[T]) DeleteWhere(criteria *T) (int, error) {
	db_result := service._db.Where(criteria).Delete(new(T))
	if db_result.Error != nil {
		return 0, db_result.Error
	}
	return int(db_result.RowsAffected), nil
}

func (service *CrudServiceImpl[T]) DeleteByQuery(query string, paramValues ...any) (int, error) {
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
