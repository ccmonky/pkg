package inithook

import "fmt"

type AttrSetter[T any] func(value T)

func RegisterAttrSetter[T any](attr string, setter AttrSetter[T]) {
	setters[attr] = append(setters[attr], setter)
}

func SetAttr[T any](attr string, value T) error {
	for i, setter := range setters[attr] {
		s, ok := setter.(AttrSetter[T])
		if !ok {
			return fmt.Errorf("attr %s setters[%d] type should be AttrSetter[%T] but got %T", attr, i, *new(T), setter)
		}
		s(value)
	}
	return nil
}

var setters = map[string][]any{}
