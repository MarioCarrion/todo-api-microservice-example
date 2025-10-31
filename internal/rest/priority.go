package rest

import (
	"encoding/json"

	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

// NewPriority converts the received domain type to a rest type, when the argument is unknown "none" is used.
func NewPriority(p internal.Priority) Priority {
	switch p {
	case internal.PriorityNone:
		return PriorityNone
	case internal.PriorityLow:
		return PriorityLow
	case internal.PriorityMedium:
		return PriorityMedium
	case internal.PriorityHigh:
		return PriorityHigh
	}

	return PriorityNone
}

// Convert returns the domain type defining the internal representation, when Priority is unknown "none" is
// used.
func (p *Priority) Convert() internal.Priority {
	switch *p {
	case PriorityNone:
		return internal.PriorityNone
	case PriorityLow:
		return internal.PriorityLow
	case PriorityMedium:
		return internal.PriorityMedium
	case PriorityHigh:
		return internal.PriorityHigh
	}

	return internal.PriorityNone
}

// Validate ...
func (p *Priority) Validate() error {
	switch *p {
	case PriorityNone, PriorityLow, PriorityMedium, PriorityHigh:
		return nil
	}

	return internal.NewErrorf(internal.ErrorCodeInvalidArgument, "unknown value")
}

// MarshalJSON ...
func (p *Priority) MarshalJSON() ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "Validate")
	}

	b, err := json.Marshal(string(*p))
	if err != nil {
		return nil, internal.WrapErrorf(err, internal.ErrorCodeUnknown, "json.Marshal")
	}

	return b, nil
}

// UnmarshalJSON ...
func (p *Priority) UnmarshalJSON(b []byte) error {
	var str string
	if err := json.Unmarshal(b, &str); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "json.Unmarshal")
	}

	conv := Priority(str)

	if err := conv.Validate(); err != nil {
		return internal.WrapErrorf(err, internal.ErrorCodeInvalidArgument, "Validate")
	}

	*p = Priority(str)

	return nil
}
