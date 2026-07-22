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
	// UnixTime represents a [time.Time] object that is serialized as a Unix
	// timestamp.
	UnixTime struct {
		time.Time
	}

	// RFC3339Time represents a [time.Time] object that is serialized as an
	// RFC3339 timestamp.
	RFC3339Time struct {
		time.Time
	}
)

// UnmarshalJSON is the method that satisfies the Unmarshaller interface.
func (u *UnixTime) UnmarshalJSON(b []byte) error {
	var timestamp int64

	err := json.Unmarshal(b, &timestamp)
	if err != nil {
		return fmt.Errorf("could not unmarshal the timestamp into an integer: %w", err)
	}

	u.Time = time.Unix(timestamp, 0)

	return nil
}

// MarshalJSON turns our [time.Time] back into an integer.
func (u UnixTime) MarshalJSON() ([]byte, error) { // lint:allow_param
	return fmt.Appendf(nil, "%d", u.Unix()), nil
}

// UnmarshalJSON is the method that satisfies the Unmarshaller interface.
func (r *RFC3339Time) UnmarshalJSON(b []byte) error {
	s := strings.Trim(string(b), "\"")

	pt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("failed to parse RFC3339 time: %w", err)
	}

	*r = RFC3339Time{pt}

	return nil
}

// MarshalJSON turns our [time.Time] back into an integer.
func (r *RFC3339Time) MarshalJSON() ([]byte, error) { // lint:allow_param
	return fmt.Appendf(nil, "%q", r.UTC().Format(time.RFC3339)), nil
}
