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
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"strconv"

	v1a1 "github.com/jmckind/apache-bench-operator/pkg/apis/httpd/v1alpha1"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
)

const (
	// defaultContainerImage is the container image to use when one is not specified in the CR.
	defaultContainerImage = "httpd@sha256:223b88ef9a99261b07d2025d43799f45cace9b7b208195078b42cc2b922e453c" // 2.4.43-alpine
)

// addJobResultsToStatus will add the output from each Job pod to the ApacheBench status.
func (r *ReconcileApacheBench) addJobResultsToStatus(cr *v1a1.ApacheBench, job *batchv1.Job) error {
	clientset, err := kubernetes.NewForConfig(r.config)
	if err != nil {
		return err
	}

	opts := metav1.ListOptions{
		LabelSelector: fmt.Sprintf("job-name=%s", job.Name),
	}

	podList, err := clientset.CoreV1().Pods(job.Namespace).List(opts)
	if err != nil {
		return err
	}

	results := make([]string, 0)
	for _, pod := range podList.Items {
		logs, err := r.getPodLogs(clientset, pod)
		if err != nil {
			return err
		}
		results = append(results, string(logs))
	}
	cr.Status.Results = results

	return nil
}

// addStatusError will add the given error message to the Status.Errors property on the given ApacheBench.
// The value will not be added if it already exists.
func addStatusError(cr *v1a1.ApacheBench, msg string) {
	found := false

	for _, s := range cr.Status.Errors {
		if s == msg {
			found = true
			break
		}
	}

	if !found {
		cr.Status.Errors = append(cr.Status.Errors, msg)
	}
}

// fetchObject will retrieve the object with the given namespace and name using the Kubernetes API.
// The result will be stored in the given object.
func (r *ReconcileApacheBench) fetchObject(namespace string, name string, obj runtime.Object) error {
	return r.client.Get(context.TODO(), types.NamespacedName{Namespace: namespace, Name: name}, obj)
}

