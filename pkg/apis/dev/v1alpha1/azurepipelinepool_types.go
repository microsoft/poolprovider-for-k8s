package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        corev1 "k8s.io/api/core/v1"
)

// AzurePipelinePoolSpec defines the desired state of AzurePipelinePool
// +k8s:openapi-gen=true
type AzurePipelinePoolSpec struct {
        ControllerName string `json:"controllerImage"`
        BuildkitReplicaCount int32 `json:"buildkitReplicaCount"`
	AgentPools []AgentPoolSpec `json:"agentPools"`
        Initialized bool  `json:"initialized"`
}

type AgentPoolSpec struct {
	PoolName string      `json:"name"`
	PoolSpec *corev1.PodSpec `json:"spec"`
}


// AzurePipelinePoolStatus defines the observed state of AzurePipelinePool
// +k8s:openapi-gen=true
type AzurePipelinePoolStatus struct {
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzurePipelinePool is the Schema for the azurepipelinepools API
// +k8s:openapi-gen=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=azurepipelinepools,scope=Namespaced
type AzurePipelinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzurePipelinePoolSpec   `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzurePipelinePoolList contains a list of AzurePipelinePool
type AzurePipelinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzurePipelinePool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzurePipelinePool{}, &AzurePipelinePoolList{})
}
