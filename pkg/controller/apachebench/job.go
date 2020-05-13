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

package apachebench

import (
	"context"
	"strconv"

	v1a1 "github.com/jmckind/apache-bench-operator/pkg/apis/httpd/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

const (
	// defaultContainerImage is the container image to use when one is not specified in the CR.
	defaultContainerImage = "httpd@sha256:223b88ef9a99261b07d2025d43799f45cace9b7b208195078b42cc2b922e453c" // 2.4.43-alpine
)

// fetchObject will retrieve the object with the given namespace and name using the Kubernetes API.
// The result will be stored in the given object.
func fetchObject(client client.Client, namespace string, name string, obj runtime.Object) error {
	return client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, obj)
}

// getCommand will return the command for the given ApacheBench.
func getCommand(cr *v1a1.ApacheBench) []string {
	cmd := make([]string, 0)
	cmd = append(cmd, "ab")

	if cr.Spec.Concurrency > 1 {
		cmd = append(cmd, "-c")
		cmd = append(cmd, strconv.FormatUint(uint64(cr.Spec.Concurrency), 10))
	}

	if cr.Spec.Requests > 1 {
		cmd = append(cmd, "-n")
		cmd = append(cmd, strconv.FormatUint(uint64(cr.Spec.Requests), 10))
	}

	cmd = append(cmd, cr.Spec.URL)
	return cmd
}

// getImage will return the container image to use for the given ApacheBench.
func getImage(cr *v1a1.ApacheBench) string {
	img := cr.Spec.Image
	if len(img) <= 0 {
		img = defaultContainerImage
	}
	return img
}

// isObjectFound will perform a basic check that the given object exists via the Kubernetes API.
// If an error occurs as part of the check, the function will return false.
func isObjectFound(client client.Client, namespace string, name string, obj runtime.Object) bool {
	if err := fetchObject(client, namespace, name, obj); err != nil {
		return false
	}
	return true
}

// newJob returns a new Job instance for the given ApacheBench.
func newJob(cr *v1a1.ApacheBench) *batchv1.Job {
	return &batchv1.Job{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
	}
}

// newPodSpec returns a new PodSpec for the given ApacheBench.
func newPodSpec(cr *v1a1.ApacheBench) corev1.PodSpec {
	pod := corev1.PodSpec{
		Containers: []corev1.Container{{
			Command:         getCommand(cr),
			Image:           getImage(cr),
			ImagePullPolicy: corev1.PullIfNotPresent,
			Name:            "benchmark",
		}},
		RestartPolicy: corev1.RestartPolicyOnFailure,
	}
	return pod
}

// newPodTemplateSpec returns a new PodTemplateSpec for the given ApacheBench.
func newPodTemplateSpec(cr *v1a1.ApacheBench) corev1.PodTemplateSpec {
	return corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Spec: newPodSpec(cr),
	}
}

// reconcileJob will ensure that the Job for the given ApacheBench is present.
func (r *ReconcileApacheBench) reconcileJobs(cr *v1a1.ApacheBench) error {
	job := newJob(cr)
	if isObjectFound(r.client, cr.Namespace, job.Name, job) {
		if job.Status.Succeeded > 0 && cr.Status.Phase != string(batchv1.JobComplete) {
			// Mark status Phase as Complete
			cr.Status.Phase = string(batchv1.JobComplete)
			return r.client.Status().Update(context.TODO(), cr)
		}
		return nil // Job not complete, move along...
	}

	if cr.Spec.Job != nil {
		job.Spec = *cr.Spec.Job
	}

	job.Spec.Template = newPodTemplateSpec(cr)

	if err := controllerutil.SetControllerReference(cr, job, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.TODO(), job)
}
