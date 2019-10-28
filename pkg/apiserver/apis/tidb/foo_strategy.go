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

package tidb

import (
	"context"
	"sort"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

func (FooStrategy) PrepareForCreate(ctx context.Context, obj runtime.Object) {
	o := obj.(*Foo)
	if o.Spec.ConfigMapName == "" {
		o.Spec.ConfigMapName = "default"
	}
}

func (FooStrategy) Validate(ctx context.Context, obj runtime.Object) field.ErrorList {
	o := obj.(*Foo)
	errors := field.ErrorList{}
	specPath := field.NewPath("spec")
	if len(o.Spec.Replicas) < 1 {
		errors = append(errors, field.Invalid(specPath.Child("Replicas"), o.Spec.Replicas, "Replicas cannot be empty"))
	}
	return errors
}


func (FooStrategy) Canonicalize(obj runtime.Object) {
	o := obj.(*Foo)
	sort.Ints(o.Spec.Replicas)
}
