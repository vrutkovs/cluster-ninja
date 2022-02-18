package main

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sAPI) listDeployments(namespace string) ([]string, error) {
	var result []string
	ctx := context.TODO()
	deployments, err := k.c.AppsV1().Deployments(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching deployments")
	}
	for i := 0; i < len(deployments.Items); i++ {
		result = append(result, deployments.Items[i].Name)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no deployments found")
	}
	return result, nil
}
