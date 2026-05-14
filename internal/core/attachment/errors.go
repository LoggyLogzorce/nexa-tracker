package attachment

import "errors"

var (
	ErrCreateDirectory        = errors.New("error creating directory")
	ErrWriteFile              = errors.New("error writing file")
	ErrCreateAttachmentRecord = errors.New("error creating attachment")
	ErrRemoveFile             = errors.New("error removing file")
	ErrAttachmentNotFound     = errors.New("attachment not found")
	ErrNotAttachmentOwner     = errors.New("user is not the attachment owner")
)
