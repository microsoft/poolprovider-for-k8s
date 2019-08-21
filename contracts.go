package main

type PodResponse struct {
    Status string
    Message  string
}

type AgentConfigurationData struct {
	AgentSettings string
	AgentCredentials string
	AgentVersion string
	AgentDownloadUrls string
}

type AgentRequest struct {
    AgentId string
	AgentPool  string
	AccountId string
	FailRequestUrl string
	AppendRequestMessageUrl string
	IsScheduled bool
	IsPublic bool
	AgentConfiguration AgentConfigurationData
	AgentSpec string
}

type AcquireAgentResponse struct {
	AgentData string
	Success bool
}

type ReleaseAgentRequest struct {
	AgentId string
	AccountId string
	AgentPool string
	AgentData string
}