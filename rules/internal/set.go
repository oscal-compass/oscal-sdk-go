// SPDX-License-Identifier: Apache-2.0

package internal

// Set represents a set data structure.
type Set[T comparable] map[T]struct{}

// NewSet returns an initialized set.
func NewSet[T comparable]() Set[T] {
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
