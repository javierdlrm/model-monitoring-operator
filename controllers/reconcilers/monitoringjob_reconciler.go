package reconcilers

import (
	"context"

	monitoringv1beta1 "github.com/javierdlrm/model-monitoring-operator/api/v1beta1"
	"github.com/javierdlrm/model-monitoring-operator/constants"
	"github.com/javierdlrm/model-monitoring-operator/controllers/resources"

	"github.com/go-logr/logr"

	"k8s.io/client-go/tools/record"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"

	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	sparkv1beta2 "github.com/GoogleCloudPlatform/spark-on-k8s-operator/pkg/apis/sparkoperator.k8s.io/v1beta2"
)

// MonitoringJobReconciler defines a reconciler for Monitoring job
type MonitoringJobReconciler struct {
	Client   client.Client
	Scheme   *runtime.Scheme
	Log      logr.Logger
	Recorder record.EventRecorder
	Builder  *resources.MonitoringJobBuilder
}

// NewMonitoringJobReconciler creates a new reconciler for Monitoring job
func NewMonitoringJobReconciler(client client.Client, scheme *runtime.Scheme, log logr.Logger, recorder record.EventRecorder,
	config *corev1.ConfigMap) *MonitoringJobReconciler {

	return &MonitoringJobReconciler{
		Client:   client,
		Scheme:   scheme,
		Log:      log,
		Recorder: recorder,
		Builder:  resources.NewMonitoringJobBuilder(config, log),
	}
}

// Reconcile a given ModelMonitor declarative config
func (r *MonitoringJobReconciler) Reconcile(modelMonitor *monitoringv1beta1.ModelMonitor) error {
	monitoringJobName := constants.DefaultMonitoringJobName(modelMonitor.Name)
	r.Log = r.Log.WithValues("monitoringJob", modelMonitor.Namespace+"/"+modelMonitor.Name, "monitoringJobName", monitoringJobName)

	var sparkApp *sparkv1beta2.SparkApplication
	var err error
	sparkApp, err = r.Builder.CreateMonitoringJobSparkApp(monitoringJobName, modelMonitor)
	if err != nil {
		return err
	}

	if sparkApp == nil {
		if err = r.finalizeSparkApp(monitoringJobName, modelMonitor.Namespace); err != nil {
			return err
		}
		// TODO: Modify status
		return nil
	}

	// _, err => status, err
	if _, err := r.reconcileSparkApp(modelMonitor, sparkApp); err != nil {
		return err
	}
	// TODO: Modify status
	return nil
}

func (r *MonitoringJobReconciler) finalizeSparkApp(monitoringJobName string, namespace string) error {
	existing := &sparkv1beta2.SparkApplication{}
	if err := r.Client.Get(context.TODO(), types.NamespacedName{Name: monitoringJobName, Namespace: namespace}, existing); err != nil {
		if !errors.IsNotFound(err) {
			return err
		}
	} else {
		r.Log.Info("Deleting Spark Application", "namespace", namespace, "name", monitoringJobName)
		if err := r.Client.Delete(context.TODO(), existing, client.PropagationPolicy(metav1.DeletePropagationBackground)); err != nil {
			if !errors.IsNotFound(err) {
				return err
			}
		}
	}
	return nil
}

func (r *MonitoringJobReconciler) reconcileSparkApp(modelMonitor *monitoringv1beta1.ModelMonitor, desired *sparkv1beta2.SparkApplication) (*sparkv1beta2.SparkApplicationStatus, error) {
	// Set ModelMonitor as owner of desired spark app
	if err := controllerutil.SetControllerReference(modelMonitor, desired, r.Scheme); err != nil {
		return nil, err
	}

	// Create spark app if does not exist
	existing := &sparkv1beta2.SparkApplication{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: desired.Name, Namespace: desired.Namespace}, existing)
	if err != nil {
		if errors.IsNotFound(err) {
			// Create service account, role and role binding if do not exist
			sa, role, roleBinding, err := r.Builder.Permissions.CreateServiceAccountRoleAndBinding(modelMonitor)
			if err = r.reconcileSparkAppPermissions(sa, role, roleBinding); err != nil {
				return &desired.Status, err
			}

			r.Log.Info("Creating Spark Application", "namespace", desired.Namespace, "name", desired.Name)
			return &desired.Status, r.Client.Create(context.TODO(), desired)
		}
		return nil, err
	}

	// Return if no differences to reconcile.
	if sparkAppsSemanticEquals(desired, existing) {
		r.Log.Info("No differences found")
		return &existing.Status, nil
	}

	r.Log.Info("Updating Spark Application", "namespace", desired.Namespace, "name", desired.Name)
	existing.Spec = desired.Spec
	existing.ObjectMeta.Labels = desired.ObjectMeta.Labels
	if err := r.Client.Update(context.TODO(), existing); err != nil {
		return &existing.Status, err
	}

	return &existing.Status, nil
}

func (r *MonitoringJobReconciler) reconcileSparkAppPermissions(desiredServiceAccount *corev1.ServiceAccount, desiredRole *rbacv1.Role, desiredRoleBinding *rbacv1.RoleBinding) error {
	namespace := desiredServiceAccount.Namespace

	// Names
	serviceAccountName := constants.DefaultServiceAccountName(constants.MonitoringJobAssignee)
	roleName := constants.DefaultRoleName(constants.MonitoringJobAssignee)
	roleBindingName := constants.DefaultRoleBindingName(constants.MonitoringJobAssignee)

	// Create service account if does not exist
	existingSA := &corev1.ServiceAccount{}
	err := r.Client.Get(context.TODO(), types.NamespacedName{Name: serviceAccountName, Namespace: namespace}, existingSA)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Creating Service Account", "namespace", namespace, "name", serviceAccountName)
			if err = r.Client.Create(context.TODO(), desiredServiceAccount); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	// Create role if does not exist
	existingR := &rbacv1.Role{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: roleName, Namespace: namespace}, existingR)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Creating Role", "namespace", namespace, "name", roleName)
			if err = r.Client.Create(context.TODO(), desiredRole); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	// Create role binding if does not exist
	existingRB := &rbacv1.RoleBinding{}
	err = r.Client.Get(context.TODO(), types.NamespacedName{Name: roleBindingName, Namespace: namespace}, existingRB)
	if err != nil {
		if errors.IsNotFound(err) {
			r.Log.Info("Creating Role Binding", "namespace", namespace, "name", roleBindingName)
			if err = r.Client.Create(context.TODO(), desiredRoleBinding); err != nil {
				return err
			}
		} else {
			return err
		}
	}
	return nil
}

func sparkAppsSemanticEquals(desired *sparkv1beta2.SparkApplication, sparkApp *sparkv1beta2.SparkApplication) bool {
	return equality.Semantic.DeepEqual(desired.Spec, sparkApp.Spec) &&
		equality.Semantic.DeepEqual(desired.ObjectMeta.Labels, sparkApp.ObjectMeta.Labels)
}
