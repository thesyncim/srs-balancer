package cloud

import (
	"math"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"sort"
)

type EdgeInfo struct {
	*HeartbeatReq
	Coordinates Coordinates
	LastHBeat   time.Time
	CurrentBw   int64
	Continent   Continent
}

type Stats struct {
	c *Cluster
}

func (s *Stats) Nodes() []*EdgeInfo {

	s.c.mu.Lock()
	defer s.c.mu.Unlock()
	var nodes []*EdgeInfo
	for i := range s.c.Nodes {
		nodes = append(nodes, s.c.Nodes[i])
	}

	return nodes
}
func (s *Stats) Node(id string) *EdgeInfo {
	s.c.mu.Lock()
	defer s.c.mu.Unlock()

	if val, ok := s.c.Nodes[id]; ok {
		return val
	}
	return nil
}

type Cluster struct {
	Nodes         map[string]*EdgeInfo
	Coordinator   *Coordinator
	mu            sync.RWMutex
	underloadEgde string
}

//NewCluster returns a servers instance
func NewCluster() *Cluster {
	log.Debug("Started")
	s := &Cluster{}
	s.Coordinator = NewCoordinator(s)
	s.Nodes = make(map[string]*EdgeInfo)
	go s.RemoveDead()
	return s
}

func (s *Cluster) Ids() []string {
	var ids []string
	s.mu.Lock()
	for ip := range s.Nodes {
		ids = append(ids, ip)
	}
	s.mu.Unlock()
	return ids
}

func (s *Cluster) GetMinLoadEdge() string {

	s.mu.Lock()
	defer s.mu.Unlock()
	return s.underloadEgde
}

//IsUnderload  return the IsUnderload Continent if any(doesnt return Continents < 2 active nodes)
func (s *Cluster) IsUnderload() (bool, Continent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var loadbyCo = map[Continent]struct {
		load            float64
		totalNodes      float64
		MinLoadedEdgeIp string
		MinLoadedBW     int64
	}{}

	var visitedContinent = map[Continent]bool{}

	//get load by continent
	for i := range s.Nodes {

		tmp := loadbyCo[s.Nodes[i].Continent]

		//
		if _, ok := visitedContinent[s.Nodes[i].Continent]; !ok {
			visitedContinent[s.Nodes[i].Continent] = true

			tmp.MinLoadedBW = math.MaxInt64
		}

		tmp.load += float64(s.Nodes[i].CurrentBw)
		tmp.totalNodes++
		if s.Nodes[i].CurrentBw < tmp.MinLoadedBW {
			tmp.MinLoadedBW = s.Nodes[i].CurrentBw
			tmp.MinLoadedEdgeIp = s.Nodes[i].IP

		}
	}

	//check if Underload

	for i := range loadbyCo {
		//calculate the medium load/server
		loadbytes := loadbyCo[i].load / loadbyCo[i].totalNodes
		load := loadbytes / MaxEdgeBW * StopEdgesThreshold

		if load <= 1 && loadbyCo[i].totalNodes > 1 {
			s.underloadEgde = loadbyCo[i].MinLoadedEdgeIp
			return true, i
		}
	}

	return false, Undefined
}

//IsOverload check return the Overloaded Continent if any
func (s *Cluster) IsOverload() (bool, Continent) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	var loadbyCo map[Continent]*struct {
		load       float64
		totalNodes float64
	}
	//get load by continent
	for i := range s.Nodes {
		loadbyCo[s.Nodes[i].Continent].load += float64(s.Nodes[i].CurrentBw)
		loadbyCo[s.Nodes[i].Continent].totalNodes++
	}

	//check if overloaded

	for i := range loadbyCo {
		loadbytes := loadbyCo[i].load / loadbyCo[i].totalNodes
		load := loadbytes / MaxEdgeBW * StartEdgesThreshold
		if load >= 1 {
			return true, i
		}
	}
	return false, Undefined
}

type server struct {
	distance float64
	ip       string
	load     int64
}
type servers []server

// Len is part of sort.Interface.
func (s servers) Len() int {
	return len(s)
}

// Swap is part of sort.Interface.
func (s servers) Swap(i, j int) {
	s[i], s[j] = s[j], s[i]
}

// Less is part of sort.Interface. It is implemented by calling the "by" closure in the sorter.
func (s servers) Less(i, j int) bool {
	return s[i].distance < s[j].distance
}

func (s *Cluster) GetEdgeIP(userIP string) string {
	//holds map-> distance to ip
	var distancebyIPmap = make(map[float64]string)

	//get coords based on ip
	ipLoc, _ := IPToCoords(userIP)

	//calculate user distance to each server
	var distance float64
	s.mu.RLock()

	var Servers servers

	for ip := range s.Nodes {
		coord := s.Nodes[ip].Coordinates
		distance = Distance(coord.Latitude, ipLoc.Longitude, ipLoc.Latitude, coord.Longitude)

		//we are excluding overloaded servers
		if float64(s.Nodes[ip].CurrentBw) > MaxEdgeBW*StartEdgesThreshold {
			continue
		}
		Servers = append(Servers, server{
			distance: distance,
			ip:       ip,
			load:     s.Nodes[ip].CurrentBw,
		})
		//fill distance -> ip
		distancebyIPmap[distance] = ip
	}
	s.mu.RUnlock()

	log.WithFields(logrus.Fields{
		"userIP":     userIP,
		"IPLocation": ipLoc,
		"DistanceKM": distance,
	}).Debug("GetNearIp")

	//TODO pick a random or lowest  (could not find a free one)
	if Servers.Len() == 0 {
		//maps iteration are random
		for ip := range s.Nodes {
			return ip
		}

	}

	//Servers sorted by distance
	sort.Sort(Servers)

	//pick 30% off the  closest servers and chose the one less overloaded

	number := math.Floor((float64(Servers.Len()) / 3) + .5)
	if number == 0 {
		number = 1
	}

	load := int64(math.MaxInt64)
	index := 0

	for i := 0; i < int(number); i++ {
		if Servers[i].load < int64(load) {
			sload := Servers[i].load
			load = int64(sload)
			index = i
		}
	}

	return Servers[index].ip

}

