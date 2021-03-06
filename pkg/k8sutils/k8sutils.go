package k8sutils

import (
	"fmt"

	"k8s.io/api/core/v1"
	storagev1 "k8s.io/api/storage/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	k8shelper "k8s.io/kubernetes/pkg/api/v1/helper"
)

// GetK8sClient instantiates a k8s client
func GetK8sClient() (*kubernetes.Clientset, error) {
	k8sClient, err := loadClientFromServiceAccount()
	if err != nil {
		return nil, err
	}

	if k8sClient == nil {
		return nil, ErrK8SApiAccountNotSet
	}

	return k8sClient, nil
}

// loadClientFromServiceAccount loads a k8s client from a ServiceAccount
// specified in the pod running
func loadClientFromServiceAccount() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	k8sClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return k8sClient, nil
}

// GetService Gets the service by the name
func GetService(svcName string, svcNS string) (*v1.Service, error) {
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	if svcName == "" {
		return nil, fmt.Errorf("Cannot return service obj without service name")
	}
	svc, err := client.CoreV1().Services(svcNS).Get(svcName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return svc, nil

}

// GetPod Gets the pod by the name
func GetPod(podName string, namespace string) (*v1.Pod, error) {
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	if podName == "" {
		return nil, fmt.Errorf("Cannot return pod obj without pod name")
	}

	pod, err := client.CoreV1().Pods(namespace).Get(podName, metav1.GetOptions{})
	if err != nil {
		return nil, err
	}
	return pod, nil
}

// GetAllPods Get all Pods in the cluster
func GetAllPods() (*v1.PodList, error) {
	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	podList, err := client.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return nil, err
	}
	return podList, nil
}

// GetPVC Gets the PVC by the name
func GetPVC(pvcName string, namespace string) (*v1.PersistentVolumeClaim, error) {
	if pvcName == "" {
		return nil, fmt.Errorf("Empty PVC name")
	}

	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	return client.CoreV1().PersistentVolumeClaims(namespace).Get(pvcName, metav1.GetOptions{})
}

// DeletePod Deletes the pod by the name
func DeletePod(podName string, namespace string, force bool) error {
	client, err := GetK8sClient()
	if err != nil {
		return err
	}

	if podName == "" {
		return fmt.Errorf("Cannot delete pod without pod name")
	}

	deleteOptions := metav1.DeleteOptions{}
	if force {
		gracePeriodSec := int64(0)
		deleteOptions.GracePeriodSeconds = &gracePeriodSec

	}

	return client.CoreV1().Pods(namespace).Delete(podName, &deleteOptions)
}

// GetStorageClassName Gets the storage class name for a PVC
func GetStorageClassName(pvc *v1.PersistentVolumeClaim) string {
	return k8shelper.GetPersistentVolumeClaimClass(pvc)
}

// GetStorageClass Gets the storage class by name
func GetStorageClass(storageClassName string, namespace string) (*storagev1.StorageClass, error) {
	if storageClassName == "" {
		return nil, fmt.Errorf("Empty storage class name")
	}

	client, err := GetK8sClient()
	if err != nil {
		return nil, err
	}

	return client.StorageV1().StorageClasses().Get(storageClassName, metav1.GetOptions{})
}
