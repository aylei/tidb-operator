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

package defaulting

import (
	"github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
	corev1 "k8s.io/api/core/v1"
)

const (
	defaultTiDBImage   = "pingcap/tidb"
	defaultTiKVImage   = "pingcap/tikv"
	defaultPDImage     = "pingcap/pd"
	defaultBinlogImage = "pingcap/tidb-binlog"
)

func SetTidbClusterDefault(tc *v1alpha1.TidbCluster) {
	if tc.Spec.TiDB.BaseImage == "" {
		tc.Spec.TiDB.BaseImage = defaultTiDBImage
	}
	if tc.Spec.TiKV.BaseImage == "" {
		tc.Spec.TiKV.BaseImage = defaultTiKVImage
	}
	if tc.Spec.PD.BaseImage == "" {
		tc.Spec.PD.BaseImage = defaultPDImage
	}
	if tc.Spec.Pump != nil && tc.Spec.Pump.BaseImage == "" {
		tc.Spec.Pump.BaseImage = defaultBinlogImage
	}
	if tc.Spec.TiDB.Config == nil {
		tc.Spec.TiDB.Config = &v1alpha1.TiDBConfig{}
	}
	if tc.Spec.TiKV.Config == nil {
		tc.Spec.TiKV.Config = &v1alpha1.TiKVConfig{}
	}
	if tc.Spec.PD.Config == nil {
		tc.Spec.PD.Config = &v1alpha1.PDConfig{}
	}
	if string(tc.Spec.ImagePullPolicy) == "" {
		tc.Spec.ImagePullPolicy = corev1.PullIfNotPresent
	}
}
