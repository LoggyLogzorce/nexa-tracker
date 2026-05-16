package project

import "errors"

var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrProjectAccessDenied = errors.New("access to project denied")
	ErrGetOwner            = errors.New("error getting owner")
	ErrInvalidStatus       = errors.New("invalid status")
	ErrInvalidPriority     = errors.New("invalid priority")
	ErrDataIntegrity       = errors.New("data integrity error")
)
