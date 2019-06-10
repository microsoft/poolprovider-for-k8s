package main

import (
	"strings"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func GetPodNames() string {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return err.Error()
	}

	// creates the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err.Error()
	}

	pods, err := clientset.CoreV1().Pods("").List(metav1.ListOptions{})
	if err != nil {
		return err.Error()
	}

	var ret strings.Builder
	ret.WriteString("Name of one pod is: ")
	ret.WriteString(pods.Items[0].GetName())

	return ret.String();
}