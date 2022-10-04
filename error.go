package short

import "fmt"

type ConflictError struct{}

func (e *ConflictError) Error() string {
	return "conflict - id already exist in collection"
}

type IdNotFoundError struct {
	id string
}

func (e *IdNotFoundError) Error() string {
	return fmt.Sprintf("the id %s not found", e.id)
}
