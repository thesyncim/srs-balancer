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

	heartbeat := router.Group("/heartbeat")
	heartbeat.POST("ping", Heartbeat)

	balancer := router.Group("/balancer")
	balancer.GET("hls", hls)
	balancer.GET("rtmp", rtmp)

	router.Run(":80")

}
