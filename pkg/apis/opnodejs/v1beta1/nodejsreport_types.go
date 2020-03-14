package v1beta1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for
// the fields to be serialized.

// NodejsReportSpec defines the desired state of NodejsReport
type NodejsReportSpec struct {
	// Name of the pod that should write a report.
	PodName string `json:"podName"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after
	// modifying this file
	// Add custom validation using kubebuilder tags:
	// https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// NodejsReportStatus defines the observed state of NodejsReport
type NodejsReportStatus struct {
	// Result of triggering the report
	Result string `json:"result"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after
	// modifying this file
	// Add custom validation using kubebuilder tags:
	// https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodejsReport is the Schema for the nodejsreports API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nodejsreports,scope=Namespaced
type NodejsReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodejsReportSpec   `json:"spec,omitempty"`
	Status NodejsReportStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodejsReportList contains a list of NodejsReport
type NodejsReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodejsReport `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodejsReport{}, &NodejsReportList{})
}
