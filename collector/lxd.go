package collector

import (
	"os"
	"strings"
	"time"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
	log "github.com/sirupsen/logrus"
	"github.com/willie-cadete/lxdexplorer-collector/utils"
)

func (s *Collector) Connect(h string) lxd.InstanceServer {
	args := s.connectionOptions()

	cnn, err := lxd.ConnectLXD("https://"+h+":8443", args)
	if err != nil {
		log.Warnln(err)
	}
	return cnn
}

func (s *Collector) ParseContainer(c api.ContainerFull, h string) Container {

	container := Container{
		Name:     strings.ToLower(c.Name),
		Hostnode: strings.ToLower(h),
		Status:   utils.Capitalize(c.State.Status),
		OS: OS{
			Distribution: utils.Capitalize(c.Config["image.os"]),
			Version:      utils.Capitalize(c.Config["image.release"]),
		},
		ImageID: c.Config["volatile.base_image"][:6],
	}

	if utils.Capitalize(c.State.Status) == "Stopped" {
		container.Network = Network{
			IPs:     "N/A",
			Netmask: "N/A",
			CIDR:    "N/A",
		}

		return container
	}

	container.Network = Network{
		IPs:     c.State.Network["eth0"].Addresses[0].Address,
		Netmask: c.State.Network["eth0"].Addresses[0].Netmask,
		CIDR:    c.State.Network["eth0"].Addresses[0].Address + "/" + c.State.Network["eth0"].Addresses[0].Netmask,
	}

	return container

}

func (s *Collector) AddLXDTTLs() error {
	err := s.database.AddTTL("containers", "collectedat", int32(s.config.GetCollectorInterval()))
	if err != nil {
		return err
	}
	log.Infof("Collector: Setting collection interval of %d seconds", s.config.GetCollectorInterval())

	log.Infof("Collector: Setting collection history of %d days", s.config.GetCollectorRetention())
	err = s.database.AddTTL("history", "collectedat", int32(s.config.GetCollectorRetention()*60*60*24))
	return err
}

func (s *Collector) WorkerCollect() {
	start := time.Now()

	defer func() {
		log.Infoln("Collector: Collection took: ", time.Since(start))
	}()

	log.Infoln("Collector: Starting collection from", len(s.getHostnodes()), "hostnodes...")

	collectedAt := time.Now().UTC()

	for _, h := range s.getHostnodes() {
		c := s.Connect(h)
		if c == nil {
			continue
		}

		startContainer := time.Now()
		cs, _ := c.GetContainersFull()

		for _, c := range cs {

			log.Debugln("Collector: ", c.Name, h)
			log.Debugln("Collector: Parsing container", c.Name, "has status", c.State, c.Config, "from", h)
			container := s.ParseContainer(c, h)
			log.Debugln("Collector: Parsed container", container.Name, "has status", container.Status, container.Network, container.OS, container.ImageID, "from", h)
			s.database.InsertMany("containers", []interface{}{Container{CollectedAt: collectedAt, Name: container.Name, Hostnode: container.Hostnode, Status: container.Status, Network: container.Network, OS: container.OS, ImageID: container.ImageID}})
		}

		s.database.InsertMany("history", []interface{}{HostNode{CollectedAt: collectedAt, Hostname: h, Containers: cs}})
		log.Infoln("Collector: Collected", len(cs), "containers from hostnode:", h, "in", time.Since(startContainer))

	}

}

func (s *Collector) connectionOptions() *lxd.ConnectionArgs {

	TLSCertificate, _ := os.ReadFile(s.config.GetLXDTLSCertificate())
	TLSKey, _ := os.ReadFile(s.config.GetLXDTLSKey())

	args := lxd.ConnectionArgs{
		TLSClientCert:      string(TLSCertificate),
		TLSClientKey:       string(TLSKey),
		InsecureSkipVerify: !s.config.GetLXDTLSVerify(),
		SkipGetServer:      false,
	}

	return &args

}

func (s *Collector) getHostnodes() []string {

	hostnodes := s.config.GetHostnodes()
	for i, node := range hostnodes {
		hostnodes[i] = strings.ToLower(node)
	}
	return hostnodes
}
