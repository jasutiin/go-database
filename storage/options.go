package storage

import errs "github.com/jasutiin/go-database/storage/errors"

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

func (opts *Options) Validate() error {
	if opts.DbName == "" {
		return errs.ErrDbNameRequired
	}

	switch opts.StorageType {
	case StorageLSM:
		if opts.SkipListMaxLevel <= 0 {
			return errs.ErrMaxLevelLTEZero
		}

		if opts.SkipListMaxSize <= 0 {
			return errs.ErrMaxSizeLTEZero
		}

		return nil

	case StorageBTree:
		return nil

	default:
		return errs.ErrUnknownStorageType
	}
}
