package main

type AgentConfigurationData struct {
	AgentSettings     map[string]string
	AgentCredentials  AgentCredentials
	AgentVersion      string
	AgentDownloadUrls map[string]string
}

type AgentCredentials struct {
	Scheme string
	Data   map[string]string
}

type AgentRequest struct {
	AgentId                 string
	AgentPool               string
	AccountId               string
	AuthenticationToken     string
	FailRequestUrl          string
	AppendRequestMessageUrl string
	IsScheduled             bool
	IsPublic                bool
	AgentConfiguration      AgentConfigurationData
	AgentSpec               string
}

type AgentProvisionResponse struct {
	Accepted     bool
	ResponseType string
	ErrorMessage string
}

type ReleaseAgentRequest struct {
	AgentId   string
	AccountId string
	AgentPool string
	AgentData string
}
