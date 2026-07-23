// Copyright 2026, Northwood Labs, LLC <license@northwood-labs.com>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package sliceutils

import (
	"cmp"
	"regexp"
	"slices"
	"strings"
)

// ComparableT enumerates the concrete types Dedupe can operate on.
// Go 1.18/1.19 had compiler bugs around the built-in comparable
// constraint with certain type sets, so this explicit union ensures
// correct behavior across all supported Go versions (1.20+). Using a
// tilde (~) on each type permits named types derived from primitives.
type ComparableT interface {
	~byte |
		~float32 |
		~float64 |
		~int |
		~int16 |
		~int64 |
		~int8 |
		~rune |
		~string |
		~uint |
		~uint16 |
		~uint32 |
		~uint64 |
		~uintptr
}

// Dedupe removes duplicate values from a slice, returning a sorted,
// unique result. Many CLI inputs (cloud resource lists, config keys)
// arrive with duplicates from multiple sources; deduplication prevents
// redundant API calls and duplicate output lines. The in-place
// algorithm avoids allocating a second slice, which matters when lists
// are large (e.g., thousands of cloud resources).
//
// ~rune is an alias of ~int32. ~byte is an alias of ~uint8.
//
// https://hackernoon.com/how-to-remove-duplicates-in-go-slices
func Dedupe[T ComparableT](s []T) []T {
	// Short-circuit: a slice with 0 or 1 elements is already unique.
	if len(s) < 2 { // lint:allow_raw_numbers
		return s
	}

	// Sorting first allows a single linear pass to collapse duplicates
	// by comparing adjacent elements, achieving O(n log n) overall
	// instead of the O(n²) of a nested-loop approach.
	slices.SortStableFunc(s, cmp.Compare)

	uniqPointer := 0

	for i := 1; i < len(s); i++ {
		// Advance the unique pointer only when we encounter a new
		// distinct value, effectively compacting the slice in-place.
		if s[uniqPointer] != s[i] {
			uniqPointer++

			s[uniqPointer] = s[i]
		}
	}

	return s[:uniqPointer+1]
}

// StringSliceToHashmap converts a string slice into a set (map with
// empty-struct values) for O(1) membership lookups. This is used when
// a command needs to quickly test whether a value exists in a
// potentially large allow-list or deny-list without scanning the slice
// on each check.
func StringSliceToHashmap(slice []string) map[string]struct{} {
	hashmap := make(map[string]struct{})

	for i := range slice {
		hashmap[slice[i]] = struct{}{}
	}

	return hashmap
}

// FilterSubstr returns elements whose string representation contains a
// given substring. The callback lets callers define which field(s) of a
// complex type are searchable, decoupling the filter logic from any
// particular struct layout. This powers interactive search/filter UIs
// where the user types a partial name.
func FilterSubstr[T any](u []T, s string, callback func(T) string) []T {
	out := make([]T, 0, len(u))

	for i := range u {
		if strings.Contains(callback(u[i]), s) {
			out = append(out, u[i])
		}
	}

	return out
}

// FilterRegex returns elements whose string representation matches a
// regular expression. Like FilterSubstr, the callback decouples field
// selection from filtering. Regex is offered alongside substring
// filtering because power users need anchors and alternation for
// precise filtering (e.g., "^prod-" to match only production
// resources).
func FilterRegex[T any](u []T, r string, callback func(T) string) []T {
	out := make([]T, 0, len(u))
	re := regexp.MustCompile(r)

	for i := range u {
		if re.MatchString(callback(u[i])) {
			out = append(out, u[i])
		}
	}

	return out
}
