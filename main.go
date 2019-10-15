package main

import (
	"bytes"
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
	s.HandleFunc("/buildPod", func(w http.ResponseWriter, r *http.Request) { GetBuildPodHandler(w, r) })

	// Start HTTP Server with request logging
	log.Fatal(http.ListenAndServe(":8080", s))
}

func AcquireAgentHandler(resp http.ResponseWriter, req *http.Request) {
	// HTTP method should be POST and the HMAC header should be valid
	if req.Method == http.MethodPost {
		if isRequestHmacValid(req) {
			var agentRequest AgentRequest

			requestBody, err := ioutil.ReadAll(req.Body)
			json.Unmarshal(requestBody, &agentRequest)

			if err != nil {
				writeJsonResponse(resp, http.StatusBadRequest, err.Error())
			} else if agentRequest.AgentId == "" {
				writeJsonResponse(resp, http.StatusBadRequest, GetError(NoAgentIdError))
			} else {
				var pods = CreatePod(agentRequest)
				writeJsonResponse(resp, http.StatusCreated, pods)
			}
		} else {
			writeJsonResponse(resp, http.StatusForbidden, GetError(NoValidSignatureError))
		}
	} else {
		writeJsonResponse(resp, http.StatusMethodNotAllowed, GetError(InvalidRequestError))
	}
}

func ReleaseAgentHandler(resp http.ResponseWriter, req *http.Request) {
	if req.Method == http.MethodPost {
		if isRequestHmacValid(req) {
			var agentRequest ReleaseAgentRequest
			requestBody, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(requestBody, &agentRequest)

			if agentRequest.AgentId == "" {
				writeJsonResponse(resp, http.StatusBadRequest, GetError(NoAgentIdError))
			} else {
				var pods = DeletePodWithAgentId(agentRequest.AgentId)
				writeJsonResponse(resp, http.StatusCreated, pods)
			}
		} else {
			writeJsonResponse(resp, http.StatusForbidden, GetError(NoValidSignatureError))
		}
	} else {
		writeJsonResponse(resp, http.StatusMethodNotAllowed, GetError(InvalidRequestError))
	}
}

func EmptyResponeHandler(resp http.ResponseWriter, req *http.Request) {
	var emptyResponse PodResponse
	writeJsonResponse(resp, http.StatusCreated, emptyResponse)
}

func GetBuildPodHandler(resp http.ResponseWriter, req *http.Request) {

	if req.Method == http.MethodGet {
		keyHeader := "key"
		headerVal := req.Header.Get(keyHeader)
		var pods = GetBuildKitPod(headerVal)
		writeJsonResponse(resp, http.StatusCreated, pods)
	} else {
		writeJsonResponse(resp, http.StatusMethodNotAllowed, GetError(InvalidRequestError))
	}
}

func writeJsonResponse(resp http.ResponseWriter, httpStatus int, podResponse interface{}) {
	jsonData, _ := json.Marshal(podResponse)
	resp.Header().Set("Content-Type", "application/json")
	resp.WriteHeader(httpStatus)
	resp.Write(jsonData)
}

func isRequestHmacValid(req *http.Request) bool {
	azureDevOpsHeader := "X-Azure-Signature"
	headerVal := req.Header.Get(azureDevOpsHeader)
	requestBody, _ := ioutil.ReadAll(req.Body)

	// Set the body again
	req.Body = ioutil.NopCloser(bytes.NewBuffer(requestBody))

	// No header is specified
	if headerVal == "" {
		return false
	}

	// Compute HMAC for body and compare against the one sent by azure dev ops
	return ValidateHash(string(requestBody), headerVal)
}
