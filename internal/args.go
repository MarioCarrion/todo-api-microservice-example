package internal

// SearchArgs defines the arguments used for searching Task records.
type SearchArgs struct {
	Description *string
	Priority    *Priority
	IsDone      *bool
	From        int64
	Size        int64
}

// IsZero determines whether the search arguments have values or not.
func (a SearchArgs) IsZero() bool {
	return a.Description == nil &&
		a.Priority == nil &&
		a.IsDone == nil
}

// SearchResults defines the collection of tasks that were found.
type SearchResults struct {
	Tasks []Task
	Total int64
}
