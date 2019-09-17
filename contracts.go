package main

type AgentConfigurationData struct {
	AgentSettings     string
	AgentCredentials  string
	AgentVersion      string
	AgentDownloadUrls string
}

type AgentRequest struct {
	AgentId                 string
	AgentPool               string
	AccountId               string
	FailRequestUrl          string
	AppendRequestMessageUrl string
	IsScheduled             bool
	IsPublic                bool
	AgentConfiguration      AgentConfigurationData
	AgentSpec               string
}

type AgentProvisionResponse struct {
	AgentData    string
	ResponseType string
	ErrorMessage string
}

type ReleaseAgentRequest struct {
	AgentId   string
	AccountId string
	AgentPool string
	AgentData string
}
