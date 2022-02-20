package api

import (
	"context"
	"fmt"
	"math/rand"
	"time"

	"github.com/kelseyhightower/envconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type K8sAPI struct {
	c          *k8sclient.Clientset
	ctx        context.Context
	deleteOpts metav1.DeleteOptions
	namespaces []string

	NamespaceOverride string `envconfig:"NAMESPACE"`
}

func New() (*K8sAPI, error) {
	var k8sapi K8sAPI
	// Read env config
	err := envconfig.Process("", &k8sapi)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from env vars")
	}

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
	k8sapi.c = clientSet
	k8sapi.ctx = context.TODO()

	// Set shared DeleteOptions
	var zero int64 = 0
	k8sapi.deleteOpts = metav1.DeleteOptions{
		GracePeriodSeconds: &zero,
	}
	if len(k8sapi.NamespaceOverride) > 0 {
		k8sapi.namespaces = []string{k8sapi.NamespaceOverride}
	} else {
		k8sapi.updateNamespaces()
		go k8sapi.runPeriodicNamespaceUpdate()
	}

	return &k8sapi, nil
}

func (k *K8sAPI) ListResources(namespace, resourceType string) ([]string, error) {
	switch resourceType {
	case "pod":
		return k.listPods(namespace)
	case "statefulset":
		return k.listStatefulsets(namespace)
	case "deployment":
		return k.listDeployments(namespace)
	case "daemonset":
		return k.listDaemonsets(namespace)
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
	case "daemonset":
		k.killDaemonset(namespace, name)
	}
}
