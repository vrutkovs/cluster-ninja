package api

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sAPI) listDaemonsets(namespace string) ([]string, error) {
	var result []string
	daemonSets, err := k.c.AppsV1().DaemonSets(namespace).List(k.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching daemonsets")
	}
	for i := 0; i < len(daemonSets.Items); i++ {
		result = append(result, daemonSets.Items[i].Name)
	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no daemonsets found")
	}
	return result, nil
}

func (k *K8sAPI) killDaemonset(namespace, name string) error {
	if err := k.c.AppsV1().DaemonSets(namespace).Delete(k.ctx, name, k.deleteOpts); err != nil {
		return fmt.Errorf("failed to kill daemonset: %v", err)
	}
	return nil
}