//TODO validate hb request
//Set create or update an Edge server
func (s *Cluster) Set(hb *HeartbeatReq) {
	s.mu.Lock()
	coords, co := IPToCoords(hb.IP)
	ei := &EdgeInfo{
		hb,
		coords,
		time.Now(),
		int64(hb.Summaries.Data.System.ConnSrs),
		co,
	}

	s.Coordinator.SetNodeActive(hb.IP)
	s.Nodes[hb.IP] = ei

	s.Nodes[hb.IP].CurrentBw = int64(hb.Summaries.Data.System.ConnSrs)

	s.mu.Unlock()
	log.WithField("hertBeatReq", hb).Debug("received HB request")
}

func (s *Cluster) LoadByContinent(co Continent) float64 {
	s.mu.Lock()

	var load int64
	var ctr int
	for i := range s.Nodes {
		if s.Nodes[i].Continent == co {
			ctr++
			load += s.Nodes[i].CurrentBw
		}
	}

	return float64(load) / float64(ctr) / float64(MaxEdgeBW*ctr)
}

func (s *Cluster) remove(key string) {
	delete(s.Nodes, key)
}

func (s *Cluster) RemoveDead() {
	for range time.Tick(5 * time.Second) {
		s.mu.Lock()
		for key := range s.Nodes {
			if !s.isAlive(key) || s.Nodes[key].Summaries.Code != 0 {
				log.WithField("node", s.Nodes[key]).Debug("Deleted Node")
				s.remove(key)
			}
		}
		s.mu.Unlock()
	}
}

func (s *Cluster) isAlive(edgeIP string) bool {

	lastbeat := s.Nodes[edgeIP].LastHBeat
	lastbeatBudget := lastbeat.Add(15 * time.Second)

	if lastbeatBudget.After(time.Now()) {
		return true
	}
	return false
}

//HeartbeatReq ...
type HeartbeatReq struct {
	DeviceID  string `json:"Device_id"`
	IP        string
	Summaries *struct {
		Code int `json:"code"`
		Data *struct {
			NowMs int  `json:"now_ms"`
			Ok    bool `json:"ok"`
			Self  *struct {
				Argv       string  `json:"argv"`
				CPUPercent float64 `json:"cpu_percent"`
				Cwd        string  `json:"cwd"`
				MemKbyte   int     `json:"mem_kbyte"`
				MemPercent float64 `json:"mem_percent"`
				Pid        int     `json:"pid"`
				Ppid       int     `json:"ppid"`
				SrsUptime  float64 `json:"srs_uptime"`
				Version    string  `json:"version"`
			} `json:"self"`
			System *struct {
				ConnSrs         int     `json:"conn_srs"`
				ConnSys         int     `json:"conn_sys"`
				ConnSysEt       int     `json:"conn_sys_et"`
				ConnSysTw       int     `json:"conn_sys_tw"`
				ConnSysUDP      int     `json:"conn_sys_udp"`
				CPUPercent      float64 `json:"cpu_percent"`
				Cpus            int     `json:"cpus"`
				CpusOnline      int     `json:"cpus_online"`
				DiskBusyPercent float64 `json:"disk_busy_percent"`
				DiskReadKBps    int     `json:"disk_read_KBps"`
				DiskWriteKBps   int     `json:"disk_write_KBps"`
				IldeTime        float64 `json:"ilde_time"`
				Load15m         float64 `json:"load_15m"`
				Load1m          float64 `json:"load_1m"`
				Load5m          float64 `json:"load_5m"`
				MemRAMKbyte     int     `json:"mem_ram_kbyte"`
				MemRAMPercent   float64 `json:"mem_ram_percent"`
				MemSwapKbyte    int     `json:"mem_swap_kbyte"`
				MemSwapPercent  float64 `json:"mem_swap_percent"`
				NetRecvBytes    int     `json:"net_recv_bytes"`
				NetRecviBytes   int     `json:"net_recvi_bytes"`
				NetSampleTime   int     `json:"net_sample_time"`
				NetSendBytes    int     `json:"net_send_bytes"`
				NetSendiBytes   int     `json:"net_sendi_bytes"`
				SrsRecvBytes    int     `json:"srs_recv_bytes"`
				SrsSampleTime   int     `json:"srs_sample_time"`
				SrsSendBytes    int     `json:"srs_send_bytes"`
				Uptime          float64 `json:"uptime"`
			} `json:"system"`
		} `json:"data"`
	}
}
