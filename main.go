package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

const (
	// Name of the application
	Name = "k8s-poolprovider"
	// Version of the application
	Version = "1.0.0"
)

func main() {
	// Create Redis storage
	storage := NewRedisStorage("k8s-poolprovider-redis:6379")

	// Define HTTP endpoints
	s := http.NewServeMux()
	s.HandleFunc("/definitions", func(w http.ResponseWriter, r *http.Request) { EmptyResponeHandler(w, r) })
	s.HandleFunc("/acquire", func(w http.ResponseWriter, r *http.Request) { AcquireAgentHandler(w, r) })
	s.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) { ReleaseAgentHandler(w, r) })

	// Test redis
	s.HandleFunc("/testredisdata", StorageSetHandler(storage))
	s.HandleFunc("/redisgetkeys", GetKeysHandler(storage))

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

		var pods = CreatePod("", agentRequest.AgentId)
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

func StorageSetHandler(s Storage) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		key := "some sample key"
		value := "some sample value"

		// Retrieving information from backing Redis storage
		s.Set(key, value)
		retrievedValue, _ := s.Get(key)
		fmt.Fprintf(resp, "All good. Retrieved %s", retrievedValue)
	}
}

func GetKeysHandler(s Storage) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		res, err := s.GetKeys("*")
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(resp, err.Error())
			return
		}

		resp.WriteHeader(http.StatusOK)
		fmt.Fprintln(resp, strings.Join(res, ", "))
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