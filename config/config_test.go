package config_test

import (
	"os"
	"testing"

	"github.com/willie-cadete/lxdexplorer-collector/config"

	"github.com/stretchr/testify/assert"
)

const (
	lxd_cert = "./tls/client.crt"
	lxd_key  = "./tls/client.key"
)

func TestLoadConfigWithDefaults(t *testing.T) {
	cfg, err := config.LoadConfig("NotExistingPath")

	expected := &config.Config{
		Collector: config.Collector{
			Interval:  60,
			Retention: 60,
		},
		Logging: config.Logging{
			Level: "info",
		},
		Database: config.Database{
			URI: "mongodb://localhost:27017",
		},
		LXD: config.LXD{
			TLS: config.TLS{
				Cert:   lxd_cert,
				Key:    lxd_key,
				Verify: false,
			},
			Hostnodes: []string{"https://localhost:8443"},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, cfg)
}

func TestLoadConfigWithCustomValues(t *testing.T) {
	cfg, err := config.LoadConfig("testdata")

	expected := &config.Config{
		Collector: config.Collector{
			Interval:  10,
			Retention: 10,
		},
		Logging: config.Logging{
			Level: "warn",
		},
		Database: config.Database{
			URI: "mongodb://localhost:27016",
		},
		LXD: config.LXD{
			TLS: config.TLS{
				Cert:   lxd_cert,
				Key:    lxd_key,
				Verify: true,
			},
			Hostnodes: []string{"127.0.0.2"},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, cfg)
}

func TestGetHostNodes(t *testing.T) {
	cfg, _ := config.LoadConfig("testdata")

	expected := []string{"127.0.0.2"}

	assert.Equal(t, expected, cfg.GetHostnodes())
}

func TestGetTLS(t *testing.T) {
	cfg, _ := config.LoadConfig("testdata")

	expected := config.TLS{
		Cert:   lxd_cert,
		Key:    lxd_key,
		Verify: true,
	}

	assert.Equal(t, expected, cfg.GetTLS())
}

func TestGetDatabaseURI(t *testing.T) {
	cfg, _ := config.LoadConfig("testdata")

	expected := "mongodb://localhost:27016"

	assert.Equal(t, expected, cfg.GetDatabaseURI())
}

func TestGetCollectorInterval(t *testing.T) {
	cfg, _ := config.LoadConfig("testdata")

	expected := 10

	assert.Equal(t, expected, cfg.GetCollectorInterval())
}

func TestGetCollectorRetention(t *testing.T) {
	cfg, _ := config.LoadConfig("testdata")

	expected := 10

	assert.Equal(t, expected, cfg.GetCollectorRetention())
}

func TestGetLoggingLevel(t *testing.T) {
	cfg, _ := config.LoadConfig("testdata")

	expected := "warn"

	assert.Equal(t, expected, cfg.GetLoggingLevel())
}

func TestLoadConfigWithEnvVariablesOverride(t *testing.T) {
	// Set environment variables
	_ = os.Setenv("COLLECTOR_INTERVAL", "5")
	_ = os.Setenv("COLLECTOR_RETENTION", "5")
	_ = os.Setenv("LOGGING_LEVEL", "error")
	_ = os.Setenv("DATABASE_URI", "mongodb://localhost:27015")
	_ = os.Setenv("LXD_TLS_CERTIFICATE", "./tls/client2.crt")
	_ = os.Setenv("LXD_TLS_KEY", "./tls/client2.key")
	_ = os.Setenv("LXD_TLS_VERIFY", "false")
	_ = os.Setenv("LXD_HOSTNODES", "https://localhost:8443,https://localhost:8444")

	cfg, err := config.LoadConfig("testdata")

	expected := &config.Config{
		Collector: config.Collector{
			Interval:  5,
			Retention: 5,
		},
		Logging: config.Logging{
			Level: "error",
		},
		Database: config.Database{
			URI: "mongodb://localhost:27015",
		},
		LXD: config.LXD{
			TLS: config.TLS{
				Cert:   "./tls/client2.crt",
				Key:    "./tls/client2.key",
				Verify: false,
			},
			Hostnodes: []string{"https://localhost:8443", "https://localhost:8444"},
		},
	}

	assert.Nil(t, err)
	assert.Equal(t, expected, cfg)
}
