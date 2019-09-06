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
	if(req.Method == "POST") {
		var agentRequest AgentRequest
		requestBody, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(requestBody, &agentRequest)

		if(agentRequest.AgentId == "") {
			http.Error(resp, "No AgentId sent in request body.", http.StatusCreated)
		}

		var pods = CreatePod(agentRequest.AgentId)
		WriteJsonResponse(resp, pods)
	} else {
		http.Error(resp, "Invalid request Method.", http.StatusMethodNotAllowed)
	}
}

func ReleaseAgentHandler(resp http.ResponseWriter, req *http.Request) {
	if(req.Method == "POST") {
		var agentRequest ReleaseAgentRequest
		requestBody, _ := ioutil.ReadAll(req.Body)
		json.Unmarshal(requestBody, &agentRequest)

		if(agentRequest.AgentId == "") {
			http.Error(resp, "No AgentId sent in request body.", http.StatusCreated)
		}

		var pods = DeletePodWithAgentId(agentRequest.AgentId)

		WriteJsonResponse(resp, pods)
	} else {
		http.Error(resp, "Invalid request Method.", http.StatusMethodNotAllowed)
	}
}

func EmptyResponeHandler(resp http.ResponseWriter, req *http.Request) {
	var emptyResponse PodResponse
	WriteJsonResponse(resp, emptyResponse)
}

func WriteJsonResponse(resp http.ResponseWriter, podResponse PodResponse) {
	jsonData, _ := json.Marshal(podResponse)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(http.StatusCreated)
    resp.Write(jsonData)
}