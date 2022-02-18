package api

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"os"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	k8s "k8s.io/client-go/kubernetes"
)

var blackListedNamespaces = []string{
	"openshift-console", // This may remove console deployment
	"openshift-etcd",    // This may kill etcd pod and cause outage
	"openshift-ingress", // This may remove ingress pods and backend would stop responding
	"cluster-ninja",     // Don't kill the app itself
}

var nsList = []string{}

func setNamespacesList(c *k8s.Clientset) error {
	if namespace, ok := os.LookupEnv("NAMESPACE"); ok {
		log.Println(fmt.Sprintf("Namespace override found: %s", namespace))
		nsList = []string{namespace}
		return nil
	}
	log.Println("Fetching available namespaces")
	nms, err := c.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil || nms.Items == nil || len(nms.Items) == 0 {
		return fmt.Errorf("failed to list namespaces: %v", err)
	}

	namespacesMap := map[string]bool{}
	for _, n := range nms.Items {
		namespacesMap[n.Name] = true
	}

	if _, ok := os.LookupEnv("NO_BLACKLIST"); !ok {
		// Remove blacklisted namespaces
		for n := range blackListedNamespaces {
			delete(namespacesMap, blackListedNamespaces[n])
		}
	}

	// Leave CVO alone
	delete(namespacesMap, "openshift-cluster-version")

	// Get a slice of keys
	for _, key := range reflect.ValueOf(namespacesMap).MapKeys() {
		nsList = append(nsList, key.String())
	}
	log.Printf("Namespaces: %v\n", nsList)

	return nil
}

func (k *K8sAPI) GetRandomNamespace() (string, error) {
	if len(nsList) == 0 {
		err := setNamespacesList(k.c)
		if err != nil {
			return "", fmt.Errorf("failed to fetch namespaces list: %v", err)
		}
	}

	randomNamespace := nsList[rand.Intn(len(nsList))]
	log.Println(fmt.Sprintf("random namespace: %v", randomNamespace))
	return randomNamespace, nil
}
