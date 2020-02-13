package v1alpha1

import (
	"log"
	"os"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func (c *AzurePipelinesPoolV1Alpha1Client) AzurePipelinesPool(namespace string) AzurePipelinesPoolInterface {
	return &AzurePipelinesPoolclient{
		client: c.RestClient,
		ns:     namespace,
	}
}

type AzurePipelinesPoolV1Alpha1Client struct {
	RestClient rest.Interface
}

type AzurePipelinesPoolInterface interface {
	Get(name string) (*AzurePipelinesPool, error)
	AddNewPodForCR(obj *AzurePipelinesPool, labels map[string]string) *v1.Pod
}

type AzurePipelinesPoolclient struct {
	client rest.Interface
	ns     string
}

func (c *AzurePipelinesPoolclient) Get(name string) (*AzurePipelinesPool, error) {
	log.Println("Came inside get method")
	result := &AzurePipelinesPool{}
	err := c.client.Get().
		Namespace(c.ns).Resource("azurepipelinespools").
		Name(name).Do().Into(result)
	return result, err
}

func (c *AzurePipelinesPoolclient) AddNewPodForCR(obj *AzurePipelinesPool, labels map[string]string) *v1.Pod {

	var spec *v1.PodSpec
	if IsTestingEnv() {
		spec = &v1.PodSpec{
			Containers: []v1.Container{
				{
					Name:  "vsts-agent",
					Image: "prebansa/myagent:v1",
				},
			},
		}
	} else {
		spec = FetchPodSpec(obj)
	}

	// append the RUNNING_ON environment variable
	if spec != nil && len(spec.Containers) > 0 {
		spec.Containers[0].Env = append(spec.Containers[0].Env, *GetRunningOnEnvironmentVariable())
	}

	// check if VolumeMounts is not present in the spec; then add the default one
	if spec != nil && len(spec.Containers) > 0 && spec.Containers[0].VolumeMounts == nil {
		spec.Containers[0].VolumeMounts = append(spec.Containers[0].VolumeMounts, *GetDefaultVolumeMount())
	}

	if spec != nil {
		dep := &v1.Pod{
			ObjectMeta: meta_v1.ObjectMeta{
				Labels:       labels,
				GenerateName: "azure-pipelines-agent-",
			},
			Spec: *spec,
		}
		if IsTestingEnv() {
			dep.Name = "TestAgentPod"
		}
		return dep
	}
	return nil
}

func FetchPodSpec(obj *AzurePipelinesPool) *v1.PodSpec {

	if obj.Spec.AgentPools != nil && len(obj.Spec.AgentPools) > 0 {
		// currently as demands are not supported so creating agentpod from agentspec being passed at first index
		return obj.Spec.AgentPools[0].PoolSpec
	}

	return nil
}

func GetDefaultVolumeMount() *v1.VolumeMount {

	return &v1.VolumeMount{
		Name:      "agent-creds",
		MountPath: "/azurepipelines/agent",
		ReadOnly:  true,
	}

}

func GetRunningOnEnvironmentVariable() *v1.EnvVar {
	return &v1.EnvVar{
		Name: "RUNNING_ON",
		ValueFrom: &v1.EnvVarSource{
			ConfigMapKeyRef: &v1.ConfigMapKeySelector{
				LocalObjectReference: v1.LocalObjectReference{Name: "kubernetes-config"},
				Key:                  "type",
			},
		},
	}
}

func IsTestingEnv() bool {
	testingMode := os.Getenv("IS_TESTENVIRONMENT")

	if testingMode == "true" {
		return true
	}
	return false
}
