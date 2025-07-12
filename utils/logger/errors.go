package logger

import "errors"

var (
	ErrBadRequest      = errors.New("bad Request")
	ErrInternalFailure = errors.New("internal Failure")
	ErrNotFound        = errors.New("not Found")
	ErrIsEmpty         = errors.New("empty")
	ErrNoInstance      = errors.New("no Instance Found")
	ErrSegmentFault    = errors.New("segment Fault")
	ErrNoMemory        = errors.New("no memory")
	ErrNoFreePages     = errors.New("no pages available")
	ErrProcessNil      = errors.New("process requested is nil")
	ErrNoTabla         = errors.New("no Tabla Found")
	ErrNoIndices       = errors.New("no Indices Found")
	ErrNotPresent      = errors.New("entrada not present in memory")
)