// getCommand will return the command to execute for the given ApacheBench.
func (r *ReconcileApacheBench) getCommand(cr *v1a1.ApacheBench) ([]string, error) {
	cmd := make([]string, 0)
	cmd = append(cmd, "ab")

	if cr.Spec.Authenticate {
		user, pass, err := r.getCredentialsFromSecret(cr, "request.username", "request.password")
		if err != nil {
			return nil, err
		}

		cmd = append(cmd, "-A")
		cmd = append(cmd, fmt.Sprintf("%s:%s", user, pass))
	}

	if cr.Spec.AuthenticateProxy {
		user, pass, err := r.getCredentialsFromSecret(cr, "proxy.username", "proxy.password")
		if err != nil {
			return nil, err
		}

		cmd = append(cmd, "-P")
		cmd = append(cmd, fmt.Sprintf("%s:%s", user, pass))
	}

	for key, val := range cr.Spec.Cookies {
		cmd = append(cmd, "-C")
		cmd = append(cmd, fmt.Sprintf("%s=%s", key, val))
	}

	if cr.Spec.Concurrency > 1 {
		cmd = append(cmd, "-c")
		cmd = append(cmd, strconv.FormatUint(uint64(cr.Spec.Concurrency), 10))
	}

	if len(cr.Spec.ContentType) > 0 {
		cmd = append(cmd, "-T")
		cmd = append(cmd, cr.Spec.ContentType)
	}

	if cr.Spec.DisableLengthErrors {
		cmd = append(cmd, "-l")
	}

	if cr.Spec.DisableMedian {
		cmd = append(cmd, "-S")
	}

	if cr.Spec.DisablePercentageServed {
		cmd = append(cmd, "-d")
	}

	if cr.Spec.DisableProgress {
		cmd = append(cmd, "-q")
	}

	if cr.Spec.DisableSocketExit {
		cmd = append(cmd, "-r")
	}

	if cr.Spec.EnableHEADRequests {
		cmd = append(cmd, "-i")
	}

	for key, val := range cr.Spec.Headers {
		cmd = append(cmd, "-H")
		cmd = append(cmd, fmt.Sprintf("%s=%s", key, val))
	}

	if cr.Spec.HTML.Enabled {
		cmd = append(cmd, "-w")

		if len(cr.Spec.HTML.Table) > 0 {
			cmd = append(cmd, "-x")
			cmd = append(cmd, cr.Spec.HTML.Table)
		}

		if len(cr.Spec.HTML.TD) > 0 {
			cmd = append(cmd, "-z")
			cmd = append(cmd, cr.Spec.HTML.TD)
		}

		if len(cr.Spec.HTML.TR) > 0 {
			cmd = append(cmd, "-y")
			cmd = append(cmd, cr.Spec.HTML.TR)
		}
	}

	if len(cr.Spec.HTTPMethod) > 0 {
		cmd = append(cmd, "-m")
		cmd = append(cmd, cr.Spec.HTTPMethod)
	}

	if cr.Spec.KeepAlive {
		cmd = append(cmd, "-k")
	}

	if len(cr.Spec.POSTDataKey) > 0 {
		cmd = append(cmd, "-p")
		cmd = append(cmd, fmt.Sprintf("/data/%s", cr.Spec.POSTDataKey))
	}

	if len(cr.Spec.Proxy) > 0 {
		cmd = append(cmd, "-X")
		cmd = append(cmd, cr.Spec.Proxy)
	}

	if len(cr.Spec.PUTDataKey) > 0 {
		cmd = append(cmd, "-u")
		cmd = append(cmd, fmt.Sprintf("/data/%s", cr.Spec.PUTDataKey))
	}

	if cr.Spec.Requests > 1 {
		cmd = append(cmd, "-n")
		cmd = append(cmd, strconv.FormatUint(uint64(cr.Spec.Requests), 10))
	}

	if cr.Spec.TimeLimit > 0 {
		cmd = append(cmd, "-t")
		cmd = append(cmd, strconv.FormatUint(uint64(cr.Spec.TimeLimit), 10))
	}

	if len(cr.Spec.TLS.CipherSuite) > 0 {
		cmd = append(cmd, "-Z")
		cmd = append(cmd, cr.Spec.TLS.CipherSuite)
	}

	if len(cr.Spec.TLS.Protocol) > 0 {
		cmd = append(cmd, "-f")
		cmd = append(cmd, cr.Spec.TLS.Protocol)
	}

	if cr.Spec.Verbosity > 0 {
		cmd = append(cmd, "-v")
		cmd = append(cmd, strconv.FormatUint(uint64(cr.Spec.Verbosity), 10))
	}

	if cr.Spec.WindowSize > 0 {
		cmd = append(cmd, "-b")
		cmd = append(cmd, strconv.FormatUint(uint64(cr.Spec.WindowSize), 10))
	}

	cmd = append(cmd, cr.Spec.URL)
	return cmd, nil
}

// getImage will return the container image to use for the given ApacheBench.
func getImage(cr *v1a1.ApacheBench) string {
	img := cr.Spec.Image
	if len(img) <= 0 {
		img = defaultContainerImage
	}
	return img
}

// getCredentialsFromSecret will return credential values using the given keys for the ApacheBench CR.
func (r *ReconcileApacheBench) getCredentialsFromSecret(cr *v1a1.ApacheBench, userKey string, passKey string) ([]byte, []byte, error) {
	secret := newSecret(cr)
	var failed = false

	if r.isObjectFound(cr.Namespace, secret.Name, secret) {
		_, ok := secret.Data[userKey]
		if !ok {
			failed = true
			addStatusError(cr, fmt.Sprintf("unable to locate username key '%s' in secret '%s'", userKey, secret.Name))
		}

		_, ok = secret.Data[passKey]
		if !ok {
			failed = true
			addStatusError(cr, fmt.Sprintf("unable to locate password key '%s' in secret '%s'", passKey, secret.Name))
		}
	} else {
		failed = true
		addStatusError(cr, fmt.Sprintf("unable to locate secret '%s'", secret.Name))
	}

	if failed {
		cr.Status.Phase = "Failed"
		if e := r.client.Status().Update(context.TODO(), cr); e != nil {
			return nil, nil, e
		}
		return nil, nil, errors.New("unable to locate credentials")
	}

	return secret.Data[userKey], secret.Data[passKey], nil
}

