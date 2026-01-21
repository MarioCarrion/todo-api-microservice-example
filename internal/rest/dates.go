package rest

import (
	"github.com/MarioCarrion/todo-api-microservice-example/internal"
)

// NewDates ...
func NewDates(d internal.Dates) Dates {
	return Dates{
		Start: d.Start,
		Due:   d.Due,
	}
}

// Convert returns the domain type defining the internal representation.
func (d Dates) Convert() internal.Dates {
	return internal.Dates{
		Start: d.Start,
		Due:   d.Due,
	}
}
