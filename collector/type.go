package collector

import (
	"time"

	"github.com/willie-cadete/lxdexplorer-collector/config"
)

type HostNode struct {
	CollectedAt time.Time   `bson:"collectedat"`
	Hostname    string      `bson:"hostname"`
	Containers  interface{} `bson:"containers"`
}

type Network struct {
	IPs     string `bson:"ips"`
	Netmask string `bson:"netmask"`
	CIDR    string `bson:"cidr"`
}

type OS struct {
	Distribution string `bson:"distribution"`
	Version      string `bson:"version"`
}

type Container struct {
	CollectedAt time.Time `bson:"collectedat"`
	Name        string    `bson:"name"`
	Hostnode    string    `bson:"hostnode"`
	Status      string    `bson:"status"`
	Network     Network   `bson:"network"`
	OS          OS        `bson:"os"`
	ImageID     string    `bson:"imageid"`
}

type Collector struct {
	config   *config.Config
	database Database
}

type Options struct {
	Config   *config.Config
	Database Database
}

func NewCollector(opts Options) *Collector {
	return &Collector{
		config:   opts.Config,
		database: opts.Database,
	}
}

type Database interface {
	AddTTL(collection string, field string, seconds int32) error
	InsertMany(collection string, documents []interface{})
}
