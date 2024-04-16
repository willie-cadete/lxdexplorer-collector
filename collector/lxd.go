package collector

import (
	"os"
	"time"

	log "github.com/sirupsen/logrus"

	lxd "github.com/canonical/lxd/client"
	"github.com/canonical/lxd/shared/api"
)

func (s *Collector) Connect(h string) lxd.InstanceServer {
	args := s.connectionOptions()

	cnn, err := lxd.ConnectLXD("https://"+h+":8443", args)
	if err != nil {
		log.Println(err)
	}
	return cnn
}

func (s *Collector) ParseContainer(c api.ContainerFull, h string) Container {

	if c.State.Status == "Stopped" {
		return Container{
			Name:     c.Name,
			Hostnode: h,
			Status:   c.State.Status,
			Network: Network{
				IPs:     "N/A",
				Netmask: "N/A",
				CIDR:    "N/A",
			},
			OS: OS{
				Distribution: c.Config["image.os"],
				Version:      c.Config["image.release"],
			},
			ImageID: c.Config["volatile.base_image"][:6],
		}
	}

	return Container{
		Name:     c.Name,
		Hostnode: h,
		Status:   c.State.Status,
		Network: Network{
			IPs:     c.State.Network["eth0"].Addresses[0].Address,
			Netmask: c.State.Network["eth0"].Addresses[0].Netmask,
			CIDR:    c.State.Network["eth0"].Addresses[0].Address + "/" + c.State.Network["eth0"].Addresses[0].Netmask,
		},
		OS: OS{
			Distribution: c.Config["image.os"],
			Version:      c.Config["image.release"],
		},
		ImageID: c.Config["volatile.base_image"][:6],
	}
}

func (s *Collector) AddLXDTTLs() error {
	err := s.database.AddTTL("containers", "collectedat", int32(s.config.GetCollectorInterval()))
	if err != nil {
		return err
	}
	log.Printf("Fetcher: Added TTL to containers collection: %d seconds", s.config.GetCollectorInterval())

	log.Printf("Fetcher: Added TTL to history collection: %d days", s.config.GetCollectorRetention())
	err = s.database.AddTTL("history", "collectedat", int32(s.config.GetCollectorRetention()*60*60*24))
	return err
}

func (s *Collector) WorkerCollect() {

	collectedAt := time.Now().UTC()

	for _, h := range s.getHostnodes() {
		c := s.Connect(h)
		if c == nil {
			continue
		}
		cs, _ := c.GetContainersFull()

		for _, c := range cs {
			container := s.ParseContainer(c, h)
			s.database.InsertMany("containers", []interface{}{Container{CollectedAt: collectedAt, Name: container.Name, Hostnode: container.Hostnode, Status: container.Status, Network: container.Network, OS: container.OS, ImageID: container.ImageID}})
		}
		log.Println("Fetcher: Inserted", len(cs), "containers from", h)

		s.database.InsertMany("history", []interface{}{HostNode{CollectedAt: collectedAt, Hostname: h, Containers: cs}})
		log.Println("Fetcher: Inserted", len(cs), "containers from hostnode:", h)

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

	return s.config.GetHostnodes()
}
