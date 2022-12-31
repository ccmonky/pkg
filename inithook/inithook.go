package inithook

import (
	"context"
	"fmt"
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
	settersLock.Lock()
	if setters[attr] == nil {
		setters[attr] = map[setterName]any{}
	}
	setters[attr][name] = setter
	settersLock.Unlock()
	return nil
}

// ExecuteAttrSetters execute attr setters
func ExecuteAttrSetters[T any](ctx context.Context, attr string, value T) error {
	settersLock.Lock()
	attrSetters := setters[attr]
	settersLock.Unlock()
	for name, setter := range attrSetters {
		s, ok := setter.(AttrSetter[T])
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

var (
	setters         = map[attr]map[setterName]any{}
	settersLock     sync.Mutex
	settersUsed     = map[attr]struct{}{}
	settersUsedLock sync.Mutex
)

type attr = string
type setterName = string
