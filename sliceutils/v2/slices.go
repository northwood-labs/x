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

// ComparableT is a list of types that can be compared that is semi-broken in
// 1.18 and 1.19, but fixed in 1.20.
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

// Dedupe will remove duplicate values from a variety of types.
//
// `comparable` is (supposed to be) a type that can be compared, but there were
// issues in 1.18 and 1.19 that were resolved in 1.20.
// https://blog.carlmjohnson.net/post/2023/golang-120-language-changes/
//
// ~rune is an alias of ~int32. ~byte is an alias of ~uint8.
//
// https://hackernoon.com/how-to-remove-duplicates-in-go-slices
func Dedupe[T ComparableT](s []T) []T {
	// If there are 0 or 1 items we return the slice itself.
	if len(s) < 2 { // lint:allow_raw_numbers
		return s
	}

	// Make the slice case-insensitive, ascending, sorted.
	slices.SortStableFunc(s, cmp.Compare)

	uniqPointer := 0

	for i := 1; i < len(s); i++ {
		// Compare a current item with the item under the unique pointer. If
		// they are not the same, write the item next to the right of the unique
		// pointer.
		if s[uniqPointer] != s[i] {
			uniqPointer++

			s[uniqPointer] = s[i]
		}
	}

	return s[:uniqPointer+1]
}

// StringSliceToHashmap will invert the values of a slice to keys in a hashmap.
func StringSliceToHashmap(slice []string) map[string]struct{} {
	hashmap := make(map[string]struct{})

	for i := range slice {
		hashmap[slice[i]] = struct{}{}
	}

	return hashmap
}

// FilterSubstr will filter a slice of any type by a substring. The callback is
// used to convert the type to a string for comparison.
func FilterSubstr[T any](u []T, s string, callback func(T) string) []T {
	out := make([]T, 0, len(u))

	for i := range u {
		if strings.Contains(callback(u[i]), s) {
			out = append(out, u[i])
		}
	}

	return out
}

// FilterRegex will filter a slice of any type by a regular expression. The
// callback is used to convert the type to a string for comparison.
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
