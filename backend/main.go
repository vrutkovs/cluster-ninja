package main

import (
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(K8S())
	r.GET("/gimme", handleGimme)
	r.POST("/destroy/ns/:namespace/:type/:name", handleKill)
	r.Run()
}

func K8S() gin.HandlerFunc {
	clientSet, err := NewK8sAPI()
	if err != nil {
		log.Println("Failed to login in cluster")
		log.Panic(err)
	}
	return func(c *gin.Context) {
		c.Set("k8s", clientSet)
		c.Next()
	}
}

func handleGimme(c *gin.Context) {
	_, ok := c.Keys["k8s"].(*K8sAPI)
	if !ok {
		log.Panic("Failed to get k8s api")
	}
}

func handleKill(c *gin.Context) {
	_, ok := c.Keys["k8s"].(*K8sAPI)
	if !ok {
		log.Panic("Failed to get k8s api")
	}
}
