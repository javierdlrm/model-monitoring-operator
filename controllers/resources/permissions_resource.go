package resources

import (
	"fmt"

	"github.com/go-logr/logr"

	monitoringv1beta1 "github.com/javierdlrm/model-monitoring-operator/api/v1beta1"
	"github.com/javierdlrm/model-monitoring-operator/constants"

	corev1 "k8s.io/api/core/v1"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PermissionsBuilder defines the builder for managing permissions
type PermissionsBuilder struct {
	Assignee           string
	ModelMonitorConfig *monitoringv1beta1.ModelMonitorConfig
	Log                logr.Logger
}

// NewPermissionsBuilder creates a Permission builder
func NewPermissionsBuilder(assignee string, config *corev1.ConfigMap, log logr.Logger) *PermissionsBuilder {
	modelMonitorConfig, err := monitoringv1beta1.NewModelMonitorConfig(config)
	if err != nil {
		fmt.Printf("Failed to get model monitor config %s", err.Error())
		panic("Failed to get model monitor config")
	}
	return &PermissionsBuilder{
		Assignee:           assignee,
		ModelMonitorConfig: modelMonitorConfig,
		Log:                log,
	}
}

// CreateServiceAccountRoleAndBinding creates a Service account, Role and Role binding
func (b *PermissionsBuilder) CreateServiceAccountRoleAndBinding(modelMonitor *monitoringv1beta1.ModelMonitor) (*corev1.ServiceAccount, *rbacv1.Role, *rbacv1.RoleBinding, error) {
	metadata := modelMonitor.ObjectMeta
	serviceAccountName := constants.DefaultServiceAccountName(b.Assignee)
	roleName := constants.DefaultRoleName(b.Assignee)
	roleBindingName := constants.DefaultRoleBindingName(b.Assignee)

	// Service account
	serviceAccount := &corev1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceAccountName,
			Namespace: metadata.Namespace,
			Labels:    metadata.Labels,
		},
	}
	// Role
	role := &rbacv1.Role{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleName,
			Namespace: metadata.Namespace,
			Labels:    metadata.Labels,
		},
		Rules: []rbacv1.PolicyRule{
			{
				APIGroups: []string{""},
				Resources: []string{"pods"},
				Verbs:     []string{"*"},
			},
			{
				APIGroups: []string{""},
				Resources: []string{"services"},
				Verbs:     []string{"*"},
			},
		},
	}
	// Role Binding
	roleBinding := &rbacv1.RoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name:      roleBindingName,
			Namespace: metadata.Namespace,
			Labels:    metadata.Labels,
		},
		Subjects: []rbacv1.Subject{
			{
				Kind:      constants.ServiceAccount,
				Name:      serviceAccountName,
				Namespace: metadata.Namespace,
			},
		},
		RoleRef: rbacv1.RoleRef{
			Kind:     constants.Role,
			Name:     roleName,
			APIGroup: rbacv1.GroupName,
		},
	}
	return serviceAccount, role, roleBinding, nil
}
