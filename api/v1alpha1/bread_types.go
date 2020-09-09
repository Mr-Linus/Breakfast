/*

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
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// BreadSpec defines the desired state of Bread
type BreadSpec struct {
	// INSERT ADDITIONAL SPEC FIELDS - desired state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	// Foo is an example field of Bread. Edit Bread_types.go to remove/update
	Scv       SCVSpec       `json:"scv,omitempty"`
	Framework FrameworkSpec `json:"framework,omitempty"`
	Task      TaskSpec      `json:"task,omitempty"`
}

// BreadStatus defines the observed state of Bread
type BreadStatus struct {
	// INSERT ADDITIONAL STATUS FIELD - define observed state of cluster
	// Important: Run "make" to regenerate code after modifying this file
	Phase             v1.PodPhase          `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PodPhase"`
	ContainerStatuses []v1.ContainerStatus `json:"containerStatuses,omitempty" protobuf:"bytes,8,rep,name=containerStatuses"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// Bread is the Schema for the breads API
type Bread struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              BreadSpec   `json:"spec,omitempty"`
	Status            BreadStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// BreadList contains a list of Bread
type BreadList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []Bread `json:"items"`
}

// +kubebuilder:object:generate=true
type SCVSpec struct {
	Gpu      string `json:"gpu,omitempty"`
	Memory   string `json:"memory,omitempty"`
	Clock    string `json:"clock,omitempty"`
	Priority string `json:"priority,omitempty"`
}

type FrameworkSpec struct {
	Name    string `json:"name,omitempty"`
	Version string `json:"version,omitempty"`
}

type TaskSpec struct {
	Type    string `json:"type,omitempty"`
	Command string `json:"command,omitempty"`
}

func init() {
	SchemeBuilder.Register(&Bread{}, &BreadList{})
}
