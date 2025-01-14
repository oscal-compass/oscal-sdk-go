/*
 Copyright 2024 The OSCAL Compass Authors
 SPDX-License-Identifier: Apache-2.0
*/

package set

// Set represents a set data structure.
type Set[T comparable] map[T]struct{}

// New NewSet returns an initialized set.
func New[T comparable]() Set[T] {
	return make(Set[T])
}

// Add adds item into the set s.
func (s Set[T]) Add(item T) {
	s[item] = struct{}{}
}

// Has checks if the set contains an item.
func (s Set[T]) Has(item T) bool {
	_, ok := s[item]
	return ok
}

// Intersect returns a new Set representing the intersection
// between two sets.
func (s Set[T]) Intersect(other Set[T]) Set[T] {
	newSet := New[T]()
	for elem := range s {
		if _, ok := other[elem]; ok {
			newSet.Add(elem)
		}
	}
	return newSet
}
