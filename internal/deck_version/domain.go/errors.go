package domain

import "errors"

var ErrDeckIdRequired = errors.New("deck id is required")

var (
	ErrDeckIdInvalid            = errors.New("deck id must be a valid ObjectID")
	ErrDeckNotFound             = errors.New("deck not found")
	ErrDeckVersionAlreadyExists = errors.New("deck version already exists")
)
