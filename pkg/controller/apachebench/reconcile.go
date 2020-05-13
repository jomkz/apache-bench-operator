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
	v1a1 "github.com/jmckind/apache-bench-operator/pkg/apis/httpd/v1alpha1"

	batchv1 "k8s.io/api/batch/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

// reconcileResources will reconcile all ApacheBench resources.
func (r *ReconcileApacheBench) reconcileResources(cr *v1a1.ApacheBench) error {
	log.Info("reconciling jobs")
	if err := r.reconcileJobs(cr); err != nil {
		return err
	}
	return nil
}

// watchResources will register Watches for each of the supported Resources.
func watchResources(c controller.Controller) error {
	// Watch for changes to primary resource ApacheBench
	if err := c.Watch(&source.Kind{Type: &v1a1.ApacheBench{}}, &handler.EnqueueRequestForObject{}); err != nil {
		return err
	}

	// Watch for changes to Secret sub-resources owned by ApacheBench instances.
	if err := watchOwnedResource(c, &corev1.Secret{}); err != nil {
		return err
	}

	// Watch for changes to Job sub-resources owned by ApacheBench instances.
	if err := watchOwnedResource(c, &batchv1.Job{}); err != nil {
		return err
	}

	return nil
}

// watchOwnedResource will register a Watch for the given resource owned by an ApacheBench instance.
func watchOwnedResource(c controller.Controller, obj runtime.Object) error {
	return c.Watch(&source.Kind{Type: obj}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &v1a1.ApacheBench{},
	})
}
