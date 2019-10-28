// Copyright 2019. PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package v1beta1

import(
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Foo
// +k8s:openapi-gen=true
// +resource:path=foos,strategy=FooStrategy
type Foo struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   FooSpec   `json:"spec,omitempty"`
	Status FooStatus `json:"status,omitempty"`
}

// FooSpec defines the desired state of Foo
type FooSpec struct {
	Replicas      []int `json:"replicas,omitempty"`
	ConfigMapName string `json:"string,omitempty"`
}

// FooStatus defines the observed state of Foo
type FooStatus struct {
	CurrentReplicas int `json:"currentReplicas,omitempty"`
}

// DefaultingFunction sets default Foo field values
func (FooSchemeFns) DefaultingFunction(o interface{}) {
	obj := o.(*Foo)
	// set default field values here
	if len(obj.Spec.Replicas) == 0 {
		obj.Spec.Replicas = []int{1}
	}
}