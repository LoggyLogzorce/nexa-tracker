package comment

import "errors"

var (
	ErrUserNotFound    = errors.New("user not found")
	ErrCommentNotFound = errors.New("comment not found")
	ErrNotCommentOwner = errors.New("user is not the comment owner")
)
