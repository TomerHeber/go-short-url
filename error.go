package short

type ConflictError struct{}

func (e *ConflictError) Error() string {
	return "conflict - id already exist in collection"
}
