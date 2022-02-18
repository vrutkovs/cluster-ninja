package main

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sAPI) listStatefulsets(namespace string) ([]string, error) {
	var result []string
	statefulSets, err := k.c.AppsV1().StatefulSets(namespace).List(k.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching statefulsets")
	}
	for i := 0; i < len(statefulSets.Items); i++ {
		result = append(result, statefulSets.Items[i].Name)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no statefulsets found")
	}
	return result, nil
}
