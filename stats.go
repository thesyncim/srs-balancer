package main

import (
	"github.com/gin-gonic/gin"
	"github.com/kr/pretty"
)

var Stats = &stats{}

type stats struct {
}

func (s *stats) nodes(c *gin.Context) {

	pretty.Fprintf(c.Writer, "%+v\n", nodes)

}

func (s *stats) node(c *gin.Context) {

}
