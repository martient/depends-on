package pkg

import (
	"context"
	"fmt"
	"time"

	"github.com/martient/golang-utils/utils"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
)

func WaitForJob(clientset *kubernetes.Clientset, ns string, jobName string, checkInterval int) error {
	utils.LogInfo("Getting '%s' job object...\n", "JOB", jobName)
	for {
		job, err := clientset.BatchV1().Jobs(ns).Get(context.Background(), jobName, metav1.GetOptions{})
		if errors.IsNotFound(err) {
			utils.LogInfo("%s job not found. Retrying in %d seconds...\n", "JOB", jobName, checkInterval)
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		} else if err != nil {
			utils.LogError(fmt.Sprintf("Error getting %s job object, Retrying in %d seconds...\n", jobName, checkInterval), "JOB", err)
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		}

		if job.Status.Active >= 1 {
			utils.LogInfo("%s job is not completed yet. Retrying in %d seconds...\n", "JOB", jobName, checkInterval)
			time.Sleep(time.Duration(checkInterval) * time.Second)
			continue
		} else if job.Status.Succeeded >= 1 {
			utils.LogInfo("%s job succeeded\n", "JOB", jobName)
			return nil
		} else if job.Status.Failed >= 1 {
			utils.LogInfo("%s job failed\n", "JOB", jobName)
			return fmt.Errorf("%s job failed\n", jobName)
		}
	}
}
