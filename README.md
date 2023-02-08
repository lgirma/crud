# go-crud

Database CRUD operations utility for Go using [gorm](https://gorm.io).

## Installation

Use go get as follows:

```
go get github.com/lgirma/crud
```

## Usage

First setup your database entities as:

```go
type Contact struct {
	PublicId string `gorm:"primaryKey"`
	FullName string
	Code     int
	Email    string
	Phone    string
}
```

Then create a repository for that entity use the `crud.NewCrudService()` method.

```go
contactRepo = crud.NewCrudService(myDbConnection,
		func(t Contact) any { return t.PublicId },
		func(t *Contact, s any) { t.PublicId = s.(string) },
		&crud.CrudServiceOptions{},
	)
```

There we told the repository how we would read and write the primary key (which is usually a public key).

### Create

To create an entity use `Create` as:

```go
contactRepo.Create(&Contact{
    FullName: "NewCont", 
    Email: "test@mail.com",
})
```

Or use `CreateAll()` to create in batch:

```go
rowsAffected, err := contactRepo.CreateAll([]Contact{
    {FullName: "John", Email: "john@mail.com"},
    {FullName: "Peter", Email: "peter@mail.com"},
})
```

### Read

To fetch list of entities based on a criteria, use `FindWhere()` and `FindByQuery()` as:

```go
// Find contacts with the given name, 1st page, with 10 rows per page:
result, err := contactRepo.FindWhere(&Contact{FullName: "Cont-1"}, Paged(0, 10))
// result.TotalCount - the total number of results regardles of paging
// result.List - the pagenated list
// result.TotalPages - the number of pages

// Find contacts whose full name starts with 'J':
result, err := contactRepo.FindByQuery(Paged(0, 10), "full_name like ?", "J%")
```

To find a single entity based on a criteria, use `FindOneWhere()` or `FindOneByQuery()` as:

```go
// Find a contact with the given name, or return nil if it doesn't exist:
result, err := contactRepo.FindOneWhere(&Contact{Email: "test@mail.com"})

// Find the first contact whose full name starts with 'J':
result, err := contactRepo.FindOneByQuery("full_name like ?", "J%")
```

To count rows, use any of `Count` or `CountWhere` or `CountByQuery` methods as:

```go
// Count all rows
count := contactRepo.Count()

// Count rows with criteria
count, err := contactRepo.CountWhere(&Contact{Email: "test@mail.com"})

// Count rows with query
count, err := contactRepo.CountByQuery("full_name like ?", "J%")
```

### Update

To update an entity, use `Update()` as:

```go
result, err := contactRepo.FindOneWhere(&Contact{Email: "test@mail.com"})

result.Email = "test_update@gmail.com"
contactRepo.Update(result)
```

To do bulk updates, use `UpdateAll()` as:

```go
entities, err := contactRepo.FindByQuery("full_name like ?", "J%")

for i := range entities {
    entities[i].Email += ".et"
}
contactRepo.UpdateAll(entities)
```

### Delete

To delete entities using criteria, use `DeleteWhere()` or `DeleteByQuery()` as follows:

```go
rowsAffected, err := contactRepo.DeleteWhere(&Contact{FullName: "Cont-1"})

contactRepo.DeleteByQuery("full_name like ?", "J%")
```