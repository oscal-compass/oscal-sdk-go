/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package modelutils

import (
	"reflect"
	"strconv"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"

	"github.com/oscal-compass/oscal-sdk-go/internal/set"
)

// NilIfEmpty returns nil if the slice is empty, otherwise returns the original slice.
func NilIfEmpty[T any](slice *[]T) *[]T {
	if slice == nil || len(*slice) == 0 {
		return nil
	}
	return slice
}

// FindValuesByName returns a slice of values in a model associated with a key at any depth of nesting.
func FindValuesByName(model *oscalTypes.OscalModels, name string) []string {
	var results []string
	seen := set.New[uintptr]()
	var walk func(val reflect.Value, key string)
	walk = func(val reflect.Value, key string) {
		if !val.IsValid() {
			return
		}
		for (val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) && !val.IsNil() {
			if val.Kind() == reflect.Ptr {
				ptr := val.Pointer()
				if seen.Has(ptr) {
					return
				}
				seen.Add(ptr)
			}
			val = val.Elem()
		}
		switch val.Kind() {
		case reflect.String:
			if key == name {
				results = append(results, val.String())
			}
		case reflect.Struct:
			t := val.Type()
			for i := 0; i < val.NumField(); i++ {
				walk(val.Field(i), t.Field(i).Name)
			}
		case reflect.Map:
			if val.Type().Key().Kind() == reflect.String {
				for _, key := range val.MapKeys() {
					walk(val.MapIndex(key), key.String())
				}
			}
		case reflect.Slice, reflect.Array:
			for i := 0; i < val.Len(); i++ {
				walk(val.Index(i), strconv.Itoa(i))
			}
		case reflect.Ptr:
			if val.IsNil() {
				return
			}
			walk(val.Elem(), key)
		default:
			// not object-like, do nothing
		}

	}
	walk(reflect.ValueOf(model), "")
	return results
}

// HasDuplicateValuesByName aggregates all nested values of a model associated
// with a key and returns true if a value appears more than once.
func HasDuplicateValuesByName(model *oscalTypes.OscalModels, name string) bool {
	values := FindValuesByName(model, name)
	valueSet := set.New[string]()
	for _, value := range values {
		if valueSet.Has(value) {
			return true
		}
		valueSet.Add(value)
	}
	return false
}
