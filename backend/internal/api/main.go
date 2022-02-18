package api

import (
	"context"
	"math/rand"
	"time"

	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sAPI struct {
	c   *k8sclient.Clientset
	ctx context.Context
}

func New() (*K8sAPI, error) {
	// Create the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// Seed random
	rand.Seed(time.Now().Unix())

	// Create the clientset
	clientSet, err := k8sclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	ctx := context.TODO()
	return &K8sAPI{
		c:   clientSet,
		ctx: ctx,
	}, nil
}

func (k *K8sAPI) ListResources(namespace, resourceType string) ([]string, error) {
	switch resourceType {
	case "pod":
		return k.listPods(namespace)
	case "statefulset":
		return k.listStatefulsets(namespace)
	case "deployment":
		return k.listDeployments(namespace)
	default:
		return nil, nil
	}
}

func (k *K8sAPI) KillResource(namespace, resourceType, name string) {
	switch resourceType {
	case "pod":
		k.killPod(namespace, name)
	case "statefulset":
		k.killStatefulset(namespace, name)
	case "deployment":
		k.killDeployment(namespace, name)
	}
}
