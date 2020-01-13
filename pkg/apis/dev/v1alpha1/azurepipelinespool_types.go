package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        corev1 "k8s.io/api/core/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AzurePipelinesPoolSpec defines the desired state of AzurePipelinesPool
type AzurePipelinesPoolSpec struct {
	ControllerName string `json:"controllerImage"`
        BuildkitReplicaCount int32 `json:"buildkitReplicas"`
	AgentPools []AgentPoolSpec `json:"agentPools"`
	Initialized bool  `json:"initialized"`
}

type AgentPoolSpec struct {
	PoolName string      `json:"name"`
	PoolSpec *corev1.PodSpec `json:"spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzurePipelinesPool is the Schema for the azurepipelinespools API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=azurepipelinespools,scope=Namespaced
type AzurePipelinesPool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzurePipelinesPoolSpec   `json:"spec,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AzurePipelinesPoolList contains a list of AzurePipelinesPool
type AzurePipelinesPoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzurePipelinesPool `json:"items"`
}

func init() {
	SchemeBuilder.Register(&AzurePipelinesPool{}, &AzurePipelinesPoolList{})
}
