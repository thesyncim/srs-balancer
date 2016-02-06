package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
)

var Stats = &stats{}

type stats struct {
}

func (s *stats) nodes(c *gin.Context) {
   pretty.Println("IDS: ",nodes.Ids())
   pretty.Println("total: ",len(nodes.Ids()))
 for i:=range nodes.Nodes{

	 pretty.Println(nodes.Nodes[i])

 }
    
 


}

func (s *stats) node(c *gin.Context) {

}
