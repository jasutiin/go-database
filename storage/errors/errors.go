package errors

import "errors"

var (
	ErrDbNameRequired     = errors.New("DbName is required for options")
	ErrMaxLevelLTEZero    = errors.New("Skip List max level cannot be less than or equal to zero")
	ErrMaxSizeLTEZero     = errors.New("Skip List max size cannot be less than or equal to zero")
	ErrUnknownStorageType = errors.New("unknown storage type")
	ErrKeyNotFound        = errors.New("key not found")
)
