package main

import (
	"os"
	"path/filepath"

	v1alpha1 "github.com/microsoft/poolprovider-for-k8s/pkg/apis/dev/v1alpha1"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/kubernetes/fake"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
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

	if v1alpha1.IsTestingEnv() {
		if client.clientset == nil {
			client.clientset = fake.NewSimpleClientset()
		}
	} else {
		cs, _ := GetClientSet()
		client.clientset = cs
	}
	return &client
}

func SetTestingEnvironmentVariables(params ...bool) {
	os.Setenv("IS_TESTENVIRONMENT", "true")
	os.Setenv("VSTS_SECRET", "sharedsecret1234")
	client.clientset = fake.NewSimpleClientset()
	if len(params) == 0 {
		params = append(params, false)
	}
	CreateDummyPod(params[0])
}

func CreateDummyPod(isbuildkit bool) {
	cs := CreateClientSet()

	pod := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "azurepipelinesagentpod",
			Namespace: "azuredevops",
		},
		Spec: v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "agentimage",
					Image: "prebansa/myagent:v5.16",
				},
			},
		},
	}
	pod.SetLabels(map[string]string{
		"app": "azurepipelinespool-operator",
	})

	if isbuildkit {
		pod.SetLabels(map[string]string{
			"role": "buildkit",
		})
		pod.ObjectMeta.Name = "buildkitd-0"
	}

	podClient := cs.clientset.CoreV1().Pods("azuredevops")
	_, _ = podClient.Create(pod)
}
