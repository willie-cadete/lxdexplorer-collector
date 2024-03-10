package config

import (
	"fmt"
	"strings"

	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

const (
	defaultCollectorInterval  = 60
	defaultCollectorRetention = 60
	defaultLoggingLevel       = "info"
	defaultDatabaseURI        = "mongodb://localhost:27017"
	defaultLXDTLSCertificate  = "./tls/client.crt"
	defaultLXDTLSKey          = "./tls/client.key"
	defaultLXDTLSVerify       = false
)

func DefaultHostnodes() []string {
	return []string{"https://localhost:8443"}
}

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

func (c *Config) GetHostnodes() []string {
	return c.LXD.Hostnodes
}

func (c *Config) GetTLS() TLS {
	return c.LXD.TLS
}

func (c *Config) GetDatabaseURI() string {
	return c.Database.URI
}

func (c *Config) GetCollectorInterval() int {
	return c.Collector.Interval
}

func (c *Config) GetCollectorRetention() int {
	return c.Collector.Retention
}

func (c *Config) GetLoggingLevel() string {
	return c.Logging.Level
}

func (c *Config) GetLXDTLSCertificate() string {
	return c.LXD.TLS.Cert
}

func (c *Config) GetLXDTLSKey() string {
	return c.LXD.TLS.Key
}

func (c *Config) GetLXDTLSVerify() bool {
	return c.LXD.TLS.Verify
}

func setDefaults() {
	viper.SetDefault("collector.interval", defaultCollectorInterval)
	viper.SetDefault("collector.retention", defaultCollectorRetention)
	viper.SetDefault("logging.level", defaultLoggingLevel)
	viper.SetDefault("database.uri", defaultDatabaseURI)
	viper.SetDefault("lxd.tls.certificate", defaultLXDTLSCertificate)
	viper.SetDefault("lxd.tls.key", defaultLXDTLSKey)
	viper.SetDefault("lxd.tls.verify", defaultLXDTLSVerify)
	viper.SetDefault("lxd.hostnodes", DefaultHostnodes())
}

func LoadConfig(fp string) (*Config, error) {

	viper.SetConfigType("yaml")
	viper.SetConfigName("config")
	viper.AddConfigPath(fp)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Warnf("Config file not found; using default values")
			setDefaults()
		} else {
			log.Errorf("Error reading config file, %s", err)
		}
	} else {
		log.Infof("Using config file: %s", viper.ConfigFileUsed())
	}

	for _, key := range viper.AllKeys() {
		envKey := strings.ToUpper(strings.ReplaceAll(key, ".", "_"))
		err := viper.BindEnv(key, envKey)
		if err != nil {
			log.Println("config: unable to bind env: " + err.Error())
		}
	}

	var c Config
	if err := viper.Unmarshal(&c); err != nil {
		return nil, fmt.Errorf("unable to decode into struct, %v", err)
	}

	return &c, nil
}
