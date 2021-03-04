package rest

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/MarioCarrion/todo-api/internal"
)

// Dates indicates a point in time where a task starts or completes, dates are not enforced on Tasks.
type Dates struct {
	Start Time `json:"start"`
	Due   Time `json:"due"`
}

// NewDates ...
func NewDates(d internal.Dates) Dates {
	return Dates{
		Start: Time(d.Start),
		Due:   Time(d.Due),
	}
}

// Convert returns the domain type defining the internal representation.
func (d Dates) Convert() internal.Dates {
	return internal.Dates{
		Start: time.Time(d.Start),
		Due:   time.Time(d.Due),
	}
}

// Time represents an instant in time, JSON are strings using RFC3339.
type Time time.Time

// MarshalJSON ...
func (t Time) MarshalJSON() ([]byte, error) {
	str := time.Time(t).Format(time.RFC3339)

	b, err := json.Marshal(str)
	if err != nil {
		return nil, fmt.Errorf("json marshal: %w", err)
	}

	return b, nil
}

// UnmarshalJSON ...
func (t *Time) UnmarshalJSON(b []byte) error {
	var s string
	if err := json.Unmarshal(b, &s); err != nil {
		return fmt.Errorf("json unmarshal: %w", err)
	}

	tt, err := time.Parse(time.RFC3339, s)
	if err != nil {
		return fmt.Errorf("convert: %w", err)
	}

	*t = Time(tt)

	return nil
}
