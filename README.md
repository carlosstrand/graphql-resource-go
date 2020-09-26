[![Documentation](https://godoc.org/github.com/carlosstrand/graphql-resource-go?status.svg)](http://godoc.org/github.com/carlosstrand/graphql-resource-go)
[![Actions Status](https://github.com/carlosstrand/graphql-resource-go/workflows/Go/badge.svg)](https://github.com/carlosstrand/graphql-resource-go/actions)
[![Coverage Status](https://coveralls.io/repos/github/carlosstrand/graphql-resource-go/badge.svg?branch=master)](https://coveralls.io/github/carlosstrand/graphql-resource-go?branch=master)

# graphql-resource-go

Easily create CRUDs using graphql-go

### Usage:

```go
    var BookObject = graphql.NewObject(graphql.ObjectConfig{
        Name: "Book",
        Fields: graphql.Fields{
            "title": &graphql.Field{
                Type: graphql.String,
            },
            "author": &graphql.Field{
                Type: graphql.String,
            },
        },
    })

    type BooksResource struct {}

    func (b *BooksResource) List(p graphql.ResolveParams, page pagination.Page) (data interface{}, count int, err error) {
        var BooksMock = []Book{
            {
                Title: "Don Quixote",
                Author: "Miguel de Cervantes",
            },
            {
                Title: "The Stranger",
                Author: "Albert Camus",
            },
        }
        totalCount := 2
        return BooksMock, totalCount, nil
    }
    
    func (b *BooksResource) Show(p graphql.ResolveParams, id string) (interface{}, error) {
        book := Book{
            Title:  "Some book with id=" +id,
        }
        return book, nil
    }
    
    func (b *BooksResource) Create(p graphql.ResolveParams, input interface{}) (interface{}, error) {
        var i BookInput
        err := DecodeMapToStruct(input, &i)
        if err != nil {
            return nil, err
        }
        book := &Book{
            Title:  i.Title,
            Author:  i.Author,
        }
        return book, nil
    }
    
    func (b *BooksResource) Update(p graphql.ResolveParams, id string, input interface{}) (interface{}, error) {
        var i BookInput
        err := DecodeMapToStruct(input, &i)
        if err != nil {
            return nil, err
        }
        book := &Book{
            Title:  i.Title,
            Author:  i.Author,
        }
        return book, nil
    }
    
    func (b *BooksResource) Destroy(p graphql.ResolveParams, id string) (interface{}, error) {
        book := &Book{
            Title:  "Deleted Book id="+id,
        }
        return book, nil
    }

    AddResourceToSchemaConfig(schemaConfig, Options{
        Resource:          &BooksResource{},
        Type:              BookObject,
        CreateInputFields: graphql.InputObjectConfigFieldMap{
            "title": &graphql.InputObjectFieldConfig{
                Type:        graphql.String,
                Description: "Book title",
            },
            "author": &graphql.InputObjectFieldConfig{
                Type:        graphql.String,
                Description: "Book name",
            },
        },
        UpdateInputFields: graphql.InputObjectConfigFieldMap{
            "title": &graphql.InputObjectFieldConfig{
                Type:        graphql.String,
                Description: "Book title",
            },
        },
    })
```

Done! All resources should be created and you can query like:

```gql
    query {
        bookList {
            data {
                title
                author
            }
        }
        book(id: "123") {
            title
            author
        }
    }
    mutation {
        createBook(
            input: {
                title: "New Book",
                author: "Foo Bar"
            }
        ) {
            title
            author
        }
        updateBook(
            id: "432",
            input: {
                title: "New Book",
                author: "Foo Bar"
            }
        ) {
            title
            author
        }
        deleteBook(id: "432") {
            title
            author
        }
    }
```


# Resource Interface

You need implement a resource interface:

```go
type Resource interface {
    List(p graphql.ResolveParams, page pagination.Page) (data interface{}, count int, err error)
    Show(p graphql.ResolveParams, id string) (interface{}, error)
    Create(p graphql.ResolveParams, input interface{}) (interface{}, error)
    Update(p graphql.ResolveParams, id string, input interface{}) (interface{}, error)
    Destroy(p graphql.ResolveParams, id string) (interface{}, error)
}
```