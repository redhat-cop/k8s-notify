package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// EventSubscriptionSpec defines the desired state of EventSubscription
// +k8s:openapi-gen=true
type EventSubscriptionSpec struct {
	MatchObject  corev1.ObjectReference `json:"matchObject,omitempty"`
	MatchMessage string                 `json:"matchMessage,omitempty"`
	MatchReason  string                 `json:"matchReason,omitempty"`
	MatchType    string                 `json:"matchType,omitempty"`
	Notifier     string                 `json:"notifier"`
}

// EventSubscriptionStatus defines the observed state of EventSubscription
// +k8s:openapi-gen=true
type EventSubscriptionStatus struct {
	Phase string `json:"spec"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
	// Add custom validation using kubebuilder tags: https://book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// EventSubscription is the Schema for the eventsubscriptions API
// +k8s:openapi-gen=true
type EventSubscription struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   EventSubscriptionSpec   `json:"spec,omitempty"`
	Status EventSubscriptionStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventSubscriptionList contains a list of EventSubscription
type EventSubscriptionList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []EventSubscription `json:"items"`
}

func init() {
	SchemeBuilder.Register(&EventSubscription{}, &EventSubscriptionList{})
}
