package engine

import engineerrors "github.com/jasutiin/go-database/engine/errors"

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

var ErrDbNameRequired = engineerrors.ErrDbNameRequired
var ErrMaxLevelLTEZero = engineerrors.ErrMaxLevelLTEZero
var ErrMaxSizeLTEZero = engineerrors.ErrMaxSizeLTEZero
var ErrUnknownStorageType = engineerrors.ErrUnknownStorageType
var ErrKeyNotFound = engineerrors.ErrKeyNotFound

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
