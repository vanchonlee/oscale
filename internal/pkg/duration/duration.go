// Package duration provides the logic to parse and evaluate durations
package duration

import (
	"time"
)

// Duration represents a duration in string format
type Duration struct {
	DurationStr string `json:"duration,omitempty"`

	val time.Duration `json:"-"`
}

// Duration returns the time.Duration value of the duration string
func (d *Duration) Duration() (time.Duration, error) {
	if d.DurationStr == "" {
		return 0, nil
	}

	if d.val == 0 {
		var err error
		d.val, err = time.ParseDuration(d.DurationStr)
		if err != nil {
			return 0, err
		}
	}
	return d.val, nil
}

// MustDuration returns the time.Duration value of the duration string or panics
func (d *Duration) MustDuration() time.Duration {
	duration, err := d.Duration()
	if err != nil {
		panic(err)
	}
	return duration
}
