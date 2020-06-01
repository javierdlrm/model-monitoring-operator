package reconcilers

import (
	"context"
	"fmt"

	monitoringv1beta1 "github.com/javierdlrm/model-monitoring-operator/api/v1beta1"
	"github.com/javierdlrm/model-monitoring-operator/constants"
	"github.com/javierdlrm/model-monitoring-operator/controllers/resources"

	"github.com/go-logr/logr"

	"k8s.io/client-go/tools/record"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"knative.dev/pkg/kmp"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

// InferenceAdapterReconciler defines a reconciler for InferenceAdapter
type InferenceAdapterReconciler struct {
	Client   client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
	Builder  *resources.InferenceAdapterBuilder
}

// NewInferenceAdapterReconciler creates a new reconciler for InferenceAdapter
func NewInferenceAdapterReconciler(client client.Client, scheme *runtime.Scheme, log logr.Logger, recorder record.EventRecorder,
	config *corev1.ConfigMap) *InferenceAdapterReconciler {

	return &InferenceAdapterReconciler{
		Client:   client,
		Scheme:   scheme,
		Log:      log,
		Recorder: recorder,
		Builder:  resources.NewInferenceAdapterBuilder(client, config),
	}
}

// Reconcile a given ModelMonitor declarative config
func (r *InferenceAdapterReconciler) Reconcile(modelMonitor *monitoringv1beta1.ModelMonitor) error {
	log := r.Log.WithValues("inferenceadapter", modelMonitor.Namespace+"/"+modelMonitor.Name)

	serviceName := constants.DefaultInferenceAdapterName(modelMonitor.Name)

	var service *knservingv1.Service
	var err error
	service, err = r.Builder.CreateInferenceAdapterService(serviceName, modelMonitor)
	if err != nil {
		return err
	}

	if service == nil {
		if err = r.finalizeService(serviceName, modelMonitor.Namespace, log); err != nil {
			return err
		}

		// TODO: Modify status
		return nil
	}

	// _, err => status, err
	if _, err := r.reconcileService(modelMonitor, service, log); err != nil {
		return err
	}

	// TODO: Modify status
	return nil
}

func (r *InferenceAdapterReconciler) finalizeService(serviceName string, namespace string, log logr.Logger) error {
	existing := &knservingv1.Service{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Name: serviceName, Namespace: namespace}, existing); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		log.Info("Deleting Knative Service", "namespace", namespace, "name", serviceName)
		if err := r.Client.Delete(context.TODO(), existing, client.PropagationPolicy(metav1.DeletePropagationBackground)); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}
	return nil
}

func (r *InferenceAdapterReconciler) reconcileService(modelMonitor *monitoringv1beta1.ModelMonitor, desired *knservingv1.Service,
	log logr.Logger) (*knservingv1.ServiceStatus, error) {

	// Set ModelMonitor as owner of desired service
	if err := controllerutil.SetControllerReference(modelMonitor, desired, r.Scheme); err != nil {
		return nil, err
	}

	// Create service if does not exist
	existing := &knservingv1.Service{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			log.Info("Creating Knative Service", "namespace", desired.Namespace, "name", desired.Name)
			return &desired.Status, r.Client.Create(context.TODO(), desired)
		}
		return nil, err
	}

	// Return if no differences to reconcile.
	if semanticEquals(desired, existing) {
		log.Info("No differences found")
		return &existing.Status, nil
	}

	// Reconcile differences and update
	diff, err := kmp.SafeDiff(desired.Spec.ConfigurationSpec, existing.Spec.ConfigurationSpec)
	if err != nil {
		return &existing.Status, fmt.Errorf("Failed to diff Knative Service: %v", err)
	}

	log.Info("Reconciling Knative Service diff (-desired, +observed):", "diff", diff)
	log.Info("Updating Knative Service", "namespace", desired.Namespace, "name", desired.Name)

	existing.Spec.ConfigurationSpec = desired.Spec.ConfigurationSpec
	existing.ObjectMeta.Labels = desired.ObjectMeta.Labels
	if err := r.Client.Update(context.TODO(), existing); err != nil {
		return &existing.Status, err
	}

	return &existing.Status, nil
}

func semanticEquals(desired, service *knservingv1.Service) bool {
	return equality.Semantic.DeepEqual(desired.Spec.ConfigurationSpec, service.Spec.ConfigurationSpec) &&
		equality.Semantic.DeepEqual(desired.ObjectMeta.Labels, service.ObjectMeta.Labels)
}
