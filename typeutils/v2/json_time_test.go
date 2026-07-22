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
	"strings"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

type (
	testRFC3339JSON struct {
		Timestamp RFC3339Time `json:"timestamp"`
	}

	testUnixJSON struct {
		Timestamp UnixTime `json:"timestamp"`
	}
)

func tz(z string) *time.Location {
	loc, err := time.LoadLocation(z)
	if err != nil {
		panic(err)
	}

	return loc
}

func TestRFC3339Time(t *testing.T) {
	type testCase struct {
		expected time.Time
		input    string
	}

	tests := []testCase{
		{
			input:    `{"timestamp":"2018-01-01T00:00:00Z"}`,
			expected: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC),
		},
		{
			input:    `{"timestamp":"2018-01-01T00:00:00-05:00"}`,
			expected: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("America/New_York")),
		},
		{
			input:    `{"timestamp":"2018-01-01T00:00:00-06:00"}`,
			expected: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("America/Chicago")),
		},
		{
			input:    `{"timestamp":"2018-01-01T00:00:00-07:00"}`,
			expected: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("America/Denver")),
		},
		{
			input:    `{"timestamp":"2018-01-01T00:00:00-08:00"}`,
			expected: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("America/Los_Angeles")),
		},
		{
			input:    `{"timestamp":"2018-01-01T00:00:00+00:00"}`,
			expected: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("Europe/London")),
		},
		{
			input:    `{"timestamp":"2018-01-01T00:00:00+01:00"}`,
			expected: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("Europe/Paris")),
		},
		{
			input:    `{"timestamp":"2018-01-01T00:00:00+09:00"}`,
			expected: time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("Asia/Tokyo")),
		},
	}

	for i := range tests {
		tc := tests[i]

		actual := testRFC3339JSON{}

		err := json.Unmarshal([]byte(tc.input), &actual)
		if err != nil {
			t.Errorf("Failed to parse %v; %v", tc.input, err)
		}

		// Test that the time deserialized into a time.Time object.
		diff := cmp.Diff(tc.expected.Format(time.RFC3339), actual.Timestamp.Format(time.RFC3339))
		if diff != "" {
			t.Error(diff)
		}

		reserialized, err := json.Marshal(actual)
		if err != nil {
			t.Errorf("Failed to parse %v; %v", tc.input, err)
		}

		// Normalize UTC time with Zulu identifier to +00:00.
		inputWithFixedTZ := strings.Replace(tc.input, "+00:00", "Z", 1)
		reserializedWithFixedTZ := strings.Replace(string(reserialized), "+00:00", "Z", 1)

		// Test that the time deserialized and reserialized into RFC 3339 format.
		diff = cmp.Diff(reserializedWithFixedTZ, inputWithFixedTZ)
		if diff != "" {
			t.Error(diff)
		}
	}
}

func TestUnixTime(t *testing.T) {
	type testCase struct {
		expected time.Time
		input    string
	}

	tests := []testCase{
		{input: `{"timestamp":0}`, expected: time.Date(1970, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{input: `{"timestamp":315532800}`, expected: time.Date(1980, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{input: `{"timestamp":631152000}`, expected: time.Date(1990, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{input: `{"timestamp":946684800}`, expected: time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{input: `{"timestamp":1262304000}`, expected: time.Date(2010, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{input: `{"timestamp":1577836800}`, expected: time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
	}

	for i := range tests {
		tc := tests[i]

		actual := testUnixJSON{}

		err := json.Unmarshal([]byte(tc.input), &actual)
		if err != nil {
			t.Errorf("Failed to parse %v; %v", tc.input, err)
		}

		// Test that the time deserialized into a time.Time object.
		diff := cmp.Diff(tc.expected.Format(time.UnixDate), actual.Timestamp.Time.UTC().Format(time.UnixDate))
		if diff != "" {
			t.Error(diff)
		}

		reserialized, err := json.Marshal(actual)
		if err != nil {
			t.Errorf("Failed to parse %v; %v", tc.input, err)
		}

		// Test that the time deserialized and reserialized into RFC 3339 format.
		diff = cmp.Diff(string(reserialized), tc.input)
		if diff != "" {
			t.Error(diff)
		}
	}
}
