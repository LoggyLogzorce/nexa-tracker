package project

import "errors"

var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrProjectAccessDenied = errors.New("access to project denied")
)
