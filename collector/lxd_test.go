package collector

import (
	"errors"
	"testing"

	lxd "github.com/canonical/lxd/client"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	collector_mocks "github.com/willie-cadete/lxdexplorer-collector/collector/mocks"
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

func TestAddLXDTTLs(t *testing.T) {

	tests := []struct {
		name     string
		config   *config.Config
		mockFunc func(t *testing.T, cMock *collector_mocks.Database)
		asset    func(err error)
	}{
		{
			name: "TestAddLXDTTLs",
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
			asset: func(err error) {
				assert.NoError(t, err, "Expected no error")
			},
			mockFunc: func(t *testing.T, cMock *collector_mocks.Database) {
				cMock.On("AddTTL", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
			},
		},
		{
			name: "TestAddLXDTTLsError",
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
			asset: func(err error) {
				assert.Error(t, err, "Expected error")
			},
			mockFunc: func(t *testing.T, cMock *collector_mocks.Database) {
				cMock.On("AddTTL", mock.Anything, mock.Anything, mock.Anything).Return(errors.New("error to add index")).Times(1)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dbMock := collector_mocks.NewDatabase(t)
			if tt.mockFunc != nil {
				tt.mockFunc(t, dbMock)
			}

			s := NewCollector(Options{
				Config:   tt.config,
				Database: dbMock,
			})

			err := s.AddLXDTTLs()

			tt.asset(err)
		})
	}
}
