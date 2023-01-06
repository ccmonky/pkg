package jsonschema

import (
	"errors"
	"fmt"
	"reflect"
	"sync"
)

// Validator used to validate json data against schema generated according to a type(reflect by a value)
type Validator struct {
	retag     Retag
	tagMaker  TagMaker
	generator Generator
	validate  ValidateFunc
	schemas   sync.Map // map[reflect.Type][]byte
}

// NewValidator returns a new Validator
func NewValidator(opts ...ValidatorOption) (*Validator, error) {
	v := Validator{}
	for _, opt := range opts {
		opt(&v)
	}
	if v.generator == nil {
		return nil, fmt.Errorf("json schema generator not specified")
	}
	if v.validate == nil {
		return nil, fmt.Errorf("json schema validate func not specified")
	}
	return &v, nil
}

func (v *Validator) Schemas() map[string]string {
	schemas := make(map[string]string)
	fn := func(k, v any) bool {
		rtype := k.(reflect.Type)
		typeName := rtype.String()
		for ; rtype.Kind() == reflect.Ptr; rtype = rtype.Elem() {
		}
		schemas[rtype.PkgPath()+":"+typeName] = string(v.([]byte))
		return true
	}
	v.schemas.Range(fn)
	return schemas
}

// ValidatorOption control option to create Validator
type ValidatorOption func(*Validator)

// WithRetag used to set Retag which will be used to convert a value to a new type modified from the original type and tag maker
func WithRetag(retag Retag) ValidatorOption {
	return func(v *Validator) {
		v.retag = retag
	}
}

// WithTagMaker used to set TagMaker which will be used to modify T's tags dynamicly
func WithTagMaker(tm TagMaker) ValidatorOption {
	return func(v *Validator) {
		v.tagMaker = tm
	}
}

// WithGenerator used to set json schema generator
func WithGenerator(jsg Generator) ValidatorOption {
	return func(v *Validator) {
		v.generator = jsg
	}
}

// WithValidateFunc used to set json schema validate func
func WithValidateFunc(fn ValidateFunc) ValidatorOption {
	return func(v *Validator) {
		v.validate = fn
	}
}

// Validate validate json data against schema(generated according to value's Type and optional TagMaker)
func (v *Validator) Validate(value any, data []byte) error {
	rtype := reflect.TypeOf(value)
	schema, ok := v.schemas.Load(rtype)
	if ok {
		return v.validate(schema.([]byte), data)
	}
	rtypeModified := rtype
	tagMaker := v.tagMaker
	if tm, ok := value.(TagMaker); ok {
		tagMaker = tm
	}
	if v.retag != nil && tagMaker != nil {
		modified := v.retag.ConvertAny(value, tagMaker)
		rtypeModified = reflect.TypeOf(modified)
	}
	schemaBytes, err := v.generator.ReflectFromType(rtypeModified)
	if err != nil {
		return err
	}
	v.schemas.Store(rtype, schemaBytes)
	return v.validate(schemaBytes, data)
}

// refer to `github.com/maseer/retag`
type Retag interface {
	Convert(p interface{}, maker TagMaker) interface{}
	ConvertAny(p interface{}, maker TagMaker) interface{}
}

// refer to `github.com/maseer/retag`
type TagMaker interface {
	MakeTag(structureType reflect.Type, fieldIndex int) reflect.StructTag
}

// Generator used to generate json schema for value
// refer to `github.com/invopop/jsonschema`
type Generator interface {
	Reflect(v interface{}) (schema, error)
	ReflectFromType(t reflect.Type) (schema, error)
}

type schema = []byte

// ValidateFunc validate json `data` against json `schema`
// refer to `github.com/xeipuuv/gojsonschema`
type ValidateFunc func(schema, data []byte) error

type ValidateFailedError struct {
	label string
}

func NewValidateFailedErrorError(label string) *ValidateFailedError {
	return &ValidateFailedError{
		label: label,
	}
}

// Error implements the error interface.
func (e *ValidateFailedError) Error() string {
	return "jsonschema: " + e.label
}

func IsValidateFailedError(err error) bool {
	if err == nil {
		return false
	}
	var e *ValidateFailedError
	return errors.As(err, &e)
}
