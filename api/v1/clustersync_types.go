/*
Copyright 2025 Jacob Philip.

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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ClusterSyncStatus defines the observed state of ClusterSync.
type ClusterSyncStatus struct {
	// LastSyncTime is the timestamp of the last successful sync
	LastSyncTime *metav1.Time `json:"lastSyncTime,omitempty"`
	// SyncStatus indicates the current sync status
	SyncStatus string `json:"syncStatus,omitempty"`
	// ErrorMessage contains any error message if sync failed
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// ClusterSync is the Schema for the clustersyncs API.
type ClusterSync struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Status ClusterSyncStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ClusterSyncList contains a list of ClusterSync.
type ClusterSyncList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ClusterSync `json:"items"`
}

// ReportVulnerabilitiesSpec defines the desired state of ReportVulnerabilities
type ReportVulnerabilitiesSpec struct {
	Data string `json:"data,omitempty"`
}

type ReportVulnerabilitiesStatus struct {
	// Add status fields if needed
}

// ReportVulnerabilities is the Schema for the reportvulnerabilities API
// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
type ReportVulnerabilities struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ReportVulnerabilitiesSpec   `json:"spec,omitempty"`
	Status ReportVulnerabilitiesStatus `json:"status,omitempty"`
}

// ReportVulnerabilitiesList contains a list of ReportVulnerabilities
// +kubebuilder:object:root=true
type ReportVulnerabilitiesList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ReportVulnerabilities `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ClusterSync{}, &ClusterSyncList{})
	SchemeBuilder.Register(&ReportVulnerabilities{}, &ReportVulnerabilitiesList{})
}
