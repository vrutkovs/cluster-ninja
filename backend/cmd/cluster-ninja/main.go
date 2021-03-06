package main

import (
	"log"
	"math/rand"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/vrutkovs/cluster-ninja/internal/api"
)

const attempts = 10

var resourceTypes = []string{
	"pod",
	"statefulset",
	"deployment",
	"daemonset",
}

func main() {
	r := gin.Default()
	r.SetTrustedProxies(nil)
	r.Use(K8S())
	api := r.Group("/api")
	{
		api.GET("/gimme", handleGimme)
		api.POST("/destroy/ns/:namespace/:type/:name", handleKill)
	}

	r.Run()
}

func K8S() gin.HandlerFunc {
	clientSet, err := api.New()
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
	k8s, ok := c.Keys["k8s"].(*api.K8sAPI)
	if !ok {
		log.Panic("Failed to get k8s api")
	}
	var headers gin.H
	for i := 0; i < attempts; i++ {
		// Get random namespace
		namespace, err := k8s.GetRandomNamespace()
		if err != nil {
			log.Panic("Failed to get random namespace")
		}

		// Get random resource type
		resourceType := resourceTypes[rand.Intn(len(resourceTypes))]
		// Get resource names of selected type in selected namespace
		resources, err := k8s.ListResources(namespace, resourceType)
		if err != nil {
			// Restart again
			continue
		}
		name := resources[rand.Intn(len(resources))]
		headers = gin.H{
			"namespace": namespace,
			"type":      resourceType,
			"name":      name,
		}
		break
	}
	if headers == nil {
		c.JSON(http.StatusInternalServerError, gin.H{})
	} else {
		c.JSON(http.StatusOK, headers)
	}
}

func handleKill(c *gin.Context) {
	k8s, ok := c.Keys["k8s"].(*api.K8sAPI)
	if !ok {
		log.Panic("Failed to get k8s api")
	}
	namespace := c.Param("namespace")
	resourceType := c.Param("type")
	name := c.Param("name")
	go k8s.KillResource(namespace, resourceType, name)
}
