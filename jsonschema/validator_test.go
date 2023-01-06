package jsonschema_test

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"

	"github.com/ccmonky/pkg/jsonschema"
	jsgen "github.com/invopop/jsonschema"
	"github.com/sevlyar/retag"
	"github.com/xeipuuv/gojsonschema"
)

type TestType struct {
	Int    int `json:"int,omitempty"`
	String string
	Uint   uint `json:"-"`
	Embed  struct {
		A string `json:"a,omitempty"`
		B int
		C uint `json:"-"`
	}
}

func TestValidator(t *testing.T) {
	tt := TestType{}
	data := []byte(`{
		"Uint": 1,
		"Embed": {
			"C": 2
		}
	}`)
	err := validator.Validate(&tt, data)
	if err == nil {
		// - (root): int is required
		// - (root): String is required
		// - (root): Additional property Uint is not allowed
		// - Embed: a is required
		// - Embed: B is required
		// - Embed: Additional property C is not allowed
		t.Fatal("should error")
	}
	if !jsonschema.IsValidateFailedError(err) {
		t.Fatal(err)
	}
	schemas := validator.Schemas()
	// {
	// 	"$schema": "https://json-schema.org/draft/2020-12/schema",
	// 	"properties": {
	// 		"int": {
	// 			"type": "integer"
	// 		},
	// 		"String": {
	// 			"type": "string"
	// 		},
	// 		"Embed": {
	// 			"properties": {
	// 				"a": {
	// 					"type": "string"
	// 				},
	// 				"B": {
	// 					"type": "integer"
	// 				}
	// 			},
	// 			"additionalProperties": false,
	// 			"type": "object",
	// 			"required": [
	// 				"a",
	// 				"B"
	// 			]
	// 		}
	// 	},
	// 	"additionalProperties": false,
	// 	"type": "object",
	// 	"required": [
	// 		"int",
	// 		"String",
	// 		"Embed"
	// 	]
	// }
	t.Log(schemas)
	data = []byte(`{
		"int": 0,
		"String": "",
		"Embed": {
			"a": "",
			"B": 1
		}
	}`)
	err = validator.Validate(&tt, data)
	if err != nil {
		t.Fatal(err)
	}
}

type DefaultRetager struct{}

func (rt DefaultRetager) Convert(p interface{}, maker jsonschema.TagMaker) interface{} {
	return retag.Convert(indirect(p), maker)
}

func (rt DefaultRetager) ConvertAny(p interface{}, maker jsonschema.TagMaker) interface{} {
	return retag.ConvertAny(indirect(p), maker)
}

func indirect(v interface{}) interface{} {
	rtype := reflect.TypeOf(v)
	if rtype.Kind() != reflect.Ptr {
		return reflect.New(rtype).Interface()
	}
	return v
}

func DeleteJsonOmitemptyMarker() retag.TagMaker {
	return deleteJsonOmitemptyMarker{}
}

type deleteJsonOmitemptyMarker struct{}

func (m deleteJsonOmitemptyMarker) MakeTag(t reflect.Type, fieldIndex int) reflect.StructTag {
	jsonTag, ok := t.Field(fieldIndex).Tag.Lookup("json")
	if !ok {
		return ""
	}
	if strings.Contains(jsonTag, ",omitempty") {
		key := strings.Split(jsonTag, ",")[0]
		return reflect.StructTag(fmt.Sprintf(`json:"%s"`, key))
	}
	return reflect.StructTag(fmt.Sprintf(`json:"%s"`, jsonTag))
}

type DefaultGenerator struct{}

func (g DefaultGenerator) Reflect(v interface{}) ([]byte, error) {
	schema := jsgen.Reflect(v)
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func (g DefaultGenerator) ReflectFromType(t reflect.Type) ([]byte, error) {
	schema := jsgen.ReflectFromType(t)
	data, err := json.Marshal(schema)
	if err != nil {
		return nil, err
	}
	return data, nil
}

func DefaultValidate(schema, data []byte) error {
	result, err := gojsonschema.Validate(gojsonschema.NewBytesLoader(schema), gojsonschema.NewBytesLoader(data))
	if err != nil {
		return err
	}
	if result.Valid() {
		return nil
	}
	detail := ""
	for _, desc := range result.Errors() {
		detail += fmt.Sprintf("- %s\n", desc)
	}
	return jsonschema.NewValidateFailedErrorError(detail)
}

func BenchmarkRegisterTypeMultiple(b *testing.B) {
	data := []byte(`{
		"int": 0,
		"String": "",
		"Embed": {
			"a": "",
			"B": 1
		}
	}`)
	tt := TestType{}
	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		err := validator.Validate(&tt, data)
		if err != nil {
			b.Fatal(err)
		}
	}
}

var validator *jsonschema.Validator

func init() {
	var err error
	validator, err = jsonschema.NewValidator(
		jsonschema.WithRetag(DefaultRetager{}),
		jsonschema.WithTagMaker(DeleteJsonOmitemptyMarker()),
		jsonschema.WithGenerator(DefaultGenerator{}),
		jsonschema.WithValidateFunc(DefaultValidate),
	)
	if err != nil {
		panic(err)
	}
}
