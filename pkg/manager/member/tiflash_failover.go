// Copyright 2020 PingCAP, Inc.
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

package member

import (
	"fmt"
	"time"

	"github.com/pingcap/tidb-operator/pkg/apis/pingcap/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/record"
	"k8s.io/klog"
)

type tiflashFailover struct {
	tiflashFailoverPeriod time.Duration
	recorder              record.EventRecorder
}

// NewTiFlashFailover returns a tiflash Failover
func NewTiFlashFailover(tiflashFailoverPeriod time.Duration, recorder record.EventRecorder) Failover {
	return &tiflashFailover{tiflashFailoverPeriod, recorder}
}

func (tff *tiflashFailover) Failover(tc *v1alpha1.TidbCluster) error {
	ns := tc.GetNamespace()
	tcName := tc.GetName()

	for storeID, store := range tc.Status.TiFlash.Stores {
		podName := store.PodName
		if store.LastTransitionTime.IsZero() {
			continue
		}
		deadline := store.LastTransitionTime.Add(tff.tiflashFailoverPeriod)
		exist := false
		for _, failureStore := range tc.Status.TiFlash.FailureStores {
			if failureStore.PodName == podName {
				exist = true
				break
			}
		}
		if store.State == v1alpha1.TiKVStateDown && time.Now().After(deadline) && !exist {
			if tc.Status.TiFlash.FailureStores == nil {
				tc.Status.TiFlash.FailureStores = map[string]v1alpha1.TiKVFailureStore{}
			}
			if tc.Spec.TiFlash.MaxFailoverCount != nil && *tc.Spec.TiFlash.MaxFailoverCount > 0 {
				maxFailoverCount := *tc.Spec.TiFlash.MaxFailoverCount
				if len(tc.Status.TiFlash.FailureStores) >= int(maxFailoverCount) {
					klog.Warningf("%s/%s TiFlash failure stores count reached the limit: %d", ns, tcName, tc.Spec.TiFlash.MaxFailoverCount)
					return nil
				}
				tc.Status.TiFlash.FailureStores[storeID] = v1alpha1.TiKVFailureStore{
					PodName:   podName,
					StoreID:   store.ID,
					CreatedAt: metav1.Now(),
				}
				msg := fmt.Sprintf("store [%s] is Down", store.ID)
				tff.recorder.Event(tc, corev1.EventTypeWarning, unHealthEventReason, fmt.Sprintf(unHealthEventMsgPattern, "tiflash", podName, msg))
			}
		}
	}
	return nil
}

func (tff *tiflashFailover) Recover(_ *v1alpha1.TidbCluster) {
	// Do nothing now
}

type fakeTiFlashFailover struct{}

// NewFakeTiFlashFailover returns a fake Failover
func NewFakeTiFlashFailover() Failover {
	return &fakeTiFlashFailover{}
}

func (ftff *fakeTiFlashFailover) Failover(_ *v1alpha1.TidbCluster) error {
	return nil
}

func (ftff *fakeTiFlashFailover) Recover(_ *v1alpha1.TidbCluster) {
	return
}
