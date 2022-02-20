package api

import (
	"context"
	"log"
	"math/rand"
	"time"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/kubectl/pkg/util/slice"
)

var skipNamespaces = []string{
	"openshift-cluster-version", // CVO won't bring itself back
	"openshift-etcd",            // This may kill etcd pod and cause outage
	"openshift-ingress",         // This may remove ingress pods and backend would stop responding
}

func (k *K8sAPI) updateNamespaces() {
	namespaceList, err := k.c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil || namespaceList.Items == nil || len(namespaceList.Items) == 0 {
		log.Panic("failed to fetch namespaces list")
	}

	var nsList []string

	// Get a slice of keys
	for _, ns := range namespaceList.Items {
		if slice.ContainsString(skipNamespaces, ns.Name, nil) {
			continue
		}
		nsList = append(nsList, ns.Name)
	}
	k.namespaces = nsList
}

func (k *K8sAPI) runPeriodicNamespaceUpdate() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		k.updateNamespaces()
	}
}

func (k *K8sAPI) GetRandomNamespace() (string, error) {
	randomNamespace := k.namespaces[rand.Intn(len(k.namespaces))]
	return randomNamespace, nil
}
