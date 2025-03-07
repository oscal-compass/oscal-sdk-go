/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package models

// NilIfEmpty returns nil if the slice is empty, otherwise returns the original slice.
func NilIfEmpty[T any](slice *[]T) *[]T {
	if slice == nil || len(*slice) == 0 {
		return nil
	}
	return slice
}
