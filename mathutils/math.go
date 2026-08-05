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

package mathutils

import (
	"cmp"

	"golang.org/x/exp/constraints"
)

// Number is a union type as a constraint for numeric types.
//
// License: MIT
// Adapted from https://github.com/Goldziher/go-utils/blob/cebd170/mathutils/mathutils.go
type Number interface {
	constraints.Integer | constraints.Float
}

// Clamp restricts a value to be within a specified range. If value < min,
// returns min. If value > max, returns max. Otherwise returns value.
//
// License: MIT
// Adapted from https://github.com/Goldziher/go-utils/blob/cebd170/mathutils/mathutils.go
func Clamp[T cmp.Ordered](value, min, max T) T {
	if value < min {
		return min
	}

	if value > max {
		return max
	}

	return value
}

// InRange checks if a value is within a range [min, max] (inclusive).
//
// License: MIT
// Adapted from https://github.com/Goldziher/go-utils/blob/cebd170/mathutils/mathutils.go
func InRange[T cmp.Ordered](value, min, max T) bool {
	return value >= min && value <= max
}
