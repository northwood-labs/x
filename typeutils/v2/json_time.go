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

package typeutils

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type (
	// UnixTime wraps time.Time so that JSON (un)marshaling uses a Unix
	// epoch integer instead of the default RFC 3339 string. Many APIs
	// (GitHub, AWS, POSIX tools) expose timestamps as integers; this
	// type lets consumers decode them directly into a time.Time without
	// manual conversion at every call site.
	UnixTime struct {
		time.Time
	}

	// RFC3339Time wraps time.Time so that JSON (un)marshaling always
	// uses RFC 3339 format with explicit timezone offset. Unlike the
	// default encoding/json behavior (which already uses RFC 3339),
	// this type strips surrounding quotes during Unmarshal and
	// normalizes to UTC on Marshal, providing consistent round-trip
	// behavior for APIs that return timezone-aware strings.
	RFC3339Time struct {
		time.Time
	}
)

// UnmarshalJSON decodes a JSON integer (Unix epoch seconds) into the
// embedded time.Time. This satisfies the json.Unmarshaler interface so
// the standard json.Unmarshal call handles UnixTime fields
// automatically.
func (u *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64

	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return fmt.Errorf("could not unmarshal the timestamp into an integer: %w", err)
	}

	u.Time = time.Unix(timestamp, 0)

	return nil
}

// MarshalJSON encodes the time as a bare integer (Unix epoch seconds)
// so downstream consumers receive the same integer format they sent.
func (u UnixTime) MarshalJSON() ([]byte, error) { // lint:allow_param
	return fmt.Appendf(nil, "%d", u.Unix()), nil
}

// UnmarshalJSON decodes a JSON string in RFC 3339 format into the
// embedded time.Time. Surrounding quotes are stripped before parsing
// because encoding/json delivers the raw token including quotes.
func (r *RFC3339Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	pt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("failed to parse RFC3339 time: %w", err)
	}

	*r = RFC3339Time{pt}

	return nil
}

// MarshalJSON encodes the time as a quoted RFC 3339 string normalized
// to UTC, ensuring consistent output regardless of the original
// timezone.
func (r *RFC3339Time) MarshalJSON() ([]byte, error) { // lint:allow_param
	return fmt.Appendf(nil, "%q", r.UTC().Format(time.RFC3339)), nil
}
