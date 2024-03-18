package main

import (
	"time"

	"github.com/willie-cadete/lxdexplorer-collector/collector"
	"github.com/willie-cadete/lxdexplorer-collector/config"
	"github.com/willie-cadete/lxdexplorer-collector/database"
)

func main() {
	conf, err := config.LoadConfig("/Users/willie/Documents/projects/lxdexplorer-collector/")
	if err != nil {
		panic(err)
	}

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
