package main

import (
	"testing"
	"strings"
	"os"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"io/ioutil"

	"github.com/ghodss/yaml"
)

func TestCreatePod(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	testPod := CreatePod(agentrequest)

	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

}

func TestCreatePodMustCreateSecret(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	testPod := CreatePod(agentrequest)

	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

	cs := CreateClientSet()
	secretClient := cs.clientset.CoreV1().Secrets("azuredevops")
	
	// Get the secret with this agentId
	secrets, _ := secretClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + agentrequest.AgentId})
	if secrets == nil || len(secrets.Items) == 0 {
		t.Errorf("Could not find secret with AgentId "+agentrequest.AgentId)
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

func TestAgentSpecificationMustAssignRandomNameToAgentPod(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	agentVersion := "2.158.0"
	SetTestingEnvironmentVariables()
	pod, _ :=getAgentSpecification(agentrequest.AgentId, agentVersion)

	if (!strings.Contains(pod.ObjectMeta.GenerateName,"vsts-agent-")){
		t.Errorf("Agent pod not generated with random name")
	}

}

func TestCreateSecretMustHaveAllDataValues(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()
	cs := CreateClientSet()

	testSecret := createSecret(cs,agentrequest)

	if _, ok := testSecret.Data[".agent"] ; !ok {
		t.Errorf("Secret doesn't have .agent data")
	}

	if _, ok := testSecret.Data[".credentials"] ; !ok {
		t.Errorf("Secret doesn't have .credentials data")
	}

	if _, ok := testSecret.Data[".url"] ; !ok {
		t.Errorf("Secret doesn't have .url data")
	}

	if _, ok := testSecret.Data[".agentVersion"] ; !ok {
		t.Errorf("Secret doesn't have .agentVersion data")
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

func TestDeletePodMustDeleteSecret(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	testPod := CreatePod(agentrequest)
	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

	testDeletepod := DeletePodWithAgentId(agentrequest.AgentId);
	if (testDeletepod.Status != "success"){
		t.Errorf("Pod deletion passed but should have failed")
	}

	cs := CreateClientSet()
	secretClient := cs.clientset.CoreV1().Secrets("azuredevops")
	
	// Get the secret with this agentId
	secrets, _ := secretClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + agentrequest.AgentId})
	if secrets == nil || len(secrets.Items) == 01{
		t.Errorf("Secret not deleted found secret with AgentId "+agentrequest.AgentId)
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