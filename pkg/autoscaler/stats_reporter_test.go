/*
Copyright 2018 Google Inc. All Rights Reserved.
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

package autoscaler

import (
	"testing"

	"go.opencensus.io/stats/view"
)

func TestReporter_Report(t *testing.T) {
	r := &Reporter{}
	if err := r.Report(DesiredPodCountM, 10); err == nil {
		t.Error("Reporter.Report() expected an error for Report call before init. Got success.")
	}

	r, _ = NewStatsReporter("testns", "testconfig", "testrev")
	wantTags := map[string]string{
		"configuration_namespace": "testns",
		"configuration":           "testconfig",
		"revision":                "testrev",
	}

	// Send statistics only once and observe the results
	expectSuccess(t, func() error { return r.Report(DesiredPodCountM, 10) })
	expectSuccess(t, func() error { return r.Report(RequestedPodCountM, 7) })
	expectSuccess(t, func() error { return r.Report(ActualPodCountM, 5) })
	expectSuccess(t, func() error { return r.Report(PanicM, 0) })
	checkData(t, "desired_pod_count", wantTags, 10)
	checkData(t, "requested_pod_count", wantTags, 7)
	checkData(t, "actual_pod_count", wantTags, 5)
	checkData(t, "panic_mode", wantTags, 0)

	// All the stats are gauges - record multiple entries for one stat - last one should stick
	expectSuccess(t, func() error { return r.Report(DesiredPodCountM, 1) })
	expectSuccess(t, func() error { return r.Report(DesiredPodCountM, 2) })
	expectSuccess(t, func() error { return r.Report(DesiredPodCountM, 3) })
	checkData(t, "desired_pod_count", wantTags, 3)

	expectSuccess(t, func() error { return r.Report(RequestedPodCountM, 4) })
	expectSuccess(t, func() error { return r.Report(RequestedPodCountM, 5) })
	expectSuccess(t, func() error { return r.Report(RequestedPodCountM, 6) })
	checkData(t, "requested_pod_count", wantTags, 6)

	expectSuccess(t, func() error { return r.Report(ActualPodCountM, 7) })
	expectSuccess(t, func() error { return r.Report(ActualPodCountM, 8) })
	expectSuccess(t, func() error { return r.Report(ActualPodCountM, 9) })
	checkData(t, "actual_pod_count", wantTags, 9)

	expectSuccess(t, func() error { return r.Report(PanicM, 1) })
	expectSuccess(t, func() error { return r.Report(PanicM, 0) })
	expectSuccess(t, func() error { return r.Report(PanicM, 1) })
	checkData(t, "panic_mode", wantTags, 1)

	expectSuccess(t, func() error { return r.Report(PanicM, 0) })
	checkData(t, "panic_mode", wantTags, 0)
}

func expectSuccess(t *testing.T, f func() error) {
	if err := f(); err != nil {
		t.Errorf("Reporter.Report() expected success but got error %v", err)
	}
}

func checkData(t *testing.T, name string, wantTags map[string]string, wantValue int) {
	if d, err := view.RetrieveData(name); err != nil {
		t.Errorf("Reporter.Report() error = %v, wantErr %v", err, false)
	} else {
		if len(d) != 1 {
			t.Errorf("Reporter.Report() len(d) %v, want %v", len(d), 1)
		}
		for _, got := range d[0].Tags {
			if want, ok := wantTags[got.Key.Name()]; !ok {
				t.Errorf("Reporter.Report() got an extra tag %v: %v", got.Key.Name(), got.Value)
			} else {
				if got.Value != want {
					t.Errorf("Reporter.Report() expected a different tag value. key:%v, got: %v, want: %v", got.Key.Name(), got.Value, want)
				}
			}
		}

		if s, ok := d[0].Data.(*view.LastValueData); !ok {
			t.Error("Reporter.Report() expected a LastValueData type")
		} else {
			if s.Value != (float64)(wantValue) {
				t.Errorf("Reporter.Report() expected %v got %v. metric: %v", s.Value, (float64)(wantValue), name)
			}
		}
	}
}
