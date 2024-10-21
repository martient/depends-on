package pkg

import (
	"context"
	"time"

	"github.com/martient/golang-utils/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"
)

func WaitForService(clientset *kubernetes.Clientset, ns string, serviceName string, checkInterval int) error {
	utils.LogInfo("Getting '%s' service object...\n", "SERVICE", serviceName)
	for {
		service, err := clientset.CoreV1().Services(ns).Get(context.Background(), serviceName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			utils.LogInfo("%s service not found. Retrying in %d seconds...\n", "SERVICE", serviceName, checkInterval)
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		} else if err != nil {
			utils.LogInfo("Error getting %s service object: %v Retrying in %d seconds...\n", "SERVICE", serviceName, err.Error(), checkInterval)
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		}

		set := labels.Set(service.Spec.Selector)

		for {
			utils.LogInfo("Getting pods for the '%s' service...\n", "SERVICE", serviceName)
			pods, err := clientset.CoreV1().Pods(ns).List(context.Background(), metav1.ListOptions{LabelSelector: set.AsSelector().String()})
			if err != nil {
				utils.LogError("Error getting the pods", "SERVICE", err)
				return err
			}

			if len(pods.Items) < 1 {
				utils.LogInfo("No pods found for the '%s' service. Retrying...\n", "SERVICE", serviceName)
				time.Sleep(1 * time.Second)
				continue
			}

			utils.LogInfo("Checking readiness of the '%s' service pods...\n", "SERVICE", serviceName)

			readyPodFound := false

			for _, pod := range pods.Items {
				for _, cond := range pod.Status.Conditions {
					if cond.Type == "Ready" && cond.Status == "True" {
						utils.LogInfo("%s is ready.\n", "SERVICE", pod.GetName())
						readyPodFound = true
						break
					}
				}
				if readyPodFound {
					break
				}
				utils.LogInfo("%s is not ready yet. Retrying in %d seconds...\n", "SERVICE", pod.GetName(), checkInterval)
			}
			if readyPodFound {
				return nil
			}
			time.Sleep(time.Duration(checkInterval) * time.Second)
		}
	}
}
