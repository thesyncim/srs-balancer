package main

import (
	"github.com/gin-gonic/gin"

	"fmt"
)

var Stats = &stats{}

type stats struct {
}

func (s *stats) nodes(c *gin.Context) {

	c.Writer.Write([]byte(fmt.Sprintf("IDS: %s",nodes.Ids())))
	c.Writer.Write([]byte(fmt.Sprintf("total: %d ",len(nodes.Ids()))))
 for i:=range nodes.Nodes{
	 c.Writer.Write([]byte(fmt.Sprintf("%+v",nodes.Nodes[i])))

 }
    
 


}

func (s *stats) node(c *gin.Context) {

}
