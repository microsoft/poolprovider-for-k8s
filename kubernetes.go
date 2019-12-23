package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"

	"log"

	v1alpha1 "github.com/microsoft/k8s-poolprovider/pkg/apis/dev/v1alpha1"

	"github.com/ghodss/yaml"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// External callers calling into Kubernetes APIs via this package will get a PodResponse
type PodResponse struct {
	Status  string
	Message string
}

type k8s struct {
	clientset kubernetes.Interface
}

const agentIdLabel = "AgentId"

// Creates a Pod with the default image specification. The pod is labelled with the agentId passed to it.
func CreatePod(agentRequest AgentRequest, podnamespace string) AgentProvisionResponse {

	var config *rest.Config
	config, _= rest.InClusterConfig();

	crdclient, err := v1alpha1.NewClient(config)
	log.Println("rest client inside getspecification \n", crdclient)

	resp1, err := crdclient.AzurePipelinesPool(podnamespace).Get("azurepipelinespool-operator")
	if err != nil {
		log.Println("error fetching podconfig", err)
	} else {
		log.Println("resp \n", resp1)
	}

	labels := GenerateLabelsForPod(agentRequest.AgentId)

	log.Println("Add a Linux Pod using CRD")
	pod := crdclient.AzurePipelinesPool(podnamespace).AddNewPodForCR(resp1, agentRequest.AgentId, labels, "linux")

	log.Println("pod created ", pod)

	cs := CreateClientSet()

	log.Println("Starting pod creation")
	var response AgentProvisionResponse

	podClient := cs.clientset.CoreV1().Pods(podnamespace)
	webserverpod, _ := podClient.List(metav1.ListOptions{LabelSelector: "app=azurepipelinepool-operator"})

	AddOwnerRefToObject(pod, AsOwner(&webserverpod.Items[0]))
	log.Println("Webserver pod added as owner reference to agent pod ")

	log.Println("create secret called")
	secret := createSecret(cs, agentRequest, &webserverpod.Items[0])

	if err != nil {
		return getFailureResponse(response, err)
	}

	// Mount the secrets as a volume
	pod.Spec.Volumes = append(pod.Spec.Volumes, *getSecretVolume(secret.Name))
	log.Println("secrets mounted as volume")

	_, err2 := podClient.Create(pod)
	if err2 != nil {
		return getFailureResponse(response, err)
	}

	log.Println("Pod creation done")

	response.Accepted = true
	response.ResponseType = "Success"
	return response
}

func GetBuildKitPod(key string, podnamespace string) PodResponse {
	cs := CreateClientSet()

	var response PodResponse

	listOptions := metav1.ListOptions{
		LabelSelector: "role=buildkit",
	}

	podClient := cs.clientset.CoreV1().Pods(podnamespace)
	podlist, err2 := podClient.List(listOptions)
	if err2 != nil {
		return getFailure(response, err2)
	}
	log.Println("Fetched list of pods configured as buildkit stateful pods")
	var nodes []string

	for _, items := range podlist.Items {
		s := items.GetName()
		if s != "" {
			nodes = append(nodes, s)
		}
	}

	log.Println("Fetching the target pod using consistent hash")
	chosen := ComputeConsistentHash(nodes, key)
	response.Status = "success"
	response.Message = chosen
	return response
}

func DeletePodWithAgentId(agentId string, podnamespace string) PodResponse {
	cs := CreateClientSet()
	var response PodResponse

	podClient := cs.clientset.CoreV1().Pods(podnamespace)

	secretClient := cs.clientset.CoreV1().Secrets(podnamespace)

	// Get the secret with this agentId
	secrets, _ := secretClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + agentId})
	if secrets == nil || len(secrets.Items) == 0 {
		return getFailure(response, errors.New("Could not find secret with AgentId "+agentId))
	}

	// Get the pod with this agentId
	pods, _ := podClient.List(metav1.ListOptions{LabelSelector: agentIdLabel + "=" + agentId})
	if pods == nil || len(pods.Items) == 0 {
		return getFailure(response, errors.New("Could not find running pod with AgentId "+agentId))
	}

	err1 := secretClient.Delete(secrets.Items[0].GetName(), &metav1.DeleteOptions{})
	if err1 != nil {
		return getFailure(response, err1)
	}
	log.Println("Delete agent secret done")

	err2 := podClient.Delete(pods.Items[0].GetName(), &metav1.DeleteOptions{})
	if err2 != nil {
		return getFailure(response, err2)
	}
	log.Println("Delete agent pod done")

	response.Status = "success"
	response.Message = "Deleted " + pods.Items[0].GetName() + " and secret " + secrets.Items[0].GetName()
	return response
}

func getAgentSecret() *v1.Secret {
	var secret v1.Secret

	log.Println("Reading agent-secret.yaml")
	dat, _ := ioutil.ReadFile("agentpods/agent-secret.yaml")
	var secretYaml = string(dat)
	yaml.Unmarshal([]byte(secretYaml), &secret)

	return &secret
}

func createSecret(cs *k8s, request AgentRequest, m *v1.Pod) *v1.Secret {
	secret := getAgentSecret()

	log.Println("parsing secret data from agent request")
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
	secret.Data[".agentVersion"] = ([]byte(request.AgentConfiguration.AgentVersion))

	AddOwnerRefToObject(secret, AsOwner(m))
	log.Println("WebServer pod added as Owner reference to secret")

	secretClient := cs.clientset.CoreV1().Secrets(podnamespace)
	secret2, err := secretClient.Create(secret)

	if err != nil {
		secret2.Name = "newname"
	}
	log.Println("Secret creation done")
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

func GenerateLabelsForPod(agentId string) map[string]string {
	return map[string]string{agentIdLabel: agentId}
}

func AddOwnerRefToObject(obj metav1.Object, ownerRef metav1.OwnerReference) {
	obj.SetOwnerReferences(append(obj.GetOwnerReferences(), ownerRef))
}

// asOwner returns an OwnerReference set as the memcached CR
func AsOwner(m *v1.Pod) metav1.OwnerReference {
	falseVar := false
	return metav1.OwnerReference{
		APIVersion: "apps/v1",
		Kind:       "Pod",
		Name:       m.Name,
		UID:        m.UID,
		Controller: &falseVar,
	}
}