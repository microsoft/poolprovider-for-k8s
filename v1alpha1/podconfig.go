package v1alpha1

import (
	"log"

	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/rest"
)

func (c *PodConfigV1Alpha1Client) PodConfigs(namespace string) PodConfigInterface {
	return &PodConfigclient{
		client: c.restClient,
		ns:     namespace,
	}
}

type PodConfigV1Alpha1Client struct {
	restClient rest.Interface
}

type PodConfigInterface interface {
	Create(obj *PodConfig) (*PodConfig, error)
	Update(obj *PodConfig) (*PodConfig, error)
	Delete(name string, options *meta_v1.DeleteOptions) error
	Get(name string) (*PodConfig, error)
	AddNewPodForCR(obj *PodConfig, agentId string, labels map[string]string, poolName string) *v1.Pod
}

type PodConfigclient struct {
	client rest.Interface
	ns     string
}

func (c *PodConfigclient) Create(obj *PodConfig) (*PodConfig, error) {
	result := &PodConfig{}
	err := c.client.Post().
		Namespace(c.ns).Resource("podconfigs").
		Body(obj).Do().Into(result)
	return result, err
}

func (c *PodConfigclient) Update(obj *PodConfig) (*PodConfig, error) {
	result := &PodConfig{}
	err := c.client.Put().
		Namespace(c.ns).Resource("podconfigs").
		Body(obj).Do().Into(result)
	return result, err
}

func (c *PodConfigclient) Delete(name string, options *meta_v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).Resource("podconfigs").
		Name(name).Body(options).Do().
		Error()
}

func (c *PodConfigclient) Get(name string) (*PodConfig, error) {
	log.Println("Came insidde get method")
	result := &PodConfig{}
	err := c.client.Get().
		Namespace(c.ns).Resource("podconfigs").
		Name(name).Do().Into(result)
	return result, err
}

func (c *PodConfigclient) AddNewPodForCR(obj *PodConfig, agentId string, labels map[string]string, poolname string) *v1.Pod {

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

func FetchPodSpec(obj *PodConfig, poolname string) *v1.PodSpec {
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
