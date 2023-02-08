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