package v1alpha1

import (
	"log"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func (c *AzurePipelinesPoolV1Alpha1Client) AzurePipelinesPool(namespace string) AzurePipelinesPoolInterface {
	return &AzurePipelinesPoolclient{
		client: c.restClient,
		ns:     namespace,
	}
}

type AzurePipelinesPoolV1Alpha1Client struct {
	restClient rest.Interface
}

type AzurePipelinesPoolInterface interface {
	Get(name string) (*AzurePipelinesPool, error)
	AddNewPodForCR(obj *AzurePipelinesPool, agentId string, labels map[string]string, poolName string) *v1.Pod
}

type AzurePipelinesPoolclient struct {
	client rest.Interface
	ns     string
}

func (c *AzurePipelinesPoolclient) Get(name string) (*AzurePipelinesPool, error) {
	log.Println("Came insidde get method")
	result := &AzurePipelinesPool{}
	err := c.client.Get().
		Namespace(c.ns).Resource("azurepipelinespool").
		Name(name).Do().Into(result)
	return result, err
}

func (c *AzurePipelinesPoolclient) AddNewPodForCR(obj *AzurePipelinesPool, agentId string, labels map[string]string, poolname string) *v1.Pod {

	spec := FetchPodSpec(obj, poolname)

	dep := &v1.Pod{
		ObjectMeta: meta_v1.ObjectMeta{
			Labels:       labels,
			GenerateName: "azure-pipelines-agent-",
		},
		Spec: *spec,
	}

	return dep
}

func FetchPodSpec(obj *AzurePipelinesPool, poolname string) *v1.PodSpec {

	if obj.Spec.AgentPools != nil {
		for i := range obj.Spec.AgentPools {
			if obj.Spec.AgentPools[i].PoolName == poolname {
				return obj.Spec.AgentPools[i].PoolSpec
			}
		}
	}

	return nil
}
