package database

import "github.com/willie-cadete/lxdexplorer-collector/config"

type Database struct {
	config *config.Config
}

type Options struct {
	Config *config.Config
}

func NewDatabase(opts Options) *Database {
	return &Database{
		config: opts.Config,
	}
}
