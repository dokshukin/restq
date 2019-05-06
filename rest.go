package main

import (
	"os"

	"github.com/gin-gonic/gin"
)

func restAPI() {
	bindIP := "0.0.0.0"
	if len(os.Getenv("RESTQ_BIND_IP")) > 0 {
		bindIP = os.Getenv("RESTQ_BIND_IP")
	}

	bindPort := "8080"
	if len(os.Getenv("RESTQ_BIND_PORT")) > 0 {
		bindPort = os.Getenv("RESTQ_BIND_PORT")
	}

	prefixURI := "/"
	if len(os.Getenv("RESTQ_PREFIX_URI")) > 0 {
		prefixURI = os.Getenv("RESTQ_PREFIX_URI")
	}

	gin.SetMode(gin.ReleaseMode)
	r := gin.Default()
	r.POST(prefixURI+"/:queue", postQueueHandler)
	r.Run(bindIP + ":" + bindPort)
}
