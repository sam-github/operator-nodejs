package v1alpha1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for
// the fields to be serialized.

// NodejsDiagnosticReportSpec defines the desired state of NodejsDiagnosticReport
type NodejsDiagnosticReportSpec struct {
	// Name of the pod that should write a diagnostic report.
	PodName string `json:"podName"`
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after
	// modifying this file
	// Add custom validation using kubebuilder tags:
	// https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// NodejsDiagnosticReportStatus defines the observed state of NodejsDiagnosticReport
type NodejsDiagnosticReportStatus struct {
	// Result of triggering the report
	Result string `json:"result"`
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "operator-sdk generate k8s" to regenerate code after
	// modifying this file
	// Add custom validation using kubebuilder tags:
	// https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodejsDiagnosticReport is the Schema for the nodejsdiagnosticreports API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=nodejsdiagnosticreports,scope=Namespaced
type NodejsDiagnosticReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   NodejsDiagnosticReportSpec   `json:"spec,omitempty"`
	Status NodejsDiagnosticReportStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodejsDiagnosticReportList contains a list of NodejsDiagnosticReport
type NodejsDiagnosticReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []NodejsDiagnosticReport `json:"items"`
}

func init() {
	SchemeBuilder.Register(&NodejsDiagnosticReport{}, &NodejsDiagnosticReportList{})
}
