package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
)

// StuntDoublePolicySpec defines the desired state of StuntDoublePolicy
type StuntDoublePolicySpec struct {
	EnforcementMode string `json:"enforcementMode,omitempty"`
	Network         NetworkRules `json:"network,omitempty"`
}

type NetworkRules struct {
	Allow []string `json:"allow,omitempty"`
	Deny  []string `json:"deny,omitempty"`
}

// StuntDoublePolicyStatus defines the observed state of StuntDoublePolicy
type StuntDoublePolicyStatus struct {
	ActiveAgents int `json:"activeAgents,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// StuntDoublePolicy is the Schema for the stuntdoublepolicies API
type StuntDoublePolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   StuntDoublePolicySpec   `json:"spec,omitempty"`
	Status StuntDoublePolicyStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// StuntDoublePolicyList contains a list of StuntDoublePolicy
type StuntDoublePolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []StuntDoublePolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&StuntDoublePolicy{}, &StuntDoublePolicyList{})
}

// DeepCopyObject implements runtime.Object
func (in *StuntDoublePolicy) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(StuntDoublePolicy)
	in.DeepCopyInto(out)
	return out
}

func (in *StuntDoublePolicy) DeepCopyInto(out *StuntDoublePolicy) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	out.Spec = in.Spec
	out.Status = in.Status
}

// DeepCopyObject implements runtime.Object
func (in *StuntDoublePolicyList) DeepCopyObject() runtime.Object {
	if in == nil {
		return nil
	}
	out := new(StuntDoublePolicyList)
	in.DeepCopyInto(out)
	return out
}

func (in *StuntDoublePolicyList) DeepCopyInto(out *StuntDoublePolicyList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]StuntDoublePolicy, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}
