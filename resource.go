package resource

import (
    "fmt"
    pagination "github.com/carlosstrand/graphql-pagination-go"
    "github.com/graphql-go/graphql"
    "strings"
)

type Resource interface {
    List(p graphql.ResolveParams, page pagination.Page) (data interface{}, count int, err error)
    Show(p graphql.ResolveParams, id string) (interface{}, error)
    Create(p graphql.ResolveParams, input interface{}) (interface{}, error)
    Update(p graphql.ResolveParams, id string, input interface{}) (interface{}, error)
    Destroy(p graphql.ResolveParams, id string) (interface{}, error)
}

type Options struct {
    Resource          Resource
    Type              graphql.Type
    CreateInputFields graphql.InputObjectConfigFieldMap
    UpdateInputFields graphql.InputObjectConfigFieldMap
}

func AddResourceToSchemaConfig(schemaConfig graphql.SchemaConfig, opts Options) {
    name := strings.ToLower(opts.Type.Name())
    nameTitle := strings.Title(name)
    f := resourceFactory{opts}
    schemaConfig.Query.AddFieldConfig(name + "List", f.MakeList())
    schemaConfig.Query.AddFieldConfig(name, f.MakeShow())
    schemaConfig.Mutation.AddFieldConfig("create"+nameTitle, f.MakeCreate())
    schemaConfig.Mutation.AddFieldConfig("update"+nameTitle, f.MakeUpdate())
    schemaConfig.Mutation.AddFieldConfig("delete"+nameTitle, f.MakeDelete())
}

type resourceFactory struct {
    opts Options
}

func (f *resourceFactory) MakeList() *graphql.Field {
    return pagination.Paginated(&pagination.PaginatedField{
        Name: f.opts.Type.Name() + "List",
        Type: f.opts.Type,
        DataAndCountResolve: func(p graphql.ResolveParams, page pagination.Page) (interface{}, int, error) {
            return f.opts.Resource.List(p, page)
        },
    })
}

func (f *resourceFactory) MakeShow() *graphql.Field {
    return &graphql.Field{
        Type: f.opts.Type,
        Args: graphql.FieldConfigArgument{
            "id": &graphql.ArgumentConfig{
                Type: graphql.NewNonNull(graphql.String),
            },
        },
        Resolve: func(p graphql.ResolveParams) (i interface{}, e error) {
            id := p.Args["id"].(string)
            return f.opts.Resource.Show(p, id)
        },
    }
}


func (f *resourceFactory) MakeCreate() *graphql.Field {
    return &graphql.Field{
        Type: f.opts.Type,
        Args: graphql.FieldConfigArgument{
            "input": &graphql.ArgumentConfig{
                Type: graphql.NewInputObject(graphql.InputObjectConfig{
                    Name: "Create"+strings.Title(f.opts.Type.Name())+"Input",
                    Fields: f.opts.CreateInputFields,
                }),
            },
        },
        Resolve: func(p graphql.ResolveParams) (interface{}, error) {
            return f.opts.Resource.Create(p, p.Args["input"])
        },
        Description: fmt.Sprintf("Create %s", f.opts.Type),
    }
}

func (f *resourceFactory) MakeUpdate() *graphql.Field {
    return &graphql.Field{
        Type: f.opts.Type,
        Args: graphql.FieldConfigArgument{
            "id": &graphql.ArgumentConfig{
                Type: graphql.NewNonNull(graphql.String),
            },
            "input": &graphql.ArgumentConfig{
                Type: graphql.NewInputObject(graphql.InputObjectConfig{
                    Name: "Update"+strings.Title(f.opts.Type.Name())+"Input",
                    Fields: f.opts.UpdateInputFields,
                }),
            },
        },
        Resolve: func(p graphql.ResolveParams) (interface{}, error) {
            id := p.Args["id"].(string)
            return f.opts.Resource.Update(p, id, p.Args["input"])
        },
        Description: fmt.Sprintf("Update %s by id", f.opts.Type),
    }
}

func (f *resourceFactory) MakeDelete() *graphql.Field {
    return &graphql.Field{
        Type: f.opts.Type,
        Args: graphql.FieldConfigArgument{
            "id": &graphql.ArgumentConfig{
                Type: graphql.NewNonNull(graphql.String),
            },
        },
        Resolve: func(p graphql.ResolveParams) (interface{}, error) {
            id := p.Args["id"].(string)
            return f.opts.Resource.Destroy(p, id)
        },
        Description: fmt.Sprintf("Delete %s by id", f.opts.Type),
    }
}