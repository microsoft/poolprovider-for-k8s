package v1alpha1

import (
	v1 "k8s.io/api/core/v1"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type PodConfig struct {
	meta_v1.TypeMeta   `json:",inline"`
	meta_v1.ObjectMeta `json:"metadata"`
	//Spec               *v1.PodSpec `json:"spec"`
	Spec PodConfigSpec `json:"spec"`
}
type PodConfigSpec struct {
	//Image string `json:"image"`
	AgentPools []AgentPoolSpec `json:"agentPools"`
	//Podspec v1.PodSpec `yaml:"podspec"`
}

/*type PodConfigSpec struct {
	Template corev1.PodTemplateSpec `json:"template"`
}*/

type AgentPoolSpec struct {
	PoolName string      `json:"name"`
	PoolSpec *v1.PodSpec `json:"spec"`
}

type PodConfigList struct {
	meta_v1.TypeMeta `json:",inline"`
	meta_v1.ListMeta `json:"metadata"`
	Items            []PodConfig `json:"items"`
}
