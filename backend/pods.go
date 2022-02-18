package main

import (
	"fmt"

	coreapi "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func (k *K8sAPI) listPods(namespace string) ([]string, error) {
	var result []string
	pods, err := k.c.CoreV1().Pods(namespace).List(k.ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("error fetching pods")
	}
	for i := 0; i < len(pods.Items); i++ {
		var pod = pods.Items[i]

		// Skip failed, succeded and unknown pods
		switch pod.Status.Phase {
		case coreapi.PodFailed, coreapi.PodSucceeded, coreapi.PodUnknown:
			continue
		}

		result = append(result, pod.Name)

	}
	if len(result) == 0 {
		return nil, fmt.Errorf("no pods found")
	}
	return result, nil
}

func (k *K8sAPI) killPod(namespace, name string) error {
	if err := k.c.CoreV1().Pods(namespace).Delete(k.ctx, name, metav1.DeleteOptions{}); err != nil {
		return fmt.Errorf("failed to kill pod: %v", err)
	}
	return nil
}
