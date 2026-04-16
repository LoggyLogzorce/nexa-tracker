package priority

import "errors"

var (
	ErrPriorityNotFound    = errors.New("priority not found")
	ErrPriorityTitleExists = errors.New("priority with this title already exists")
	ErrColorFormat         = errors.New("invalid color format")
	ErrProjectNotFound     = errors.New("project not found")
	ErrProjectAccessDenied = errors.New("access to project denied")
)
