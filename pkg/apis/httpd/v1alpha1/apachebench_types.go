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

// ApacheBenchHTMLSpec defines the options for HTML output.
type ApacheBenchHTMLSpec struct {
	// Enabled toggles the printing of results in HTML tables.
	// Default table is two columns wide, with a white background.
	Enabled bool `json:"enabled"`

	// Table is the value to use as attributes for <table>. Attributes are inserted <table here >.
	Table string `json:"table,omitempty"`

	// TD is the value to use as attributes for <td>.
	TD string `json:"td,omitempty"`

	// TR is the value to use as attributes for <tr>.
	TR string `json:"tr,omitempty"`
}

// ApacheBenchSpec defines the desired state of ApacheBench
type ApacheBenchSpec struct {
	// Authenticate enables authentication for requests.
	// The "request.username" and "request.password" properties should be present in the Secret that is referenced by
	// the SecretName property.
	Authenticate bool `json:"authenticate,omitempty"`

	// AuthenticateProxy enables authentication for proxied requests.
	// The "proxy.username" and "proxy.password" properties should be present in the Secret that is referenced by
	// the SecretName property.
	AuthenticateProxy bool `json:"authenticateProxy,omitempty"`

	// Cookies is a map of key-value pairs to add as Cookie: lines to the request.
	Cookies map[string]string `json:"cookies,omitempty"`

	// Concurrency is the number of multiple requests to perform at a time. Default is one request at a time.
	Concurrency uint32 `json:"concurrency,omitempty"`

	// ConfigMapName is the name of a ConfigMap that contains POST or PUT data for requests.
	ConfigMapName string `json:"configMapName,omitempty"`

	// ContentType is the Content-type header to use for POST/PUT data, eg. application/x-www-form-urlencoded.
	// Default is text/plain.
	ContentType string `json:"contentType,omitempty"`

	// OmitLengthErrors disables errors if the length of the responses is not constant. This can be useful for dynamic pages.
	DisableLengthErrors bool `json:"disableLengthErrors,omitempty"`

	// DisableMedian disables the display of the median and standard deviation values.
	// Also disables the warning/error messages when the average and median are more than one or two times the standard
	// deviation apart. Default to the min/avg/max values. (legacy support).
	DisableMedian bool `json:"disableMedian,omitempty"`

	// OmitPercentageServed disables the "percentage served within XX [ms] table". (legacy support).
	DisablePercentageServed bool `json:"disablePercentageServed,omitempty"`

	// DisableProgress disables the progress count every 10% or 100 requests when processing more than 150 requests.
	DisableProgress bool `json:"disableProgress,omitempty"`

	// DisableSocketExit disables exit on socket receive errors.
	DisableSocketExit bool `json:"disableSocketExit,omitempty"`

	// EnableHEADRequests enables HEAD requests instead of GET.
	EnableHEADRequests bool `json:"enableHEADRequests,omitempty"`

	// Headers is a map of key-value pairs to add as headers to the request.
	Headers map[string]string `json:"headers,omitempty"`

	// HTML defines the HTML output options.
	HTML ApacheBenchHTMLSpec `json:"html,omitempty"`

	// HTTPMethod is a custom HTTP method for the requests.
	HTTPMethod string `json:"httpMethod,omitempty"`

	// Image is the container image (including tag) to use.
	Image string `json:"image,omitempty"`

	// Job is the JobSpec to override the default behavior of the benchmark Job.
	Job *batchv1.JobSpec `json:"job,omitempty"`

	// KeepAlive enables the HTTP KeepAlive feature, i.e., perform multiple requests within one HTTP session.
	KeepAlive bool `json:"keepAlive,omitempty"`

	// POSTDataKey is the name of the key in the ConfigMap specified in the ConfigMapName property that contains data
	// to POST with each request.
	POSTDataKey string `json:"postDataKey,omitempty"`

	// Proxy is the proxy server for the requests in the form proxy[:port].
	Proxy string `json:"proxy,omitempty"`

	// PUTDataKey is the name of the key in the ConfigMap specified in the ConfigMapName property that contains data
	// to PUT with each request.
	PUTDataKey string `json:"putDataKey,omitempty"`

	// Requests is the number of requests to perform for the benchmarking session.
	// The default is to just perform a single request which usually leads to non-representative benchmarking results.
	Requests uint32 `json:"requests,omitempty"`

	// SecretName is the name of the Secret containing authentication credentials and/or the client certificate.
	SecretName string `json:"secretName,omitempty"`

	// TimeLimit is the maximum number of seconds to spend for benchmarking.
	// This implies a 50000 value for Requests. Use this to benchmark the server within a fixed total amount of time.
	// Per default there is no timelimit.
	TimeLimit uint32 `json:"timeLimit,omitempty"`

	// Timeout is the maximum number of seconds to wait before the socket times out. Default is 30 seconds.
	Timeout uint32 `json:"timeout,omitempty"`

	// TLS defines the options for TLS connections.
	TLS ApacheBenchTLSSpec `json:"tls,omitempty"`

	// URL is the HTTP endpoint to benchmark.
	URL string `json:"url"`

	// Verbosity is the verbosity level.
	// 4 and above prints information on headers.
	// 3 and above prints response codes (404, 200, etc.).
	// 2 and above prints warnings and info.
	Verbosity uint32 `json:"verbosity,omitempty"`

	// WindowSize is the size of TCP send/receive buffer, in bytes.
	WindowSize uint32 `json:"windowSize,omitempty"`
}

// ApacheBenchStatus defines the observed state of ApacheBench
type ApacheBenchStatus struct {
	// Errors contains any errors that prevented the completion of the benchmark Job(s).
	Errors []string `json:"errors,omitempty"`

	// Phase is a simple, high-level summary of where the ApacheBench is in its lifecycle.
	// There are five possible phase values:
	// Pending: The ApacheBench has been accepted by the Kubernetes system.
	// Running: At least one or more ApacheBench Jobs are currently running.
	// Complete: All of the ApacheBench Jobs have completed successfully.
	// Failed: At least one ApacheBench Job has experienced a failure.
	// Unknown: For some reason the state of the ApacheBench could not be obtained.
	Phase string `json:"phase"`

	// Results contains the result output from each benchmark Job.
	Results []string `json:"results,omitempty"`
}

// ApacheBenchTLSSpec defines the options for TLS connections.
type ApacheBenchTLSSpec struct {
	// CipherSuite is the SSL/TLS cipher suite (See openssl ciphers).
	CipherSuite string `json:"cipherSuite,omitempty"`

	// Protocol is the SSL/TLS protocol.
	// (SSL2, SSL3, TLS1, TLS1.1, TLS1.2, or ALL). TLS1.1 and TLS1.2
	Protocol string `json:"protocol,omitempty"`
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
