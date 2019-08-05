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
	s.HandleFunc("/storageget", StorageGetHandler(storage))
	s.HandleFunc("/storageset", StorageSetHandler(storage))
	s.HandleFunc("/storageping", PingHandler(storage))
	s.HandleFunc("/storagegetkeys", GetKeysHandler(storage))

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

func StorageGetHandler(s Storage) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		key := req.URL.Query()["key"]
		res, err := s.Get(key[0])
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(resp, err.Error())
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintln(resp, res)
	}
}

func StorageSetHandler(s Storage) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		key := req.URL.Query()["key"]
		value := req.URL.Query()["value"]
		fmt.Fprintf(resp, "key is %s, value is %s", key[0], value[0]);
		if key == nil || key[0] == "" {
			key[0] = "fu"
		}

		if value == nil || value[0] == "" {
			value[0] = "bar"
		}

		err := s.Set(key[0], value[0])
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(resp, "Error" + err.Error())
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintf(resp, "Value set")
	}
}

func PingHandler(s Storage) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		res := s.Init()

		resp.WriteHeader(http.StatusOK)
		fmt.Fprintln(resp, res)
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