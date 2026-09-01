package lsm

import "fmt"

type Options struct {
	DbName           string
	SkipListMaxLevel int
	SkipListMaxSize  int
}

func (opts *Options) validate() error {
	if opts == nil {
		return fmt.Errorf("LSM options are required")
	}
	if opts.DbName == "" {
		return fmt.Errorf("database name is required")
	}
	if opts.SkipListMaxLevel <= 0 {
		return fmt.Errorf("skip list max level must be positive")
	}
	if opts.SkipListMaxSize <= 0 {
		return fmt.Errorf("skip list max size must be positive")
	}

	return nil
}