// getPodLogs will return the log output in bytes for the given Pod.
func (r *ReconcileApacheBench) getPodLogs(clientset *kubernetes.Clientset, pod corev1.Pod) ([]byte, error) {
	opts := corev1.PodLogOptions{} // Set size limit on result?

	req := clientset.CoreV1().Pods(pod.Namespace).GetLogs(pod.Name, &opts)
	logs, err := req.Stream()
	if err != nil {
		return nil, err
	}
	defer logs.Close()

	buf := new(bytes.Buffer)
	_, err = io.Copy(buf, logs)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

// getVolumeMounts will return the VolumeMounts for the given ApacheBench.
func getVolumeMounts(cr *v1a1.ApacheBench) []corev1.VolumeMount {
	vms := make([]corev1.VolumeMount, 0)

	if len(cr.Spec.POSTDataKey) < 0 || len(cr.Spec.PUTDataKey) < 0 {
		vms = append(vms, corev1.VolumeMount{
			Name:      "data",
			MountPath: "/data",
		})
	}

	return vms
}

// getVolumes will return the Volumes for the given ApacheBench.
func getVolumes(cr *v1a1.ApacheBench) []corev1.Volume {
	vs := make([]corev1.Volume, 0)

	if len(cr.Spec.POSTDataKey) < 0 || len(cr.Spec.PUTDataKey) < 0 {
		vs = append(vs, corev1.Volume{
			Name: "data",
			VolumeSource: corev1.VolumeSource{
				ConfigMap: &corev1.ConfigMapVolumeSource{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: cr.Spec.ConfigMapName,
					},
				},
			},
		})
	}

	return vs
}

// isObjectFound will perform a basic check that the given object exists via the Kubernetes API.
// If an error occurs as part of the check, the function will return false.
func (r *ReconcileApacheBench) isObjectFound(namespace string, name string, obj runtime.Object) bool {
	if err := r.fetchObject(namespace, name, obj); err != nil {
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
func (r *ReconcileApacheBench) newPodSpec(cr *v1a1.ApacheBench) (*corev1.PodSpec, error) {
	cmd, err := r.getCommand(cr)
	if err != nil {
		return nil, err
	}

	pod := corev1.PodSpec{
		Containers: []corev1.Container{{
			Command:         cmd,
			Image:           getImage(cr),
			ImagePullPolicy: corev1.PullIfNotPresent,
			Name:            "benchmark",
			VolumeMounts:    getVolumeMounts(cr),
		}},
		RestartPolicy: corev1.RestartPolicyOnFailure,
		Volumes:       getVolumes(cr),
	}

	return &pod, nil
}

// newPodTemplateSpec returns a new PodTemplateSpec for the given ApacheBench.
func (r *ReconcileApacheBench) newPodTemplateSpec(cr *v1a1.ApacheBench) (*corev1.PodTemplateSpec, error) {
	podSpec, err := r.newPodSpec(cr)
	if err != nil {
		return nil, err
	}

	ts := corev1.PodTemplateSpec{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    cr.Labels,
		},
		Spec: *podSpec,
	}

	return &ts, err
}

// newSecret returns a new Secret instance using the Name and Namespace from the given ApacheBench.
func newSecret(cr *v1a1.ApacheBench) *corev1.Secret {
	return &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Spec.SecretName,
			Namespace: cr.Namespace,
		},
	}
}

// reconcileJob will ensure that the Job for the given ApacheBench is present.
func (r *ReconcileApacheBench) reconcileJobs(cr *v1a1.ApacheBench) error {
	job := newJob(cr)
	if r.isObjectFound(cr.Namespace, job.Name, job) {
		if job.Status.Succeeded > 0 && job.Status.Succeeded == *job.Spec.Parallelism && cr.Status.Phase != string(batchv1.JobComplete) {
			// Mark status Phase as Complete
			cr.Status.Phase = string(batchv1.JobComplete)

			// Add results to the CR status
			if err := r.addJobResultsToStatus(cr, job); err != nil {
				return err
			}

			return r.client.Status().Update(context.TODO(), cr)
		}
		return nil // Job not complete, move along...
	}

	if cr.Spec.Job != nil {
		job.Spec = *cr.Spec.Job
	}

	template, err := r.newPodTemplateSpec(cr)
	if err != nil {
		return err
	}
	job.Spec.Template = *template

	if err := controllerutil.SetControllerReference(cr, job, r.scheme); err != nil {
		return err
	}
	return r.client.Create(context.TODO(), job)
}
