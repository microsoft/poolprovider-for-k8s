package main

import (
	"testing"
	"os"
	v1 "k8s.io/api/core/v1"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

func TestCreatePod(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	//clientset := fake.NewSimpleClientset()
	testPod := CreatePod(agentrequest)
	//secret := createSecret(&clientset,agentrequest)
	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

}

func TestCreateSecret(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()
	cs := CreateClientSet()

	testSecret := createSecret(cs,agentrequest)

	if (testSecret.Name == "newname"){
		t.Errorf("Secret creation failed")
	}

}

func TestDeletePodShouldPassIfMatchingAgentIdinAgentRequest(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	testPod := CreatePod(agentrequest)

	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

	testDeletepod := DeletePodWithAgentId(agentrequest.AgentId);
	if (testDeletepod.Status != "success"){
		t.Errorf("Pod deletion failed")
	}

}

func TestDeletePodShouldFailIfNotMatchingAgentIdinAgentRequest(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	testPod := CreatePod(agentrequest)
	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

	//Trying to delete pod with AgentId = 2
	testDeletepod := DeletePodWithAgentId("2");
	if (testDeletepod.Status != "fail"){
		t.Errorf("Pod deletion passed but should have failed")
	}

}

func TestGetBuildPodShouldReturnEmptyStringIfNoBuildKitPodPresent(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	testPod := CreatePod(agentrequest)
	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

	testDeletepod := GetBuildKitPod("test");
	if (testDeletepod.Message != ""){
		t.Errorf("Test failed")
	}

}

func TestGetBuildPodShouldReturnBuildKitPodNameIfPresent(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

    CreateDummyBuildKitPod()	
	
	testGetBuildpod := GetBuildKitPod("test");
	if (testGetBuildpod.Message == ""){
		t.Errorf("Test failed")
	}

}

func CreateDummyBuildKitPod() {
	cs := CreateClientSet()
	var p1 v1.Pod

	podname := "agent-lean-dind"

	dat, _ := ioutil.ReadFile("agentpods/" + podname + ".yaml")
	var podYaml = string(dat)
	_ = yaml.Unmarshal([]byte(podYaml), &p1)
		p1.SetLabels(map[string]string{
		"role": "buildkit",
		})
	
	p1.ObjectMeta.Name = "buildkitd-0"
	podClient := cs.clientset.CoreV1().Pods("azuredevops")
	_, _ = podClient.Create(&p1)
}

func SetTestingEnvironmentVariables() {
	os.Setenv("COUNTTEST","1")
}