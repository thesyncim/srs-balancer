package main

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/itsjamie/gin-cors"
)

func main() {
	router := gin.Default()


	router.Use(cors.Middleware(cors.Config{
		Origins:         "*",
		Methods:         "GET, PUT, POST, DELETE",
		RequestHeaders:  "Origin, Authorization, Content-Type",
		ExposedHeaders:  "",
		MaxAge:          50 * time.Second,
		Credentials:     true,
		ValidateHeaders: false,
	}))

	router.GET("/crossdomain.xml", func(c *gin.Context) {
		c.Writer.Write([]byte(`<cross-domain-policy>
  <allow-access-from domain="*"/>
</cross-domain-policy>`))
	})
    
 

	heartbeat := router.Group("/heartbeat")
	heartbeat.POST("ping", Heartbeat)
 
	balancer := router.Group("/balancer")
	balancer.GET("hls.smil", hls)
	balancer.GET("rtmp.smil", rtmp)

	stats := router.Group("/stats") 
	stats.GET("/nodes", Stats.nodes)
	stats.GET("/nodes/:id", nil)

	router.Run(":80")

}
