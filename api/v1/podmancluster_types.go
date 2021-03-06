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

package v1

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// ClusterFinalizer allows PodmanClusterReconciler to clean up resources associated with PodmanCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "podmancluster.infrastructure.cluster.x-k8s.io"
)

// PodmanClusterSpec defines the desired state of PodmanCluster
type PodmanClusterSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file

	// Foo is an example field of PodmanCluster. Edit podmancluster_types.go to remove/update
	Foo string `json:"foo,omitempty"`
}

// PodmanClusterStatus defines the observed state of PodmanCluster
type PodmanClusterStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Conditions defines current service state of the ByoCluster.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

//+kubebuilder:object:root=true
//+kubebuilder:subresource:status

// PodmanCluster is the Schema for the podmanclusters API
type PodmanCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   PodmanClusterSpec   `json:"spec,omitempty"`
	Status PodmanClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// PodmanClusterList contains a list of PodmanCluster
type PodmanClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []PodmanCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&PodmanCluster{}, &PodmanClusterList{})
}

// GetConditions gets the condition for the ByoCluster status
func (podCluster *PodmanCluster) GetConditions() clusterv1.Conditions {
	return podCluster.Status.Conditions
}

// SetConditions sets the conditions for the ByoCluster status
func (podCluster *PodmanCluster) SetConditions(conditions clusterv1.Conditions) {
	podCluster.Status.Conditions = conditions
}
