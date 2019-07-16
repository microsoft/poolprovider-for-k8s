package main

import (
	"fmt"
	"log"
	"net/http"
)

const (
	// Name of the application
	Name = "Argus"
	// Version of the application
	Version = "1.0.0"
)

func main() {
	// Define HTTP endpoints
	s := http.NewServeMux()
	s.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) { KubernetesCreateHandler(w, r) })
	s.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) { KubernetesDeleteHandler(w, r) })

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
