// +build !ignore_autogenerated

/*
Copyright The Kubernetes Authors.

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

// Code generated by conversion-gen. DO NOT EDIT.

package v1beta1

import (
	tidb "github.com/pingcap/tidb-operator/pkg/apiserver/apis/tidb"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Foo)(nil), (*tidb.Foo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_Foo_To_tidb_Foo(a.(*Foo), b.(*tidb.Foo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*tidb.Foo)(nil), (*Foo)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_tidb_Foo_To_v1beta1_Foo(a.(*tidb.Foo), b.(*Foo), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*FooList)(nil), (*tidb.FooList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_FooList_To_tidb_FooList(a.(*FooList), b.(*tidb.FooList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*tidb.FooList)(nil), (*FooList)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_tidb_FooList_To_v1beta1_FooList(a.(*tidb.FooList), b.(*FooList), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*FooSpec)(nil), (*tidb.FooSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_FooSpec_To_tidb_FooSpec(a.(*FooSpec), b.(*tidb.FooSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*tidb.FooSpec)(nil), (*FooSpec)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_tidb_FooSpec_To_v1beta1_FooSpec(a.(*tidb.FooSpec), b.(*FooSpec), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*FooStatus)(nil), (*tidb.FooStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1beta1_FooStatus_To_tidb_FooStatus(a.(*FooStatus), b.(*tidb.FooStatus), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*tidb.FooStatus)(nil), (*FooStatus)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_tidb_FooStatus_To_v1beta1_FooStatus(a.(*tidb.FooStatus), b.(*FooStatus), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1beta1_Foo_To_tidb_Foo(in *Foo, out *tidb.Foo, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_v1beta1_FooSpec_To_tidb_FooSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_v1beta1_FooStatus_To_tidb_FooStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1beta1_Foo_To_tidb_Foo is an autogenerated conversion function.
func Convert_v1beta1_Foo_To_tidb_Foo(in *Foo, out *tidb.Foo, s conversion.Scope) error {
	return autoConvert_v1beta1_Foo_To_tidb_Foo(in, out, s)
}

func autoConvert_tidb_Foo_To_v1beta1_Foo(in *tidb.Foo, out *Foo, s conversion.Scope) error {
	out.ObjectMeta = in.ObjectMeta
	if err := Convert_tidb_FooSpec_To_v1beta1_FooSpec(&in.Spec, &out.Spec, s); err != nil {
		return err
	}
	if err := Convert_tidb_FooStatus_To_v1beta1_FooStatus(&in.Status, &out.Status, s); err != nil {
		return err
	}
	return nil
}

// Convert_tidb_Foo_To_v1beta1_Foo is an autogenerated conversion function.
func Convert_tidb_Foo_To_v1beta1_Foo(in *tidb.Foo, out *Foo, s conversion.Scope) error {
	return autoConvert_tidb_Foo_To_v1beta1_Foo(in, out, s)
}

func autoConvert_v1beta1_FooList_To_tidb_FooList(in *FooList, out *tidb.FooList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]tidb.Foo, len(*in))
		for i := range *in {
			if err := Convert_v1beta1_Foo_To_tidb_Foo(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_v1beta1_FooList_To_tidb_FooList is an autogenerated conversion function.
func Convert_v1beta1_FooList_To_tidb_FooList(in *FooList, out *tidb.FooList, s conversion.Scope) error {
	return autoConvert_v1beta1_FooList_To_tidb_FooList(in, out, s)
}

func autoConvert_tidb_FooList_To_v1beta1_FooList(in *tidb.FooList, out *FooList, s conversion.Scope) error {
	out.ListMeta = in.ListMeta
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]Foo, len(*in))
		for i := range *in {
			if err := Convert_tidb_Foo_To_v1beta1_Foo(&(*in)[i], &(*out)[i], s); err != nil {
				return err
			}
		}
	} else {
		out.Items = nil
	}
	return nil
}

// Convert_tidb_FooList_To_v1beta1_FooList is an autogenerated conversion function.
func Convert_tidb_FooList_To_v1beta1_FooList(in *tidb.FooList, out *FooList, s conversion.Scope) error {
	return autoConvert_tidb_FooList_To_v1beta1_FooList(in, out, s)
}

func autoConvert_v1beta1_FooSpec_To_tidb_FooSpec(in *FooSpec, out *tidb.FooSpec, s conversion.Scope) error {
	// WARNING: in.Replicas requires manual conversion: inconvertible types ([][]int vs int)
	// WARNING: in.ConfigMapName requires manual conversion: does not exist in peer-type
	return nil
}

func autoConvert_tidb_FooSpec_To_v1beta1_FooSpec(in *tidb.FooSpec, out *FooSpec, s conversion.Scope) error {
	// WARNING: in.Replicas requires manual conversion: inconvertible types (int vs [][]int)
	return nil
}

func autoConvert_v1beta1_FooStatus_To_tidb_FooStatus(in *FooStatus, out *tidb.FooStatus, s conversion.Scope) error {
	out.CurrentReplicas = in.CurrentReplicas
	return nil
}

// Convert_v1beta1_FooStatus_To_tidb_FooStatus is an autogenerated conversion function.
func Convert_v1beta1_FooStatus_To_tidb_FooStatus(in *FooStatus, out *tidb.FooStatus, s conversion.Scope) error {
	return autoConvert_v1beta1_FooStatus_To_tidb_FooStatus(in, out, s)
}

func autoConvert_tidb_FooStatus_To_v1beta1_FooStatus(in *tidb.FooStatus, out *FooStatus, s conversion.Scope) error {
	out.CurrentReplicas = in.CurrentReplicas
	return nil
}

// Convert_tidb_FooStatus_To_v1beta1_FooStatus is an autogenerated conversion function.
func Convert_tidb_FooStatus_To_v1beta1_FooStatus(in *tidb.FooStatus, out *FooStatus, s conversion.Scope) error {
	return autoConvert_tidb_FooStatus_To_v1beta1_FooStatus(in, out, s)
}
