package engine

import "errors"

type Options struct {
	DbName           string
	SkipListMaxLevel int
	SkipListMaxSize  int
}

var ErrDbNameRequired = errors.New("DbName is required for options")
var ErrMaxLevelBelowZero = errors.New("Skip List max level cannot be less than zero")
var ErrMaxSizeBelowZero = errors.New("Skip List max size cannot be less than zero")

func (opts *Options) Validate() error {
	if opts.DbName == "" {
		return ErrDbNameRequired
	}

	if opts.SkipListMaxLevel < 0 {
		return ErrMaxLevelBelowZero
	}

	if opts.SkipListMaxSize < 0 {
		return ErrMaxSizeBelowZero
	}

	return nil
}
