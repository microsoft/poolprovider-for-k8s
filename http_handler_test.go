package main

import (
	"testing"
	"bytes"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"net/http"
	"net/http/httptest"
	"encoding/json"

)

func TestAcquireHandlerShouldBeSuccessful(t *testing.T) {

	var response AgentProvisionResponse
	var jsonStr = []byte(`{"AgentId":"1"}`)
	
	SetupCustomResource()
	req, _ := http.NewRequest("POST", "/acquire", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(AcquireAgentHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusCreated { //Must be 201
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusCreated, status)
	} else {

		cs := CreateClientSet()
		podClient := cs.clientset.CoreV1().Pods("azuredevops")
		pods, _ := podClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + "1"})

		if pods == nil || len(pods.Items) == 0 {
			t.Errorf("Http Aquire Call failed")
		} else {

			json.Unmarshal([]byte(resp.Body.String()), &response)

			if response.Accepted != true {
				t.Errorf("Http Aquire Call failed")
			} else if response.ResponseType != "Success" {
				t.Errorf("Http Aquire Call failed")
			} else if response.ErrorMessage != "" {
				t.Errorf("Http Aquire Call failed")
			}
		}
	}
}

func TestAcquireHandlerShouldFailIfGetRequest(t *testing.T) {
	SetupCustomResource()

	var jsonStr = []byte(`{"AgentId":"1"}`)

	req, _ := http.NewRequest("GET", "/acquire", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(AcquireAgentHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusMethodNotAllowed { //Must be 405
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusMethodNotAllowed, status)
	}
}

func TestAcquireHandlerShouldFailIfHmacNotValid(t *testing.T) {
	SetupCustomResource()

	var jsonStr = []byte(`{"AgentId":"12"}`)
	
	req, _ := http.NewRequest("POST", "/acquire", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")

	// wrong encoding
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(AcquireAgentHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusForbidden { //Must be 403
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusForbidden, status)
	}
}

func TestReleaseHandlerShouldBeSuccessful(t *testing.T) {
	SetupCustomResource()
	var agentrequest AgentRequest
	agentrequest.AgentId = "1"
	testPod := CreatePod(agentrequest, "azuredevops")

	if (testPod.Accepted != true){
		t.Errorf("Pod creation failed")
	}

	var response PodResponse
	var jsonStr = []byte(`{"AgentId":"1"}`)
	
	req, _ := http.NewRequest("POST", "/release", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(ReleaseAgentHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusCreated { //Must be 201
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusCreated, status)
	} else {

		// Now check the pod is deleted or not
		cs := CreateClientSet()
		podClient := cs.clientset.CoreV1().Pods("azuredevops")
		pods, _ := podClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + "1"})

		if pods == nil || len(pods.Items) == 1 {
			t.Errorf("Http Release Call failed")
		} else {
			json.Unmarshal([]byte(resp.Body.String()), &response)

			if response.Status != "success" {
				t.Errorf("Http Release Call failed")
			} else if response.Message == "" {
				t.Errorf("Http Release Call failed")
			}
		}
	}
}

func TestGetBuildPodHandlerShouldBeSuccessful(t *testing.T) {
	SetTestingEnvironmentVariables()
	CreateDummyBuildKitPod()

	var response PodResponse
	var jsonStr = []byte("")
	req, _ := http.NewRequest("GET", "/buildPod", bytes.NewBuffer(jsonStr))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("X-Azure-Signature", "4f6a97c5aa13477ed775dd20cdd7cf44477e310ba683545144d5112e77d88be967a43791da0696ded702a32ca0c190ab831dcd9204521b9a9ebe413066699ef9")

	resp := httptest.NewRecorder()
	http.HandlerFunc(GetBuildPodHandler).ServeHTTP(resp,req)

	if status := resp.Code; status != http.StatusCreated { //Must be 201
		t.Errorf("Status code differs. Expected %d. Got %d", http.StatusCreated, status)
	} else {
			json.Unmarshal([]byte(resp.Body.String()), &response)

			if response.Status != "success" {
				t.Errorf("Http buildPod Call failed")
			} else if response.Message != "buildkitd-0" {
				t.Errorf("Http buildPod Call failed")
			}
	}
}