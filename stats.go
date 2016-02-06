package main

import (
	"github.com/gin-gonic/gin"

	"fmt"
	"github.com/mongodb/mongo-tools/common/json"
)

var Stats = &stats{}

type stats struct {
}

func (s *stats) nodes(c *gin.Context) {

	c.Writer.WriteString(fmt.Sprintf("IDS: %s \r\n",nodes.Ids()))
	c.Writer.WriteString(fmt.Sprintf("total nodes: %d \r\n",len(nodes.Ids())))

	var totalConn int
 for i:=range nodes.Nodes{

	 c.Writer.WriteString(fmt.Sprintf("Ip %s : %d \r\n",nodes.Nodes[i].IP,nodes.Nodes[i].HeartbeatReq.Summaries.Data.System.ConnSysTw))
	 totalConn+=nodes.Nodes[i].HeartbeatReq.Summaries.Data.System.ConnSysTw

 }


	for i:=range nodes.Nodes{

		json.NewEncoder(c.Writer).Encode(nodes.Nodes[i])
		c.Writer.WriteString("\r\n")

	}

	c.Writer.WriteString(fmt.Sprintf("total Conn: %d \r\n",totalConn))


}

func (s *stats) node(c *gin.Context) {

}
