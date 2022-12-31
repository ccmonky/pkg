package inithook

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
	"sync"
)

// some builtin attrs
const (
	AppName = "app_name"
	Version = "version"
)

// AttrSetter used to set attr in library
type AttrSetter[T any] func(ctx context.Context, value T) error

// RegisterAttrSetter used to register AttrSetter, and should be used in library init
// NOTE: the execution order of setters cannot be guaranteed!
func RegisterAttrSetter[T any](attr, name string, setter AttrSetter[T]) error {
	if setter == nil {
		return fmt.Errorf("inithook: nil attr setter")
	}
	attrConstructorsLock.Lock()
	if _, ok := attrConstructors[attr]; !ok {
		attrConstructors[attr] = NewConstructor[T]()
	}
	attrConstructorsLock.Unlock()
	settersLock.Lock()
	if setters[attr] == nil {
		setters[attr] = map[setterName]any{}
	}
	setters[attr][name] = genericAttrSetter(func(ctx context.Context, value any) error {
		return setter(ctx, value.(T))
	})
	settersLock.Unlock()
	return nil
}

func ExecuteAllAttrSetters(ctx context.Context, attrsData map[attr]json.RawMessage) error {
	for attr, data := range attrsData {
		fn := GetAttrConstructor(attr)
		if fn == nil {
			return fmt.Errorf("inithook: attr %s constructor not found", attr)
		}
		value := fn()
		err := json.Unmarshal(data, &value)
		if err != nil {
			return fmt.Errorf("inithook: attr %s unmarshal with data %v failed: %v", attr, data, err)
		}
		err = ExecuteAttrSetters(ctx, attr, value)
		if err != nil {
			return err
		}
	}
	return nil
}

// ExecuteAttrSetters execute attr setters
func ExecuteAttrSetters[T any](ctx context.Context, attr string, value T) error {
	settersLock.Lock()
	attrSetters := setters[attr]
	settersLock.Unlock()
	for name, setter := range attrSetters {
		s, ok := setter.(genericAttrSetter)
		if !ok {
			return fmt.Errorf("inithook: attr %s setter %s should be AttrSetter[%T] but got %T", attr, name, *new(T), setter)
		}
		err := s(ctx, value)
		if err != nil {
			return fmt.Errorf("inithook: attr %s setters %s executed failed: %v", attr, name, setter)
		}
	}
	settersUsedLock.Lock()
	settersUsed[attr] = struct{}{}
	settersUsedLock.Unlock()
	return nil
}

// AttrsNotSetted return a slice of attr which has not seted, used to alert
func AttrsNotSetted() []attr {
	var attrNotUsed []attr
	settersLock.Lock()
	settersUsedLock.Lock()
	for attr := range setters {
		if _, ok := settersUsed[attr]; !ok {
			attrNotUsed = append(attrNotUsed, attr)
		}
	}
	settersLock.Unlock()
	settersUsedLock.Unlock()
	return attrNotUsed
}

// NewConstructor returns a constructor which create a new T's instance,
// and New will indirect reflect.Ptr recursively to ensure not return nil pointer
func NewConstructor[T any]() func() any {
	typ := reflect.TypeOf(new(T)).Elem()
	constructorsCacheLock.Lock()
	defer constructorsCacheLock.Unlock()
	if fn, ok := constructorsCache[typ]; ok {
		return fn
	}
	fn := func() any {
		var level int
		typ := reflect.TypeOf(new(T)).Elem()
		for ; typ.Kind() == reflect.Ptr; typ = typ.Elem() {
			level++
		}
		if level == 0 {
			return *new(T)
		}
		value := reflect.Zero(typ)
		for i := 0; i < level; i++ {
			p := reflect.New(value.Type())
			p.Elem().Set(value)
			value = p
		}
		return value.Interface().(T)
	}
	constructorsCache[typ] = fn
	return fn
}

func GetAttrConstructor(attr string) func() any {
	attrConstructorsLock.Lock()
	defer attrConstructorsLock.Unlock()
	return attrConstructors[attr]
}

var (
	setters               = map[attr]map[setterName]any{}
	settersLock           sync.Mutex
	attrConstructors      = map[attr]func() any{}
	attrConstructorsLock  sync.Mutex
	constructorsCache     = map[reflect.Type]func() any{}
	constructorsCacheLock sync.Mutex
	settersUsed           = map[attr]struct{}{}
	settersUsedLock       sync.Mutex
)

type attr = string
type setterName = string

type genericAttrSetter func(ctx context.Context, value any) error
