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

// Package typeutils provides custom types that extend Go's standard
// library with serialization behaviors the stdlib does not offer out
// of the box. The primary use case is JSON time handling: APIs return
// timestamps as either Unix integers or RFC 3339 strings, but Go's
// time.Time only supports RFC 3339 natively via encoding/json. These
// wrapper types let structs declare the wire format declaratively via
// their type, so JSON marshaling and unmarshaling happen correctly
// without custom per-field logic in every consumer.
package typeutils
