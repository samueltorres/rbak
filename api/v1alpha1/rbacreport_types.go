/*
Copyright 2022.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1alpha1

import (
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Rules struct {
	Namespace string   `json:"namespace,omitempty"`
	APIGroups []string `json:"apiGroups,omitempty"`
	Resources []string `json:"resources,omitempty"`
	Verbs     []string `json:"verbs,omitempty"`
	Versions  []string `json:"versions,omitempty"`
}

// RbacReportSpec defines the desired state of RbacReport
type RbacReportSpec struct {
	Subject rbacv1.Subject `json:"subject,omitempty"`
	Rules   []Rules        `json:"rules,omitempty"`
}

// RbacReportStatus defines the observed state of RbacReport
type RbacReportStatus struct {
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status
//+kubebuilder:resource:scope=Cluster
//+kubebuilder:printcolumn:name="SubjectKind",type="string",JSONPath=".spec.subject.kind"
//+kubebuilder:printcolumn:name="SubjectName",type="string",JSONPath=".spec.subject.name"
//+kubebuilder:printcolumn:name="SubjectNamespace",type="string",JSONPath=".spec.subject.namespace"

// RbacReport is the Schema for the rbacreports API
type RbacReport struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RbacReportSpec   `json:"spec,omitempty"`
	Status RbacReportStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// RbacReportList contains a list of RbacReport
type RbacReportList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []RbacReport `json:"items"`
}

func init() {
	SchemeBuilder.Register(&RbacReport{}, &RbacReportList{})
}
