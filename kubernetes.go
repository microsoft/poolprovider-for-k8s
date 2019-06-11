package main

import (
	"github.com/ghodss/yaml"
	
	"k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func CreatePod() string {
	config, err := rest.InClusterConfig()
	if err != nil {
		return err.Error()
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return err.Error()
	}

	var podYaml = `
apiVersion: v1
kind: Pod
metadata:
   name: vsts-agent-dind
spec:
  containers:
  - name: vsts-agent
    image: microsoft/vsts-agent:ubuntu-16.04-docker-18.06.1-ce-standard
    env:
    - name: VSTS_ACCOUNT
      valueFrom:
        secretKeyRef:
          name: vsts
          key: VSTS_ACCOUNT
    - name: VSTS_TOKEN
      valueFrom:
        secretKeyRef:
          name: vsts
          key: VSTS_TOKEN
    - name: VSTS_POOL
      value: divman
    - name: DOCKER_HOST
      value: tcp://localhost:2375
  - name: dind-daemon
    image: docker:18.09.6-dind
    securityContext:
      privileged: true
    volumeMounts:
    - name: daemon-storage
      mountPath: /var/lib/docker
  volumes:
  - name: daemon-storage
    emptyDir: {}
	`

	var p1 v1.Pod
	err1 := yaml.Unmarshal([]byte(podYaml), &p1)
	if err1 != nil {
		return err1.Error()
	}

	// running the app in the default namespace. Pass namespace to pods method.
	podClient := clientset.CoreV1().Pods("")
	pod, err2 := podClient.Create(&p1)
	if err2 != nil {
		return err2.Error()
	}

	return pod.GetName()
}