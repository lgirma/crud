# crud

Database CRUD operations utility for Go using [gorm](https://gorm.io).

## Table of contents

- [crud](#crud)
  - [Table of contents](#table-of-contents)
  - [Installation](#installation)
  - [Usage](#usage)
    - [Create](#create)
    - [Read](#read)
    - [Update](#update)
    - [Delete](#delete)
    - [Options](#options)
  - [REST API](#rest-api)
  - [Features](#features)

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
// NewCrudService(connection, PK_getter, PK_setter, options)
contactRepo = crud.NewCrudService(
  myDbConnection,
  func(e Contact) string { return e.PublicId },
  func(t *Contact, a string) { t.PublicId = a },
  &crud.CrudServiceOptions[Contact, string]{}
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

To fetch list of entities based on a criteria, use `GetAll()`, `FindAll()` and `FindAllWhere()` as:

```go
// Find all contacts with the default paging (first 10 rows)
result, err := contactRepo.GetAll()

// Find the first 5 contacts
result2, err := contactRepo.GetAll(Paged(0, 5))

// Find contacts with the given name, 1st page, with 10 rows per page:
result, err := contactRepo.FindAll(&Contact{FullName: "Cont-1"}, Paged(0, 10))
// result.TotalCount - the total number of results regardles of paging
// result.List - the pagenated list
// result.TotalPages - the number of pages

// Find contacts whose full name starts with 'J':
result, err := contactRepo.FindAllWhere("full_name like ?", "J%")

// Same with explicit paging:
result, err := contactRepo.FindAllWhere("full_name like ?", "J%", Paged(0, 5))
```

To find a single entity based on a criteria, use `FindOne()` or `FindOneWhere()` as:

```go
// Find the first contact or return nil if there isn't any
result, err := contactRepo.FindOne()

// Find a contact with the given name, or return nil if it doesn't exist:
result, err := contactRepo.FindOne(&Contact{Email: "test@mail.com"})

// Find the first contact whose full name starts with 'J':
result, err := contactRepo.FindOneWhere("full_name like ?", "J%")
```

To count rows, use any of `Count` or `CountWhere` methods as:

```go
// Count all rows
count := contactRepo.Count()

// Count rows with criteria
count, err := contactRepo.Count(&Contact{Email: "test@mail.com"})

// Count rows with query
count, err := contactRepo.CountWhere("full_name like ?", "J%")
```

### Update

To update an entity, use `Update()` as:

```go
result, err := contactRepo.FindOne(&Contact{Email: "test@mail.com"})

result.Email = "test_update@gmail.com"
contactRepo.Update(result)
```

To do bulk updates, use `UpdateAll()` as:

```go
entities, err := contactRepo.FindAllWhere("full_name like ?", "J%")

for i := range entities {
    entities[i].Email += ".et"
}
contactRepo.UpdateAll(entities)
```

To do bulk updates using queries, use `UpdateWhere()` as:

```go
// Equivalent to: UPDATE contacts SET Code = 5 WHERE full_name like 'J%'
rowsAffected, err := contactRepo.UpdateWhere(&Contact{Code: 5}, "full_name like ?", "J%")
```

### Delete

To delete entities using criteria, use `Delete()` or `DeleteWhere()` as follows:

```go
// Equivalent to: DELTE FROM contacts WHERE full_name = 'Cont-1'
rowsAffected, err := contactRepo.DeleteWher(&Contact{FullName: "Cont-1"})

// Equivalent to: DELTE FROM contacts WHERE full_name LIKE 'J%'
contactRepo.DeleteWhere("full_name like ?", "J%")
```

### Options

You can supply options struct when creating the CRUD service as:

```go
contactRepo = crud.NewCrudService(
  myDbConnection,
  func(e Contact) string { return e.PublicId },
  func(t *Contact, a string) { t.PublicId = a },
  &crud.CrudServiceOptions[Contact, string]{
    DefaultPageSize: 25,
    PublicIdColumnName: "uuid",    
    DisableAutoIdGeneration: true,
    LookupQuery: "full_name like ? or email like ?",
  }
)
```

## REST API

You can start a REST API for your CRUD service based on gin gonic, as:

```go
tagsRepo := crud.NewCrudService(/* constructor */)
r := gin.Default()

// Create a rest api at api/tags end-point for the entity Tag with a public Id of type string
crud.AddCrudGinRestApi[Tag, string]("api/tags", r, tagsRepo, nil)
```

Then you will have these end-points automatically:

| End-point | Method | Description | Example URL | Example Body |
|-----------|--------|-------------|-------------|--------------|
| `api/tags` | GET | Paginated list | `api/tags` | - |
| `api/tags` | GET | Paginated list | `api/tags?limit=2&page=3&sort=name:asc,age:desc` | - |
| `api/tags/:publicId` | GET | Details of a single item | `api/tags/e61bc045` | - |
| `api/tags/:publicIds` | DELETE | Deletes items with the given public IDs | `api/tags/e61bc045,cb837345` | - |
| `api/tags` | POST | Creates the given list of items in the body | `api/tags` | `[{"Name": "finance"}, {"Name": "technology"}]` |
| `api/tags` | PUT | Updates the given list of items in the body | `api/tags` | `[{"Id": "cb837345", "Name": "books"}]` |

## Features

- [x] CRUD Service
  - [x] Create
  - [x] Read
  - [x] Update
  - [x] Delete
- [x] Sort
- [x] CRUD Web Api
- [ ] Custom Filters
- [ ] Validation
- [ ] Error handling
  - [ ] Separate 404s and 400s instead of 500
- [ ] Metadata
- [ ] Consistent casing: snake, camel