package task

import "errors"

var (
	ErrDataIntegrity        = errors.New("data integrity error")
	ErrAssigneeNotInProject = errors.New("assignee not in project")
	ErrStatusNotInProject   = errors.New("status not in project")
	ErrPriorityNotInProject = errors.New("priority not in project")
)
