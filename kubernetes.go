package main

import (
	"io/ioutil"
	"github.com/ghodss/yaml"
	
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CreatePod(podname, agentId string) PodResponse {
	cs, err := getInClusterClientSet()
	var response PodResponse
	if err != nil {
		response.Status = "failure"
		response.Message = err.Error()
		return response
	}

	var podYaml = getAgentSpecification(podname)
	var p1 v1.Pod
	err1 := yaml.Unmarshal([]byte(podYaml), &p1)
	if err1 != nil {
		response.Status = "failure"
		response.Message = "unmarshal error: " + err1.Error()
		return response
	}

	// Set the agentId as label if specified
	if(agentId != "") {
		p1.Labels["AgentId"] = agentId
	}

	podClient := cs.CoreV1().Pods("azuredevops")
	pod, err2 := podClient.Create(&p1)
	if err2 != nil {
		response.Status = "failure"
		response.Message = "podclient create error: " + err2.Error()
		return response
	}

	response.Status = "success"
	response.Message = "Pod created: " + pod.GetName()
	return response
}

func DeletePod(podname string) PodResponse {
	cs, err := getInClusterClientSet()
	response := PodResponse { "failure", "" }
	if err != nil {
		response.Message = err.Error()
		return response
	}

	podClient := cs.CoreV1().Pods("azuredevops")

	err2 := podClient.Delete(podname, &metav1.DeleteOptions{})
	if err2 != nil {
		response.Message = "podclient delete error: " + err2.Error()
		return response
	}

	response.Status = "success"
	response.Message = "Deleted " + podname
	return response
}

func DeletePodWithAgentId(agentId string) PodResponse {
	cs, err := getInClusterClientSet()
	response := PodResponse { "failure", "" }
	if err != nil {
		response.Message = err.Error()
		return response
	}

	podClient := cs.CoreV1().Pods("azuredevops")

	// Get the pod with this agentId
	pods, _ := podClient.List(metav1.ListOptions{LabelSelector: agentId})
	if(pods == nil || len(pods.Items) == 0) {
		response.Message = "Could not find running pod with AgentId" + agentId
		return response
	}

	err2 := podClient.Delete(pods.Items[0].GetName(), &metav1.DeleteOptions{})
	if err2 != nil {
		response.Message = "podclient delete error: " + err2.Error()
		return response
	}

	response.Status = "success"
	response.Message = "Deleted " + pods.Items[0].GetName()
	return response
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
	if(podname == "") {
		podname = "agent-dind"
	}

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