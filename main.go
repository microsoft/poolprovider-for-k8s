package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"io/ioutil"

	"github.com/gorilla/handlers"
)

const (
	// Name of the application
	Name = "divman's GoServer"
	// Version of the application
	Version = "1.0.0"
)

func main() {
	var (
		listenFlag = flag.String("listen", EnvOrDefault("SIMPLE_WEBSERVER_LISTEN", ":8082"), "Address + Port to listen on. Format ip:port. Environment variable: SIMPLE_WEBSERVER_LISTEN")
		redisFlag  = flag.String("redis", EnvOrDefault("SIMPLE_WEBSERVER_REDIS", ":6379"), "Address + Port where a redis server is listening. Environment variable: SIMPLE_WEBSERVER_REDIS")
	)
	flag.Parse()

	// Create Redis storage
	storage := NewRedisStorage(*redisFlag)

	// Define HTTP endpoints
	s := http.NewServeMux()
	s.HandleFunc("/ping", PingHandler(storage))
	s.HandleFunc("/version", VersionHandler)
	s.HandleFunc("/payload", PayloadHandler)
	s.HandleFunc("/create", func(w http.ResponseWriter, r *http.Request) { KubernetesCreateHandler(w, r, storage) })
	s.HandleFunc("/delete", func(w http.ResponseWriter, r *http.Request) { KubernetesDeleteHandler(w, r, storage) })

	// Bootstrap logger
	logger := log.New(os.Stdout, "", log.LstdFlags)
	logger.Printf("Starting webserver and listen on %s", *listenFlag)

	// Start HTTP Server with request logging
	loggingHandler := handlers.LoggingHandler(os.Stdout, s)
	log.Fatal(http.ListenAndServe(*listenFlag, loggingHandler))
}

// PingHandler handles request to the "/ping" endpoint.
// It will send a PING request to Redis and return the response
// of the NoSQL database.
// The response is obvious: "pong" :)
func PingHandler(s Storage) http.HandlerFunc {
	return func(resp http.ResponseWriter, req *http.Request) {
		res, err := s.Ping()
		if err != nil {
			resp.WriteHeader(http.StatusInternalServerError)
			fmt.Fprintf(resp, err.Error())
			return
		}
		resp.WriteHeader(http.StatusOK)
		fmt.Fprintln(resp, res)
	}
}

// VersionHandler handles request to the "/version" endpoint.
// It prints the Name and Version of this app.
func VersionHandler(resp http.ResponseWriter, req *http.Request) {
	resp.WriteHeader(http.StatusOK)
	fmt.Fprintf(resp, "%s v%s\n", Name, Version)
}

// PayloadHandler handles request to the "/payload" endpoint.
// It is a debug route to dump the complete request incl. method, header and body.
func PayloadHandler(resp http.ResponseWriter, req *http.Request) {
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	resp.WriteHeader(http.StatusOK)
	log.Printf("Method: %s\n", req.Method)
	fmt.Fprintf(resp, "Method: %s\n", req.Method)

	if len(req.Header) > 0 {
		log.Println("Headers:")
		fmt.Fprint(resp, "Headers:\n")
		for key, values := range req.Header {
			for _, val := range values {
				log.Printf("%s: %s\n", key, val)
				fmt.Fprintf(resp, "%s: %s\n", key, val)
			}
		}
	}

	log.Printf("Payload: %s", string(body))
	fmt.Fprintf(resp, "Payload: %s", string(body))
}

func KubernetesCreateHandler(resp http.ResponseWriter, req *http.Request, s Storage) {
	_, err := ioutil.ReadAll(req.Body)
	if err != nil {
		resp.WriteHeader(http.StatusInternalServerError)
		return
	}

	var pods = CreatePod()

	fmt.Fprintf(resp, "Pods: %s", pods)
}

func KubernetesDeleteHandler(resp http.ResponseWriter, req *http.Request, s Storage) {
	podname := req.URL.Query()["podname"][0]
	if podname == "" {
		fmt.Fprintf(resp, "Provide pod name as ?podname=somename");
		return
	}

	var pods = DeletePod(podname)

	fmt.Fprintf(resp, "Response: %s", pods)
}

// EnvOrDefault will read env from the environment.
// If the environment variable is not set in the environment
// fallback will be returned.
// This function can be used as a value for flag.String to enable
// env var support for your binary flags.
func EnvOrDefault(env, fallback string) string {
	value := fallback
	if v := os.Getenv(env); v != "" {
		value = v
	}

	return value
}
