package cloud

import (
	_ "log"
	"sync"
	"time"
)

type provider int

const (
	StartEdgesThreshold float64 = 0.8
	StopEdgesThreshold  float64 = 0.375
	MinServers = 1
	MaxEdgeBW = 200 / 8 * 1024 * 1024 //200Mb

	OvhProvider provider = iota
	DOProvider
)

//Coordinator has the responsability to start and stop nodes based on Edge Stats
type Coordinator struct {
	TotalBW             float64
	StartEdgesThreshold float64
	StopEdgesThreshold  float64
	cluster             *Cluster
	startNodeMu         sync.Mutex
	stopNodeMu          sync.Mutex
	Providers           map[provider]Provider
	activeNodes         map[string]bool
	activeNodesMu       sync.Mutex
}

func NewCoordinator(cl *Cluster) *Coordinator {
	c := &Coordinator{
		StartEdgesThreshold: StartEdgesThreshold,
		StopEdgesThreshold:  StopEdgesThreshold,
		cluster:             cl,
		Providers: map[provider]Provider{
			DOProvider: new(digitalOcean),
		},
		activeNodes: make(map[string]bool),
	}

	go func(c *Coordinator) {
		time.Sleep(15 * time.Second)
		c.activeNodesMu.Lock()
		activenodes := len(c.activeNodes)
		c.activeNodesMu.Unlock()
		if activenodes < 1 {
			c.startNodeMu.Lock()
			err := c.Providers[DOProvider].Authenticate()
			if err != nil {
				c.startNodeMu.Unlock()
				log.Fatalln(err)
			}

			_, err = c.Providers[DOProvider].StartNode(defaultContinent)
			if err != nil {
				c.startNodeMu.Unlock()
				log.Println(err)
			}
			c.startNodeMu.Unlock()
		}

	}(c)
	return c
}

func (c *Coordinator) SetNodeActive(ip string) {
	c.activeNodesMu.Lock()
	c.activeNodes[ip] = true
	c.activeNodesMu.Unlock()
}

func (c *Coordinator) Monitor() {

	go func() {
		for range time.Tick(15 * time.Second) {
			c.startNodeMu.Lock()
			if isoverloaded, continent := c.cluster.IsOverload(); isoverloaded {
				err := c.Providers[DOProvider].Authenticate()
				if err != nil {
					c.startNodeMu.Unlock()
					log.Fatalln(err)
				}

				_, err = c.Providers[DOProvider].StartNode(continent)
				if err != nil {
					c.startNodeMu.Unlock()
					log.Println(err)
				}
				c.startNodeMu.Unlock()

			}

		}
	}()

	for range time.Tick(15 * time.Second) {
		c.stopNodeMu.Lock()
		if isunderloaded, _ := c.cluster.IsUnderload(); isunderloaded {
			err := c.Providers[DOProvider].Authenticate()
			if err != nil {
				//TODO do not die
				c.stopNodeMu.Unlock()
				log.Fatalln(err)
			}

			edgetostop := c.cluster.GetMinLoadEdge()
			err = c.Providers[DOProvider].StopNode(edgetostop)
			if err != nil {
				c.stopNodeMu.Unlock()
				log.Println(err)
			}
			delete(c.activeNodes, edgetostop)
		}

	}

}
