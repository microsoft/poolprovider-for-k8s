package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)

const (
	// Name of the application
	Name = "divman's GoServer"
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
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pods = CreatePod()

	fmt.Fprintf(resp, "Pods: %s", pods)
}

func KubernetesDeleteHandler(resp http.ResponseWriter, req *http.Request) {
	podname := req.URL.Query()["podname"][0]
	if podname == "" {
		fmt.Fprintf(resp, "Provide pod name as ?podname=somename");
		return
	}

	var pods = DeletePod(podname)

	fmt.Fprintf(resp, "Response: %s", pods)
}
