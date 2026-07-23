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

// Package structutils provides reflection-based converters that turn
// arbitrary Go structs into maps or flat key-value slices. Structured
// logging (slog) and generic serialization layers need data as
// map[string]any or []any rather than typed structs. Writing manual
// conversion functions for every struct is tedious and error-prone;
// this package uses reflection to handle any struct shape — including
// nested and pointer fields — so callers get correct, recursive
// conversion without boilerplate.
package structutils
