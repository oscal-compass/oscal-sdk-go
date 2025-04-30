/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package modelutils

import (
	"reflect"
	"strconv"

	oscalTypes "github.com/defenseunicorns/go-oscal/src/types/oscal-1-1-3"
)

// NilIfEmpty returns nil if the slice is empty, otherwise returns the original slice.
func NilIfEmpty[T any](slice *[]T) *[]T {
	if slice == nil || len(*slice) == 0 {
		return nil
	}
	return slice
}

func FindValuesByName(model *oscalTypes.OscalModels, name string) []string {
	var results []string
	seen := make(map[uintptr]bool)
	var walk func(val reflect.Value, key string)
	walk = func(val reflect.Value, key string) {
		if !val.IsValid() {
			return
		}
		for (val.Kind() == reflect.Ptr || val.Kind() == reflect.Interface) && !val.IsNil() {
			if val.Kind() == reflect.Ptr {
				ptr := val.Pointer()
				if seen[ptr] {
					return
				}
				seen[ptr] = true
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

func HasDuplicateValuesByName(model *oscalTypes.OscalModels, name string) bool {
	values := FindValuesByName(model, name)
	valueMap := make(map[string]bool)
	for _, value := range values {
		if valueMap[value] {
			return false
		}
		valueMap[value] = true
	}
	return true
}
