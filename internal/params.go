package internal

// CreateParams defines the arguments used for creating Task records.
type CreateParams struct {
	Description string
	Priority    Priority
	Dates       Dates
}

// Validate indicates whether the fields are valid or not.
func (c CreateParams) Validate() error {
	t := Task{
		Description: c.Description,
		Priority:    c.Priority,
		Dates:       c.Dates,
	}

	if err := t.Validate(); err != nil {
		return WrapErrorf(err, ErrorCodeInvalidArgument, "Validate")
	}

	return nil
}

//-

// SearchParams defines the arguments used for searching Task records.
type SearchParams struct {
	Description *string
	Priority    *Priority
	IsDone      *bool
	From        int64
	Size        int64
}

// IsZero determines whether the search arguments have values or not.
func (a SearchParams) IsZero() bool {
	return a.Description == nil &&
		a.Priority == nil &&
		a.IsDone == nil
}

// SearchResults defines the collection of tasks that were found.
type SearchResults struct {
	Tasks []Task
	Total int64
}