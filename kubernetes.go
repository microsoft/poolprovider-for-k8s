package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"github.com/ghodss/yaml"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// External callers calling into Kubernetes APIs via this package will get a PodResponse
type PodResponse struct {
	Status  string
	Message string
}

const agentIdLabel = "AgentId"

// Creates a Pod with the default image specification. The pod is labelled with the agentId passed to it.
func CreatePod(agentRequest AgentRequest) AgentProvisionResponse {
	cs, err := GetClientSet()

	var response AgentProvisionResponse
	if err != nil {
		return getFailureResponse(response, err)
	}

	secret := createSecret(cs, agentRequest)
	pod, err := getAgentSpecification(agentRequest.AgentId)
	if err != nil {
		return getFailureResponse(response, err)
	}

	// Mount the secrets as a volume
	pod.Spec.Volumes = append(pod.Spec.Volumes, *getSecretVolume(secret.Name))

	//append(p1.Spec.Containers[0].Env, v1.EnvVar{Name: "VSTS_TOKEN", Value: token})

	podClient := cs.CoreV1().Pods("azuredevops")
	_, err2 := podClient.Create(pod)
	if err2 != nil {
		return getFailureResponse(response, err)
	}

	response.Accepted = true
	response.ResponseType = "Success"
	return response
}

func DeletePod(podname string) PodResponse {
	cs, err := GetClientSet()
	response := PodResponse{"failure", ""}
	if err != nil {
		response.Message = err.Error()
		return response
	}

	podClient := cs.CoreV1().Pods("azuredevops")

	err2 := podClient.Delete(podname, &metav1.DeleteOptions{})
	if err2 != nil {
		response.Message = "podclient delete error: " + err2.Error()
		return response
	}

	response.Status = "success"
	response.Message = "Deleted " + podname
	return response
}

func DeletePodWithAgentId(agentId string) PodResponse {
	cs, err := GetClientSet()
	var response PodResponse
	if err != nil {
		return getFailure(response, err)
	}

	podClient := cs.CoreV1().Pods("azuredevops")

    secretClient := cs.CoreV1().Secrets("azuredevops")
	
	// Get the secret with this agentId
	secrets, _ := secretClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + agentId})
	if secrets == nil || len(secrets.Items) == 0 {
		return getFailure(response, errors.New("Could not find running pod with AgentId"+agentId))
	}
	
	// Get the pod with this agentId
	pods, _ := podClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + agentId})
	if pods == nil || len(pods.Items) == 0 {
		return getFailure(response, errors.New("Could not find running pod with AgentId"+agentId))
	}

	err1 := secretClient.Delete(secrets.Items[0].GetName(), &metav1.DeleteOptions{})
	if err1 != nil {
		return getFailure(response, err1)
	}

	err2 := podClient.Delete(pods.Items[0].GetName(), &metav1.DeleteOptions{})
	if err2 != nil {
		return getFailure(response, err2)
	}

	response.Status = "success"
	response.Message = "Deleted " + pods.Items[0].GetName() + " and secret "+ secrets.Items[0].GetName()
	return response
}

func getAgentSpecification(agentId string) (*v1.Pod, error) {
	// Defaulting to use the DIND image, the podname can be exposed as a parameter and the user can then select which
	// image will be used to create the agent.
	podname := "agent-lean-dind"

	// If pod is to be created in a different namespace
	// then secrets need to be created in the same namespace, i.e. VSTS_TOKEN and VSTS_ACCOUNT
	// kubectl create secret generic vsts --from-literal=VSTS_TOKEN=<token> --from-literal=VSTS_ACCOUNT=<accountname>
	dat, _ := ioutil.ReadFile("agentpods/" + podname + ".yaml")

	var p1 v1.Pod
	var podYaml = string(dat)
	_ = yaml.Unmarshal([]byte(podYaml), &p1)

	if agentId != "" {
		// Set the agentId as label if specified
		p1.SetLabels(map[string]string{
			agentIdLabel: agentId,
		})
	}

	return &p1, nil
}

func getAgentSecret() *v1.Secret {
	var secret v1.Secret

	dat, _ := ioutil.ReadFile("agentpods/agent-secret.yaml")
	var secretYaml = string(dat)
	yaml.Unmarshal([]byte(secretYaml), &secret)

	return &secret
}

func createSecret(cs *kubernetes.Clientset, request AgentRequest) *v1.Secret {
	secret := getAgentSecret()
	agentSettings, _ := json.Marshal(request.AgentConfiguration.AgentSettings)
	agentCredentials, _ := json.Marshal(request.AgentConfiguration.AgentCredentials)

	if request.AgentId != "" {
		// Set the agentId as label of the secret if specified
		secret.SetLabels(map[string]string{
			agentIdLabel: request.AgentId,
		})
	}

	secret.Data[".agent"] = ([]byte(string(agentSettings)))
	secret.Data[".credentials"] = ([]byte(string(agentCredentials)))
	secret.Data[".url"] = ([]byte(request.AgentConfiguration.AgentDownloadUrls["linux-x64"]))

	secretClient := cs.CoreV1().Secrets("azuredevops")
	secret2, err := secretClient.Create(secret)

	if err != nil {
		secret2.Name = "newname"
	}

	return secret2
}

func getSecretVolume(secretName string) *v1.Volume {
	return &v1.Volume{
		Name:         "agent-creds",
		VolumeSource: v1.VolumeSource{Secret: &v1.SecretVolumeSource{SecretName: secretName}},
	}
}

func getFailure(response PodResponse, err error) PodResponse {
	response.Status = "fail"
	response.Message = err.Error()
	return response
}

func getFailureResponse(response AgentProvisionResponse, err error) AgentProvisionResponse {
	response.ResponseType = "fail"
	response.ErrorMessage = err.Error()
	return response
}
