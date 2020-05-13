// Copyright 2020 Apache Bench Operator Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// 	 http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package v1alpha1

import (
	batchv1 "k8s.io/api/batch/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// NOTE: json tags are required. Any new fields you add must have json tags for the fields to be serialized.
// Important: Run "operator-sdk generate k8s" to regenerate code after modifying this file
// Add custom validation using kubebuilder tags: https://book-v1.book.kubebuilder.io/beyond_basics/generating_crd.html

// ApacheBenchSpec defines the desired state of ApacheBench
type ApacheBenchSpec struct {
	// Concurrency is the number of multiple requests to perform at a time. Default is one request at a time.
	Concurrency uint32 `json:"concurrency,omitempty"`

	// Image is the container image (including tag) to use.
	Image string `json:"image,omitempty"`

	// Job is the JobSpec to override the default behavior of the benchmark Job.
	Job *batchv1.JobSpec `json:"job,omitempty"`

	// Requests is the number of requests to perform for the benchmarking session.
	// The default is to just perform a single request which usually leads to non-representative benchmarking results.
	Requests uint32 `json:"requests,omitempty"`

	// URL is the HTTP endpoint to benchmark.
	URL string `json:"url"`
}

// ApacheBenchStatus defines the observed state of ApacheBench
type ApacheBenchStatus struct {
	// Phase is a simple, high-level summary of where the ApacheBench is in its lifecycle.
	// There are five possible phase values:
	// Pending: The ApacheBench has been accepted by the Kubernetes system.
	// Running: At least one or more ApacheBench Jobs are currently running.
	// Complete: All of the ApacheBench Jobs have completed successfully.
	// Failed: At least one ApacheBench Job has experienced a failure.
	// Unknown: For some reason the state of the ApacheBench could not be obtained.
	Phase string `json:"phase"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApacheBench is the Schema for the apachebenches API
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=apachebenches,scope=Namespaced
type ApacheBench struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ApacheBenchSpec   `json:"spec,omitempty"`
	Status ApacheBenchStatus `json:"status,omitempty"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ApacheBenchList contains a list of ApacheBench
type ApacheBenchList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ApacheBench `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ApacheBench{}, &ApacheBenchList{})
}
