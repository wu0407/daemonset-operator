package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
        appsv1 "k8s.io/api/apps/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// XdaemonsetSpec defines the desired state of Xdaemonset
type XdaemonsetSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
        appsv1.DaemonSetSpec `json:",inline"`
}

// XdaemonsetStatus defines the observed state of Xdaemonset
type XdaemonsetStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Xdaemonset is the Schema for the xdaemonsets API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=xdaemonsets,scope=Namespaced
type Xdaemonset struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   XdaemonsetSpec   `json:"spec,omitempty"`
	Status XdaemonsetStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// XdaemonsetList contains a list of Xdaemonset
type XdaemonsetList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Xdaemonset `json:"items"`
}

func init() {
	SchemeBuilder.Register(&Xdaemonset{}, &XdaemonsetList{})
}
