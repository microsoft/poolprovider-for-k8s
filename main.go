package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
)

func main() {
	// Create Redis storage
	// storage := NewRedisStorage("k8s-poolprovider-redis:6379")

	// Define HTTP endpoints
	s := http.NewServeMux()
	s.HandleFunc("/definitions", func(w http.ResponseWriter, r *http.Request) { EmptyResponeHandler(w, r) })
	s.HandleFunc("/acquire", func(w http.ResponseWriter, r *http.Request) { AcquireAgentHandler(w, r) })
	s.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) { ReleaseAgentHandler(w, r) })

	// Start HTTP Server with request logging
	log.Fatal(http.ListenAndServe(":8082", s))
}

func AcquireAgentHandler(resp http.ResponseWriter, req *http.Request) {
	// HTTP method should be POST and the HMAC header should be valid
	if (req.Method == "POST") {
		if (isRequestHmacValid(req)) {
		    var agentRequest AgentRequest
		    requestBody, _ := ioutil.ReadAll(req.Body)
		    json.Unmarshal(requestBody, &agentRequest)

		    if(agentRequest.AgentId == "") {
			    http.Error(resp, "No AgentId sent in request body.", http.StatusCreated)
		    }

		    var pods = CreatePod(agentRequest.AgentId)
		    writeJsonResponse(resp, pods)
	    } else{
			http.Error(resp, "Endpoint can only be invoked with AzureDevOps with the correct Shared Signature.", http.StatusForbidden)
		}
	} else {
		http.Error(resp, "Invalid request Method.", http.StatusMethodNotAllowed)
	}
}

func ReleaseAgentHandler(resp http.ResponseWriter, req *http.Request) {
	// HTTP method should be POST and the HMAC header should be valid
	if (req.Method == "POST") {
		if (isRequestHmacValid(req)) {
			var agentRequest ReleaseAgentRequest
		    requestBody, _ := ioutil.ReadAll(req.Body)
		    json.Unmarshal(requestBody, &agentRequest)

		    if(agentRequest.AgentId == "") {
			    http.Error(resp, "No AgentId sent in request body.", http.StatusCreated)
		    }

		    var pods = DeletePodWithAgentId(agentRequest.AgentId)
		    writeJsonResponse(resp, pods)
		} else {
			http.Error(resp, "Endpoint can only be invoked with AzureDevOps with the correct Shared Signature.", http.StatusForbidden)
		}	
	} else {
		http.Error(resp, "Invalid request Method.", http.StatusMethodNotAllowed)
	}
}

func EmptyResponeHandler(resp http.ResponseWriter, req *http.Request) {
	var emptyResponse PodResponse
	writeJsonResponse(resp, emptyResponse)
}

func writeJsonResponse(resp http.ResponseWriter, podResponse PodResponse) {
	jsonData, _ := json.Marshal(podResponse)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
    resp.Write(jsonData)
}

func isRequestHmacValid(req *http.Request) bool {
	azureDevOpsHeader := "X-Azure-Signature"
	headerVal := req.Header.Get(azureDevOpsHeader)
	requestBody, _ := ioutil.ReadAll(req.Body)
	
	// No header is specified
	if (headerVal == "") {
		return false
	}

	// Compute HMAC for body and compare against the one sent by azure dev ops
	return ValidateHash(string(requestBody), headerVal)
}