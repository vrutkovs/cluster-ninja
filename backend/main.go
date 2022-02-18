package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
)

const attempts = 10

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
	k8s, ok := c.Keys["k8s"].(*K8sAPI)
	if !ok {
		log.Panic("Failed to get k8s api")
	}
	// Get random namespace
	namespace, err := k8s.getRandomNamespace()
	if err != nil {
		log.Panic("Failed to get random namespace")
	}

	var resourceType, name string
	var resources []string
	for i := 0; i < attempts; i++ {
		// Get random resource type
		resourceType = resourceTypes[rand.Intn(len(resourceTypes))]
		// Get resource names of selected type in selected namespace
		resources, err = k8s.listResources(namespace, resourceType)
		if err != nil {
			// Restart again
			continue
		}
		name = resources[rand.Intn(len(resources))]
	}

	c.JSON(http.StatusOK, gin.H{
		"namespace": namespace,
		"type":      resourceType,
		"name":      name,
	})
}

func handleKill(c *gin.Context) {
	_, ok := c.Keys["k8s"].(*K8sAPI)
	if !ok {
		log.Panic("Failed to get k8s api")
	}
}
