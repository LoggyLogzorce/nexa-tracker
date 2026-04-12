package status

import "errors"

var (
	ErrProjectNotFound     = errors.New("project not found")
	ErrProjectAccessDenied = errors.New("access to project denied")
	ErrStatusNameExists    = errors.New("status with this name already exists")
	ErrDuplicateOrderIndex = errors.New("status with this order_index already exists")
)
