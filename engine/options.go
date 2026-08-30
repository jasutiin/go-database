package engine

import "errors"

type StorageType string

const (
	StorageLSM   StorageType = "lsm"
	StorageBTree StorageType = "btree"
)

type Options struct {
	DbName           string
	StorageType      StorageType
	SkipListMaxLevel int
	SkipListMaxSize  int
}

var ErrDbNameRequired = errors.New("DbName is required for options")
var ErrMaxLevelLTEZero = errors.New("Skip List max level cannot be less than or equal to zero")
var ErrMaxSizeLTEZero = errors.New("Skip List max size cannot be less than or equal to zero")
var ErrUnknownStorageType = errors.New("unknown storage type")

func (opts *Options) Validate() error {
	if opts.DbName == "" {
		return ErrDbNameRequired
	}

	switch opts.StorageType {
	case StorageLSM:
		if opts.SkipListMaxLevel <= 0 {
			return ErrMaxLevelLTEZero
		}

		if opts.SkipListMaxSize <= 0 {
			return ErrMaxSizeLTEZero
		}

		return nil

	case StorageBTree:
		return nil

	default:
		return ErrUnknownStorageType
	}
}
