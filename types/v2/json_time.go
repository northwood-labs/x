package types

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

type (
	// UnixTime represents a time.Time object that is serialized as a Unix timestamp.
	UnixTime struct {
		time.Time
	}

	// RFC3339Time represents a time.Time object that is serialized as an RFC3339 timestamp.
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

// MarshalJSON turns our time.Time back into an integer.
func (u UnixTime) MarshalJSON() ([]byte, error) { // lint:allow_param
	return []byte(fmt.Sprintf("%d", (u.Time.Unix()))), nil
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

// MarshalJSON turns our time.Time back into an integer.
func (r *RFC3339Time) MarshalJSON() ([]byte, error) { // lint:allow_param
	return []byte(
		fmt.Sprintf("%q", r.UTC().Format(time.RFC3339)),
	), nil
}
