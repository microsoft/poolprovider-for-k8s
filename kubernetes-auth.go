package main 

import (
	"path/filepath"
	"os"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/kubernetes/fake"
)

var client k8s

// Gets the application client set. This can be used to initialize the various Kubernetes clients.
// Uses in cliuster configuration when app is running inside the cluster, or kubeconfig file from
// home directory when running in development mode.
func GetClientSet() (*kubernetes.Clientset, error) {
	debugMode := os.Getenv("DEBUG_LOCAL")
	if debugMode != "" {
		return getOutOfClusterClientSet()
	} else {
		return getInClusterClientSet()
	}
}

func getInClusterClientSet() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func getOutOfClusterClientSet() (*kubernetes.Clientset, error) {
	kubeconfigPath := filepath.Join(homeDir(), ".kube", "config")

	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)
	if err != nil {
		return nil, err
	}

	// create the clientset
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func homeDir() string {
	if h := os.Getenv("HOME"); h != "" {
		return h
	}

	// Windows
	return os.Getenv("USERPROFILE")
}

func CreateClientSet() *k8s {
	
	//client = k8s{}
	testingMode := os.Getenv("COUNTTEST")
	
	if testingMode == "1" || testingMode == "2" {
		if testingMode == "1" {
			client.clientset = fake.NewSimpleClientset()
			os.Setenv("COUNTTEST","2")
		}
	} else {
	  	cs, _ := GetClientSet()
	  	client.clientset = cs
	}
	return &client
}

func isTestingEnv() bool {
	testingMode := os.Getenv("COUNTTEST")
	
	if testingMode == "1" || testingMode == "2" {
		return true
	}
	return false
}