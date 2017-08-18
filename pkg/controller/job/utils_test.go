/*
Copyright 2016 The Kubernetes Authors.

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

package job

import (
	"testing"
	"time"

	batch "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestIsJobFinished(t *testing.T) {
	job := &batch.Job{
		Status: batch.JobStatus{
			Conditions: []batch.JobCondition{{
				Type:   batch.JobComplete,
				Status: v1.ConditionTrue,
			}},
		},
	}

	if !IsJobFinished(job) {
		t.Error("Job was expected to be finished")
	}

	job.Status.Conditions[0].Status = v1.ConditionFalse
	if IsJobFinished(job) {
		t.Error("Job was not expected to be finished")
	}

	job.Status.Conditions[0].Status = v1.ConditionUnknown
	if IsJobFinished(job) {
		t.Error("Job was not expected to be finished")
	}
}

func Test_podByFailedTime_Sort(t *testing.T) {

	pod1 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "Pod1"},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{
				{
					Status:        v1.ConditionTrue,
					LastProbeTime: metav1.Time{Time: time.Now()},
				},
			},
		},
	}

	pod2 := &v1.Pod{
		ObjectMeta: metav1.ObjectMeta{Name: "Pod2"},
		Status: v1.PodStatus{
			Conditions: []v1.PodCondition{
				{
					Status:        v1.ConditionTrue,
					LastProbeTime: metav1.Time{Time: time.Now().Add(time.Duration(2 * time.Hour))},
				},
			},
		},
	}

	resultOnlyPod1 := &podByFailedTime{pod1}

	result2PodsSorted := &podByFailedTime{pod1, pod2}

	tests := []struct {
		name   string
		p      *podByFailedTime
		result *podByFailedTime
	}{
		{
			"empty Slice",
			&podByFailedTime{},
			&podByFailedTime{},
		},
		{
			"one Pod",
			&podByFailedTime{pod1},
			resultOnlyPod1,
		},
		{
			"two Pods already sorted",
			&podByFailedTime{pod1, pod2},
			result2PodsSorted,
		},
		{
			"two Pods not sorted",
			&podByFailedTime{pod2, pod1},
			result2PodsSorted,
		},
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tt.p.Sort()
			for i, pod := range *tt.p {
				pods := *(tt.result)
				attendedPod := pods[i]
				if pod.Name != attendedPod.Name {
					t.Errorf("Slice not sorted properly pod.Name: %s after attended pod.Name: %s", pod.Name, attendedPod.Name)
				}
			}
		})
	}
}
