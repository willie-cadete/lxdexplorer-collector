package main

import (
	"time"

	log "github.com/sirupsen/logrus"
	"github.com/willie-cadete/lxdexplorer-collector/collector"
	"github.com/willie-cadete/lxdexplorer-collector/config"
	"github.com/willie-cadete/lxdexplorer-collector/database"
)

var version string

func main() {
	// Log the version
	log.Info("LXD Explorer Collector Version: " + version)

	conf, err := config.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	logLevel, err := log.ParseLevel(conf.Logging.Level)
	if err != nil {
		panic(err)
	}

	log.SetLevel(logLevel)
	startCollector(conf)
}

func startCollector(conf *config.Config) {

	database := database.NewDatabase(database.Options{
		Config: conf,
	})

	// Create a new collector
	collect := collector.NewCollector(collector.Options{
		Config:   conf,
		Database: database,
	})

	// Start the collector

	collect.AddLXDTTLs()
	for {
		collect.WorkerCollect()
		time.Sleep(time.Duration(conf.GetCollectorInterval()) * time.Second)
	}

}
