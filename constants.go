package main

const (
	NoAgentIdError        = "No AgentId sent in request body."
	NoValidSignatureError = "Endpoint can only be invoked with AzureDevOps with the correct Shared Signature."
	InvalidRequestError   = "Invalid request Method."
	AgentLabel            = "AgentId"
)

type ErrorMessage struct {
	Error string
}

func GetError(message string) ErrorMessage {
	return ErrorMessage{Error: message}
}
