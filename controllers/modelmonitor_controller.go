/*
Copyright 2020 Javier de la Rúa Martínez.

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

package controllers

import (
	"context"

	monitoringv1beta1 "github.com/javierdlrm/model-monitoring-operator/api/v1beta1"
	"github.com/javierdlrm/model-monitoring-operator/constants"
	"github.com/javierdlrm/model-monitoring-operator/controllers/reconcilers"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/tools/record"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"

	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// ModelMonitorReconciler reconciles a ModelMonitor object
type ModelMonitorReconciler struct {
	client.Client
	Log      logr.Logger
	Scheme   *runtime.Scheme
	Recorder record.EventRecorder
}

// +kubebuilder:rbac:groups=monitoring.hops.io,resources=modelmonitors,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=monitoring.hops.io,resources=modelmonitors/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=serving.knative.dev,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=serving.knative.dev,resources=services/status,verbs=get;update;patch
// +kubebuilder:rbac:groups="",resources=events,verbs=create;patch

// Reconcile reconciles ModelMonitor object request
func (r *ModelMonitorReconciler) Reconcile(req ctrl.Request) (ctrl.Result, error) {
	ctx := context.Background()
	log := r.Log.WithValues("modelmonitor", req.NamespacedName)

	var err error

	// Fetch the ModelMonitor config
	modelMonitor := &monitoringv1beta1.ModelMonitor{}
	if err = r.Get(ctx, req.NamespacedName, modelMonitor); err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return ctrl.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Get configmap
	configMap := &corev1.ConfigMap{}
	if err = r.Get(ctx, types.NamespacedName{Name: constants.ModelMonitorConfigMapName, Namespace: constants.ModelMonitoringNamespace}, configMap); err != nil {
		log.Error(err, "Failed to find ConfigMap", "name", constants.ModelMonitorConfigMapName, "namespace", constants.ModelMonitoringNamespace)
		// Error reading the object - requeue the request.
		return ctrl.Result{}, err
	}

	// Build reconcilers
	inferenceAdapterReconciler := reconcilers.NewInferenceAdapterReconciler(r.Client, r.Scheme, r.Log, r.Recorder, configMap)

	// Reconcile InferenceAdapter
	if err = inferenceAdapterReconciler.Reconcile(modelMonitor); err != nil {
		log.Error(err, "Failed to reconcile")
		r.Recorder.Eventf(modelMonitor, corev1.EventTypeWarning, "InternalError", err.Error())
		return ctrl.Result{}, err
	}

	// Update status
	if err = r.Status().Update(ctx, modelMonitor); err != nil {
		r.Recorder.Eventf(modelMonitor, corev1.EventTypeWarning, "InternalError", err.Error())
		return ctrl.Result{}, err
	}

	return ctrl.Result{}, nil
}

// SetupWithManager creates new managed controller
func (r *ModelMonitorReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&monitoringv1beta1.ModelMonitor{}).
		Owns(&knservingv1.Service{}).
		Complete(r)
}
