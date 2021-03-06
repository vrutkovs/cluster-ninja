package api

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sAPI) listDeployments(namespace string) ([]string, error) {
	var result []string
	deployments, err := k.c.AppsV1().Deployments(namespace).List(k.ctx, metav1.ListOptions{})
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

func (k *K8sAPI) killDeployment(namespace, name string) error {
	// Don't kill frontend/backend deployments, if apiserver is down it won't be able to bring pods back
	if name == "frontend" || name == "backend" {
		return nil
	}

	if err := k.c.AppsV1().Deployments(namespace).Delete(k.ctx, name, k.deleteOpts); err != nil {
		return fmt.Errorf("failed to kill deployment: %v", err)
	}
	return nil
}
