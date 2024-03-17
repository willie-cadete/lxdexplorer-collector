package config

// Config is a struct that holds the configuration of the application
type Config struct {
	Collector Collector `mapstructure:"collector"`
	Logging   Logging   `mapstructure:"logging"`
	Database  Database  `mapstructure:"database"`
	LXD       LXD       `mapstructure:"lxd"`
}

// Collector is a struct that holds the configuration of the collector
type Collector struct {
	Interval  int `mapstructure:"interval"`
	Retention int `mapstructure:"retention"`
}

// Logging is a struct that holds the configuration of the logging
type Logging struct {
	Level string `mapstructure:"level"`
}

// Database is a struct that holds the configuration of the database
type Database struct {
	URI string `mapstructure:"uri"`
}

// LXD is a struct that holds the configuration of the LXD
type LXD struct {
	TLS       TLS      `mapstructure:"tls"`
	Hostnodes []string `mapstructure:"hostnodes"`
}

// TLS is a struct that holds the configuration of the TLS
type TLS struct {
	Cert   string `mapstructure:"certificate"`
	Key    string `mapstructure:"key"`
	Verify bool   `mapstructure:"verify"`
}
