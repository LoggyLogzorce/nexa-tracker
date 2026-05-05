package participant

import "errors"

var (
	ErrParticipantIDExists   = errors.New("participant with this id already exists")
	ErrParticipantRoleFormat = errors.New("invalid role format")
)
