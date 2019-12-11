package v1alpha1

import (
	"log"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func (c *AzurePipelinePoolV1Alpha1Client) AzurePipelinePool(namespace string) AzurePipelinePoolInterface {
	return &AzurePipelinePoolclient{
		client: c.restClient,
		ns:     namespace,
	}
}

type AzurePipelinePoolV1Alpha1Client struct {
	restClient rest.Interface
}

type AzurePipelinePoolInterface interface {
	Get(name string) (*AzurePipelinePool, error)
	AddNewPodForCR(obj *AzurePipelinePool, agentId string, labels map[string]string, poolName string) *v1.Pod
}

type AzurePipelinePoolclient struct {
	client rest.Interface
	ns     string
}

func (c *AzurePipelinePoolclient) Get(name string) (*AzurePipelinePool, error) {
	log.Println("Came insidde get method")
	result := &AzurePipelinePool{}
	err := c.client.Get().
		Namespace(c.ns).Resource("azurepipelinepools").
		Name(name).Do().Into(result)
	return result, err
}

func (c *AzurePipelinePoolclient) AddNewPodForCR(obj *AzurePipelinePool, agentId string, labels map[string]string, poolname string) *v1.Pod {

	spec := FetchPodSpec(obj, poolname)

	dep := &v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Labels:       labels,
			GenerateName: "vsts-agent-",
		},
		Spec: *spec,
	}

	return dep
}

func FetchPodSpec(obj *AzurePipelinePool, poolname string) *v1.PodSpec {
	var p1 v1.PodSpec

	if obj.Spec.AgentPools != nil {
		for i := range obj.Spec.AgentPools {
			if obj.Spec.AgentPools[i].PoolName == poolname {
				return obj.Spec.AgentPools[i].PoolSpec
			}
		}
	}

	return &p1
}
