package main

import (
	"time"

	"github.com/willie-cadete/lxdexplorer-collector/collector"
	"github.com/willie-cadete/lxdexplorer-collector/config"
)

func main() {
	conf, err := config.LoadConfig("/Users/willie/Documents/projects/lxdexplorer-collector/config.yaml")
	if err != nil {
		panic(err)
	}

	startCollector(conf)
}

func startCollector(conf *config.Config) {
	// Create a new collector
	collect := collector.NewCollector(collector.Options{
		Config: conf,
	})

	collect.AddLXDTTLs()
	for {
		collect.WorkerCollect()
		time.Sleep(time.Duration(conf.GetCollectorInterval()) * time.Second)
	}

}
