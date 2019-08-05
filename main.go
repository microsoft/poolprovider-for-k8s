package main

import (
	"fmt"
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
	s.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) { KubernetesCreateHandler(w, r) })
	s.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) { KubernetesDeleteHandler(w, r) })

	// Test redis
	s.HandleFunc("/testredisdata", StorageSetHandler(storage))
	s.HandleFunc("/redisgetkeys", GetKeysHandler(storage))

	// Start HTTP Server with request logging
	log.Fatal(http.ListenAndServe(":8082", s))
}

func KubernetesCreateHandler(resp http.ResponseWriter, req *http.Request) {
	userSpec := req.URL.Query()["agentspec"]

	// using the default agent spec
	agentSpec := "agent-dind"
	if userSpec != nil {
		agentSpec = userSpec[0]
	}

	var pods = CreatePod(agentSpec)
	fmt.Fprintf(resp, "Pods: %s", pods)
}

func KubernetesDeleteHandler(resp http.ResponseWriter, req *http.Request) {
	podname := req.URL.Query()["podname"]
	if podname == nil || podname[0] == "" {
		fmt.Fprintf(resp, "Provide pod name as ?podname=somename");
		return
	}

	var pods = DeletePod(podname[0])
	fmt.Fprintf(resp, "Response: %s", pods)
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