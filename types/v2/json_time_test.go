package types

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
	loc, _ := time.LoadLocation(z)
	return loc
}

func TestRFC3339Time(t *testing.T) {
	type testCase struct {
		input    string
		expected time.Time
	}

	tests := []testCase{
		{`{"timestamp":"2018-01-01T00:00:00Z"}`, time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{`{"timestamp":"2018-01-01T00:00:00-05:00"}`, time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("America/New_York"))},
		{`{"timestamp":"2018-01-01T00:00:00-06:00"}`, time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("America/Chicago"))},
		{`{"timestamp":"2018-01-01T00:00:00-07:00"}`, time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("America/Denver"))},
		{`{"timestamp":"2018-01-01T00:00:00-08:00"}`, time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("America/Los_Angeles"))},
		{`{"timestamp":"2018-01-01T00:00:00+00:00"}`, time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("Europe/London"))},
		{`{"timestamp":"2018-01-01T00:00:00+01:00"}`, time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("Europe/Paris"))},
		{`{"timestamp":"2018-01-01T00:00:00+09:00"}`, time.Date(2018, time.Month(1), 1, 0, 0, 0, 0, tz("Asia/Tokyo"))},
	}

	for i := range tests {
		tc := tests[i]

		actual := testRFC3339JSON{}

		err := json.Unmarshal([]byte(tc.input), &actual)
		if err != nil {
			t.Errorf("Failed to parse %v; %v", tc.input, err)
		}

		// Test that the time deserialized into a time.Time object.
		diff := cmp.Diff(tc.expected.Format(time.RFC3339), actual.Timestamp.Time.Format(time.RFC3339))
		if diff != "" {
			t.Errorf(diff)
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
			t.Errorf(diff)
		}
	}
}

func TestUnixTime(t *testing.T) {
	type testCase struct {
		input    string
		expected time.Time
	}

	tests := []testCase{
		{`{"timestamp":0}`, time.Date(1970, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{`{"timestamp":315532800}`, time.Date(1980, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{`{"timestamp":631152000}`, time.Date(1990, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{`{"timestamp":946684800}`, time.Date(2000, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{`{"timestamp":1262304000}`, time.Date(2010, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
		{`{"timestamp":1577836800}`, time.Date(2020, time.Month(1), 1, 0, 0, 0, 0, time.UTC)},
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
			t.Errorf(diff)
		}

		reserialized, err := json.Marshal(actual)
		if err != nil {
			t.Errorf("Failed to parse %v; %v", tc.input, err)
		}

		// Test that the time deserialized and reserialized into RFC 3339 format.
		diff = cmp.Diff(string(reserialized), tc.input)
		if diff != "" {
			t.Errorf(diff)
		}
	}
}
