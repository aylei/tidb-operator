/*
Copyright 2017 The Kubernetes Authors.

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

// This file was autogenerated by apiregister-gen. Do not edit it manually!

package apis

import (
	"github.com/kubernetes-incubator/apiserver-builder-alpha/pkg/builders"
	"github.com/pingcap/tidb-operator/pkg/apiserver/apis/tidb"
	_ "github.com/pingcap/tidb-operator/pkg/apiserver/apis/tidb/install"
	tidbv1alpha1 "github.com/pingcap/tidb-operator/pkg/apiserver/apis/tidb/v1alpha1"
	tidbv1beta1 "github.com/pingcap/tidb-operator/pkg/apiserver/apis/tidb/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
)

var (
	localSchemeBuilder = runtime.SchemeBuilder{
		tidbv1alpha1.AddToScheme,
		tidbv1beta1.AddToScheme,
	}
	AddToScheme = localSchemeBuilder.AddToScheme
)

// GetAllApiBuilders returns all known APIGroupBuilders
// so they can be registered with the apiserver
func GetAllApiBuilders() []*builders.APIGroupBuilder {
	return []*builders.APIGroupBuilder{
		GetTidbAPIBuilder(),
	}
}

var tidbApiGroup = builders.NewApiGroupBuilder(
	"tidb.pingcap.com",
	"github.com/pingcap/tidb-operator/pkg/apiserver/apis/tidb").
	WithUnVersionedApi(tidb.ApiVersion).
	WithVersionedApis(
		tidbv1alpha1.ApiVersion,
		tidbv1beta1.ApiVersion,
	).
	WithRootScopedKinds()

func GetTidbAPIBuilder() *builders.APIGroupBuilder {
	return tidbApiGroup
}
