package main

import (
	"io/ioutil"
	"github.com/ghodss/yaml"
	
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CreatePod(podname string) string {
	cs, err := getInClusterClientSet()
	if err != nil {
		return err.Error()
	}

	var podYaml = getAgentSpecification(podname)
	var p1 v1.Pod
	err1 := yaml.Unmarshal([]byte(podYaml), &p1)
	if err1 != nil {
		return "unmarshal error: " + err1.Error()
	}

	podClient := cs.CoreV1().Pods("azuredevops")
	pod, err2 := podClient.Create(&p1)
	if err2 != nil {
		return "podclient create error: " + err2.Error()
	}

	return pod.GetName()
}

func DeletePod(podname string) string {
	cs, err := getInClusterClientSet()
	if err != nil {
		return err.Error()
	}

	podClient := cs.CoreV1().Pods("azuredevops")

	err2 := podClient.Delete(podname, &metav1.DeleteOptions{})
	if err2 != nil {
		return "podclient delete error: " + err2.Error()
	}

	return "Deleted " + podname
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

func getAgentSpecification(podname string) string {
	// If pod is to be created in a different namespace
	// then secrets need to be created in the same namespace, i.e. VSTS_TOKEN and VSTS_ACCOUNT
	// kubectl create secret generic vsts --from-literal=VSTS_TOKEN=<token> --from-literal=VSTS_ACCOUNT=<accountname>
	dat, err := ioutil.ReadFile("agentpods/" + podname + ".yaml")
	if err != nil {
		return err.Error()
	}

	var podYaml = string(dat)
	return podYaml
}