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
	"sort"

	batch "k8s.io/api/batch/v1"
	"k8s.io/api/core/v1"
)

func IsJobFinished(j *batch.Job) bool {
	for _, c := range j.Status.Conditions {
		if (c.Type == batch.JobComplete || c.Type == batch.JobFailed) && c.Status == v1.ConditionTrue {
			return true
		}
	}
	return false
}

// podByFailedTime used to sort pod by LastTransitionTime Conditions
type podByFailedTime []*v1.Pod

func (p podByFailedTime) Len() int {
	return len(p)
}
func (p podByFailedTime) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}
func (p podByFailedTime) Less(i, j int) bool {
	ci := getTrueCondition(p[i])
	cj := getTrueCondition(p[j])

	if cj != nil && ci != nil {
		if cj.LastProbeTime.After(ci.LastProbeTime.Time) {
			return true
		}
	}

	return false
}

// Sort used to sort the element of the podByFailedTime slice
func (p *podByFailedTime) Sort() {
	sort.Sort(p)
}

func getTrueCondition(pod *v1.Pod) *v1.PodCondition {
	for _, c := range pod.Status.Conditions {
		if c.Status == v1.ConditionTrue {
			return &c
		}
	}
	return nil
}
