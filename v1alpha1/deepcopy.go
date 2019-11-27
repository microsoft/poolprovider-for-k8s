package v1alpha1

import "k8s.io/apimachinery/pkg/runtime"

// DeepCopyInto copies all properties of this object into another object of the
// same type that is provided as a pointer.
func (in *PodConfig) DeepCopyInto(out *PodConfig) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
	/*out.Spec = PodConfigSpec{
		//Podspec: in.Spec.Podspec,
		//Image: in.Spec.Image,
		if in.AgentPools != nil {
			out.AgentPools = make([]AgentPoolSpec, len(in.AgentPools))
			for i := range in.AgentPools {
				out.AgentPools[i].PoolName: in.AgentPools[i].PoolName,
				out.AgentPools[i].PoolSpec: in.AgentPools[i].PoolSpec,
			}
		}
	}*/
	out.Spec = in.Spec
	//Image: in.Spec.Image,
}

// DeepCopyObject returns a generically typed copy of an object
func (in *PodConfig) DeepCopyObject() runtime.Object {
	out := PodConfig{}
	in.DeepCopyInto(&out)

	return &out
}

// DeepCopyObject returns a generically typed copy of an object
func (in *PodConfigList) DeepCopyObject() runtime.Object {
	out := PodConfigList{}
	out.TypeMeta = in.TypeMeta
	out.ListMeta = in.ListMeta

	if in.Items != nil {
		out.Items = make([]PodConfig, len(in.Items))
		for i := range in.Items {
			in.Items[i].DeepCopyInto(&out.Items[i])
		}
	}

	return &out
}
