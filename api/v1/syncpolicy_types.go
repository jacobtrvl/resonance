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

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// SyncPolicySpec defines the desired state of SyncPolicy.
type SyncPolicySpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of SyncPolicy. Edit syncpolicy_types.go to remove/update
	Foo string `json:"foo,omitempty"`

	// RemoteClusterConfig contains the configuration for the remote cluster
	RemoteClusterConfig RemoteClusterConfig `json:"remoteClusterConfig"`
}

// RemoteClusterConfig defines the configuration for the remote cluster
type RemoteClusterConfig struct {
	// KubeconfigSecretName is the name of the secret containing the kubeconfig for the remote cluster
	KubeconfigSecretName string `json:"kubeconfigSecretName"`
	// KubeconfigSecretNamespace is the namespace where the kubeconfig secret is stored
	KubeconfigSecretNamespace string `json:"kubeconfigSecretNamespace"`
}

// SyncPolicyStatus defines the observed state of SyncPolicy.
type SyncPolicyStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// LastSyncTime is the timestamp of the last successful sync
	LastSyncTime *metav1.Time `json:"lastSyncTime,omitempty"`
	// SyncStatus indicates the current sync status
	SyncStatus string `json:"syncStatus,omitempty"`
	// ErrorMessage contains any error message if sync failed
	ErrorMessage string `json:"errorMessage,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status

// SyncPolicy is the Schema for the syncpolicies API.
type SyncPolicy struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   SyncPolicySpec   `json:"spec,omitempty"`
	Status SyncPolicyStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// SyncPolicyList contains a list of SyncPolicy.
type SyncPolicyList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SyncPolicy `json:"items"`
}

func init() {
	SchemeBuilder.Register(&SyncPolicy{}, &SyncPolicyList{})
}
