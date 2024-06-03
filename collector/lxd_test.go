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

const (
	lxd_host string = "localhost"
	lxd_cert string = "./tls/client.crt"
	lxd_key  string = "./tls/client.key"
)

func CreateConfig() *config.Config {
	return &config.Config{
		LXD: config.LXD{
			TLS: config.TLS{
				Cert:   lxd_cert,
				Key:    lxd_key,
				Verify: false,
			},
			Hostnodes: []string{lxd_host},
		},
	}
}

func TestConnect(t *testing.T) {

	tests := []struct {
		name   string
		config *config.Config
		host   string
		asset  func(lxd.InstanceServer)
	}{

		{
			name:   "TestConnectErrorInvalidHost",
			config: CreateConfig(),
			host:   "invalidhost",
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
			name:   "TestAddLXDTTLs",
			config: CreateConfig(),
			asset: func(err error) {
				assert.NoError(t, err, "Expected no error")
			},
			mockFunc: func(t *testing.T, cMock *collector_mocks.Database) {
				cMock.On("AddTTL", mock.Anything, mock.Anything, mock.Anything).Return(nil).Times(2)
			},
		},
		{
			name:   "TestAddLXDTTLsError",
			config: CreateConfig(),
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
