package main

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

var podnamespace = "azuredevops"

func main() {

	// Define HTTP endpoints
	s := http.NewServeMux()

	podnamespace = os.Getenv("POD_NAMESPACE")

	s.HandleFunc("/acquire", func(w http.ResponseWriter, r *http.Request) { AcquireAgentHandler(w, r) })
	s.HandleFunc("/release", func(w http.ResponseWriter, r *http.Request) { ReleaseAgentHandler(w, r) })
	s.HandleFunc("/buildPod", func(w http.ResponseWriter, r *http.Request) { GetBuildPodHandler(w, r) })

	// Start HTTP Server with request logging
	log.Fatal(http.ListenAndServe(":8080", s))
}

func AcquireAgentHandler(resp http.ResponseWriter, req *http.Request) {
	// HTTP method should be POST and the HMAC header should be valid
	if req.Method == http.MethodPost {
		log.Println("Recieved agent acquire request ....")
		if isRequestHmacValid(req) {
			log.Println("Hmac Validated for acquire request")
			var agentRequest AgentRequest

			requestBody, err := ioutil.ReadAll(req.Body)
			json.Unmarshal(requestBody, &agentRequest)

			if err != nil {
				writeJsonResponse(resp, http.StatusBadRequest, err.Error())
			} else if agentRequest.AgentId == "" {
				writeJsonResponse(resp, http.StatusBadRequest, GetError(NoAgentIdError))
			} else {
				log.Println("Calling create pod")
				var pods = CreatePod(agentRequest, podnamespace)
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
		log.Println("Recieved release agent request ....")
		if isRequestHmacValid(req) {
			log.Println("Hmac Validated for release request")
			var agentRequest ReleaseAgentRequest
			requestBody, _ := ioutil.ReadAll(req.Body)
			json.Unmarshal(requestBody, &agentRequest)

			if agentRequest.AgentId == "" {
				writeJsonResponse(resp, http.StatusBadRequest, GetError(NoAgentIdError))
			} else {
				log.Println("Calling delete pod")
				var pods = DeletePodWithAgentId(agentRequest.AgentId, podnamespace)
				writeJsonResponse(resp, http.StatusCreated, pods)
			}
		} else {
			writeJsonResponse(resp, http.StatusForbidden, GetError(NoValidSignatureError))
		}
	} else {
		writeJsonResponse(resp, http.StatusMethodNotAllowed, GetError(InvalidRequestError))
	}
}

func GetBuildPodHandler(resp http.ResponseWriter, req *http.Request) {

	log.Println("Recieved GetBuildPod request ....")
	if req.Method == http.MethodGet {
		if isRequestHmacValid(req) {
			log.Println("Hmac Validated for buildpod request")
			keyHeader := "key"
			headerVal := req.Header.Get(keyHeader)
			log.Println("Calling getbuildkit pod")
			var pods = GetBuildKitPod(headerVal, podnamespace)
			writeJsonResponse(resp, http.StatusCreated, pods)
		} else {
			writeJsonResponse(resp, http.StatusForbidden, GetError(NoValidSignatureError))
		}
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
