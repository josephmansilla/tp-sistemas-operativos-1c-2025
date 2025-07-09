package logger

import "errors"

var (
	ErrBadRequest      = errors.New("Bad Request")
	ErrInternalFailure = errors.New("Internal Failure")
	ErrNotFound        = errors.New("Not Found")
	ErrIsEmpty         = errors.New("Empty")
	ErrNoInstance      = errors.New("No Instance Found")
	ErrSegmentFault    = errors.New("Segment Fault")
	ErrNoMemory        = errors.New("No memory")
	ErrNoFreePages     = errors.New("No pages available")
)
