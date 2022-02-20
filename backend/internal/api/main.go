package api

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"

	"github.com/kelseyhightower/envconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8sclient "k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

	// Seed random
	rand.Seed(time.Now().Unix())

	// Read env config
	err := envconfig.Process("", &k8sapi)
	if err != nil {
		return nil, fmt.Errorf("failed to read config from env vars")
	}

	// Read kubeconfig from env var or in-cluster
	var restconfig *rest.Config
	// Try to read KUBECONFIG env var
	kubeconfig := os.Getenv("KUBECONFIG")
	if kubeconfig == "" {
		kubeconfig = filepath.Join(
			os.Getenv("HOME"), ".kube", "config",
		)
	}
	restconfig, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		// Try to read config from in-cluster configuration
		configLoader := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(
			clientcmd.NewDefaultClientConfigLoadingRules(),
			&clientcmd.ConfigOverrides{},
		)
		restconfig, err = configLoader.ClientConfig()
		if err != nil {
			return nil, err
		}
	}

	// Create the clientset
	clientSet, err := k8sclient.NewForConfig(restconfig)
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

	// Initialize namespaces list
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
