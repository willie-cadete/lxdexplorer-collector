package collector

import (
	"testing"

	lxd "github.com/canonical/lxd/client"
	"github.com/stretchr/testify/assert"
	"github.com/willie-cadete/lxdexplorer-collector/config"
)

func TestConnect(t *testing.T) {

	tests := []struct {
		name   string
		config *config.Config
		host   string
		asset  func(lxd.InstanceServer)
	}{

		{
			name: "TestConnectErrorInvalidHost",
			config: &config.Config{
				Collector: config.Collector{
					Interval:  10,
					Retention: 10,
				},
				Logging: config.Logging{
					Level: "info",
				},
				Database: config.Database{
					URI: "mongodb://localhost:27017",
				},
				LXD: config.LXD{
					TLS: config.TLS{
						Cert:   "./tls/client.crt",
						Key:    "./tls/client.key",
						Verify: false,
					},
					Hostnodes: []string{"https://localhost:8443"},
				},
			},
			host: "invalidhost",
			asset: func(result lxd.InstanceServer) {
				assert.Nil(t, result, "Expected nil, InvalidHost")
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &Collector{
				config: tt.config,
			}
			result := s.Connect(tt.host)

			tt.asset(result)
		})
	}
}
