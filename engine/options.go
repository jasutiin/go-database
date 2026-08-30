package engine

import "errors"

type Options struct {
	DbName           string
	SkipListMaxLevel int
	SkipListMaxSize  int
}

var ErrDbNameRequired = errors.New("DbName is required for options")
var ErrMaxLevelLTEZero = errors.New("Skip List max level cannot be less than or equal to zero")
var ErrMaxSizeLTEZero = errors.New("Skip List max size cannot be less than or equal to zero")

func (opts *Options) Validate() error {
	if opts.DbName == "" {
		return ErrDbNameRequired
	}

	if opts.SkipListMaxLevel <= 0 {
		return ErrMaxLevelLTEZero
	}

	if opts.SkipListMaxSize <= 0 {
		return ErrMaxSizeLTEZero
	}

	return nil
}
