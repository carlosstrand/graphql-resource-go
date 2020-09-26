package resource

import (
    pagination "github.com/carlosstrand/graphql-pagination-go"
    "github.com/graphql-go/graphql"
    "github.com/stretchr/testify/assert"
    "testing"
)

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

type BookInput struct {
    Title string `json:"title"`
    Author string `json:"author"`
}

type Book struct {
    Title string `json:"title"`
    Author string `json:"author"`
}

type BooksResource struct {}

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

func (b *BooksResource) List(p graphql.ResolveParams, page pagination.Page) (data interface{}, count int, err error) {
    return BooksMock, 2, nil
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

func createBookResourceSchema(t *testing.T) graphql.Schema {
    rootQuery := graphql.ObjectConfig{Name: "RootQuery", Fields: graphql.Fields{}}
    rootMutation := graphql.ObjectConfig{Name: "RootMutation", Fields: graphql.Fields{}}
    schemaConfig := graphql.SchemaConfig{
        Query: graphql.NewObject(rootQuery),
        Mutation: graphql.NewObject(rootMutation),
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
    return setupGraphQL(t, schemaConfig)
}

func TestBookList(t *testing.T) {
    schema := createBookResourceSchema(t)
    query := `
		{
			bookList {
				data {
                    title
                    author
                }
				count
			}
		}
	`
    params := graphql.Params{Schema: schema, RequestString: query}
    r := graphql.Do(params)
    assert.NotNil(t, r.Data)
    listField := r.Data.(map[string]interface{})["bookList"]
    assert.NotNil(t, listField)
    data := listField.(map[string]interface{})["data"].([]interface{})
    var books []Book
    err := DecodeMapToStruct(data, &books)
    assert.Nil(t, err)
    assert.Equal(t, BooksMock, books)
    assertListCount(t, r, "bookList", 2)
}

func TestBookShow(t *testing.T) {
    schema := createBookResourceSchema(t)
    query := `
		{
			book(id: "123") {
                title
                author
            }
		}
	`
    params := graphql.Params{Schema: schema, RequestString: query}
    r := graphql.Do(params)
    assert.Nil(t, r.Errors)
    assert.NotNil(t, r.Data)
    var book Book
    err := DecodeMapToStruct(r.Data.(map[string]interface{})["book"], &book)
    assert.Nil(t, err)
    assert.Equal(t, "Some book with id=123", book.Title)
}

func TestBookShowWithoutID(t *testing.T) {
    schema := createBookResourceSchema(t)
    query := `
		{
			book {
                title
                author
            }
		}
	`
    params := graphql.Params{Schema: schema, RequestString: query}
    r := graphql.Do(params)
    assert.NotNil(t, r.Errors)
    expectedError := "Field \"book\" argument \"id\" of type \"String!\" is required but not provided."
    assert.Equal(t, expectedError, r.Errors[0].Message)
}

func TestCreateBook(t *testing.T) {
    schema := createBookResourceSchema(t)
    query := `
		mutation {
			createBook(input: {
                title: "New book",
                author: "Carlos Strand"
            }) {
                title
                author
            }
		}
	`
    params := graphql.Params{Schema: schema, RequestString: query}
    r := graphql.Do(params)
    assert.Nil(t, r.Errors)
    assert.NotNil(t, r.Data)
    var book Book
    err := DecodeMapToStruct(r.Data.(map[string]interface{})["createBook"], &book)
    assert.Nil(t, err)
    assert.Equal(t, "New book", book.Title)
    assert.Equal(t, "Carlos Strand", book.Author)
}

func TestUpdateBook(t *testing.T) {
    schema := createBookResourceSchema(t)
    query := `
		mutation {
			updateBook(
                id: "456",
                input: {
                    title: "Updated book",
                }
            ) {
                title
                author
            }
		}
	`
    params := graphql.Params{Schema: schema, RequestString: query}
    r := graphql.Do(params)
    assert.Nil(t, r.Errors)
    assert.NotNil(t, r.Data)
    var book Book
    err := DecodeMapToStruct(r.Data.(map[string]interface{})["updateBook"], &book)
    assert.Nil(t, err)
    assert.Equal(t, "Updated book", book.Title)
}

func TestDeleteBook(t *testing.T) {
    schema := createBookResourceSchema(t)
    query := `
		mutation {
			deleteBook(
                id: "789",
            ) {
                title
                author
            }
		}
	`
    params := graphql.Params{Schema: schema, RequestString: query}
    r := graphql.Do(params)
    assert.Nil(t, r.Errors)
    assert.NotNil(t, r.Data)
    var book Book
    err := DecodeMapToStruct(r.Data.(map[string]interface{})["deleteBook"], &book)
    assert.Nil(t, err)
    assert.Equal(t, "Deleted Book id=789", book.Title)
}