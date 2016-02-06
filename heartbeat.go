package main

import (
	"github.com/thesyncim/srs-balancer/cloud"

	"github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
	"fmt"
)

var (
	nodes *cloud.Cluster
	log   *logrus.Logger
)

func init() {
	nodes = cloud.NewCluster()
	log = cloud.Log
}

// Heartbeat recieves posts from srs servers in order to track load and dead hosts
func Heartbeat(c *gin.Context) {

	log.WithField("context", c).Debug("Received HearBeat")

	var Hb cloud.HeartbeatReq

	err := c.BindJSON(&Hb)
	if err != nil {
		log.Println(err)
		return
	}

	nodes.Set(&Hb)
}

var hlstpl = `<smil>
	<head>
		<meta base="http://%s/live" />
	</head> 
	<body>
				<video src="livestream.m3u8" system-bitrate="1000000" width="480" height="360"/>
	</body>
</smil>`
 
var rtmptpl = `<smil>
	<head>
		<meta base="rtmp://%s/live" />
	</head> 
	<body>
				<video src="livestream" system-bitrate="1000000" width="480" height="360"/>
	</body>
</smil>`

func hls(c *gin.Context) {

	ip := nodes.GetEdgeIP(c.ClientIP())

	fmt.Fprint(c.Writer, fmt.Sprintf(hlstpl, ip))



}

func rtmp(c *gin.Context) {

	ip := nodes.GetEdgeIP(c.ClientIP())

	fmt.Fprint(c.Writer, fmt.Sprintf(rtmptpl, ip))

}
