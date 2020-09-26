package resource

import (
    "github.com/graphql-go/graphql"
    "github.com/mitchellh/mapstructure"
    "github.com/stretchr/testify/assert"
    "testing"
)

func setupGraphQL(t *testing.T, schemaConfig graphql.SchemaConfig) graphql.Schema {
    schema, err := graphql.NewSchema(schemaConfig)
    assert.Nil(t, err)
    return schema
}

// DecodeMapToStruct - Decode a map to struct. Userful for GraphQL Params.
func DecodeMapToStruct(input interface{}, output interface{}) error {
    config := &mapstructure.DecoderConfig{
        Metadata: nil,
        Result:   output,
        TagName: "json",
    }
    decoder, err := mapstructure.NewDecoder(config)
    if err != nil {
        return err
    }
    return decoder.Decode(input)
}


func assertListCount(t *testing.T, r *graphql.Result, paginatedName string, expectedCount int) {
   f := r.Data.(map[string]interface{})[paginatedName]
   count := f.(map[string]interface{})["count"].(int)
   assert.Equal(t, expectedCount, count)
}