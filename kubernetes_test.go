package main

import (
	"fmt"
	"testing"
	"os"
)

func TestCreatePod(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	os.Setenv("TESTING","True")

	//clientset := fake.NewSimpleClientset()
	pod := CreatePod(agentrequest)
	//secret := createSecret(&clientset,agentrequest)
	if (pod.Accepted == true){
		fmt.Println("Pod created", pod)
	}

}