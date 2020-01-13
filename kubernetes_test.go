package main

import (
	"testing"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const testnamespace = "azuredevops"

func TestCreatePod(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetupCustomResource()

	testPod := CreatePod(agentrequest, testnamespace)

	if testPod.Accepted != true {
		t.Errorf("Pod creation failed")
	}

}

func TestCreatePodMustCreateSecret(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetupCustomResource()

	testPod := CreatePod(agentrequest, testnamespace)

	if testPod.Accepted != true {
		t.Errorf("Pod creation failed")
	}

	cs := CreateClientSet()
	secretClient := cs.clientset.CoreV1().Secrets("azuredevops")

	// Get the secret with this agentId
	secrets, _ := secretClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + agentrequest.AgentId})
	if secrets == nil || len(secrets.Items) == 0 {
		t.Errorf("Could not find secret with AgentId " + agentrequest.AgentId)
	}
}

func TestCreateSecret(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetupCustomResource()
	cs := CreateClientSet()

	testSecret := createSecret(cs, agentrequest, nil)

	if testSecret.Name == "newname" {
		t.Errorf("Secret creation failed")
	}

}

func TestCreateSecretMustHaveAllDataValues(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetupCustomResource()
	cs := CreateClientSet()

	testSecret := createSecret(cs, agentrequest, nil)

	if _, ok := testSecret.Data[".agent"]; !ok {
		t.Errorf("Secret doesn't have .agent data")
	}

	if _, ok := testSecret.Data[".credentials"]; !ok {
		t.Errorf("Secret doesn't have .credentials data")
	}

	if _, ok := testSecret.Data[".url"]; !ok {
		t.Errorf("Secret doesn't have .url data")
	}

	if _, ok := testSecret.Data[".agentVersion"]; !ok {
		t.Errorf("Secret doesn't have .agentVersion data")
	}
}

func TestDeletePodShouldPassIfMatchingAgentIdinAgentRequest(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetupCustomResource()

	testPod := CreatePod(agentrequest, testnamespace)

	if testPod.Accepted != true {
		t.Errorf("Pod creation failed")
	}

	testDeletepod := DeletePodWithAgentId(agentrequest.AgentId, testnamespace)
	if testDeletepod.Status != "success" {
		t.Errorf("Pod deletion failed")
	}

}

func TestDeletePodShouldFailIfNotMatchingAgentIdinAgentRequest(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	testPod := CreatePod(agentrequest, testnamespace)
	if testPod.Accepted != true {
		t.Errorf("Pod creation failed")
	}

	//Trying to delete pod with AgentId = 2
	testDeletepod := DeletePodWithAgentId("2", testnamespace)
	if testDeletepod.Status != "fail" {
		t.Errorf("Pod deletion passed but should have failed")
	}

}

func TestDeletePodMustDeleteSecret(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetupCustomResource()

	testPod := CreatePod(agentrequest, testnamespace)
	if testPod.Accepted != true {
		t.Errorf("Pod creation failed")
	}

	testDeletepod := DeletePodWithAgentId(agentrequest.AgentId, testnamespace)
	if testDeletepod.Status != "success" {
		t.Errorf("Pod deletion passed but should have failed")
	}

	cs := CreateClientSet()
	secretClient := cs.clientset.CoreV1().Secrets("azuredevops")

	// Get the secret with this agentId
	secrets, _ := secretClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + agentrequest.AgentId})
	if secrets == nil || len(secrets.Items) == 01 {
		t.Errorf("Secret not deleted found secret with AgentId " + agentrequest.AgentId)
	}
}

func TestGetBuildPodShouldReturnEmptyStringIfNoBuildKitPodPresent(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables()

	testPod := CreatePod(agentrequest, testnamespace)
	if testPod.Accepted != true {
		t.Errorf("Pod creation failed")
	}

	testDeletepod := GetBuildKitPod("test", testnamespace)
	if testDeletepod.Message != "" {
		t.Errorf("Test failed")
	}

}

func TestGetBuildPodShouldReturnBuildKitPodNameIfPresent(t *testing.T) {
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	SetTestingEnvironmentVariables(true)

	testGetBuildpod := GetBuildKitPod("test", testnamespace)
	if testGetBuildpod.Message == "" {
		t.Errorf("Test failed")
	}

}
