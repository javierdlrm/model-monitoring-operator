package resources

import (
	"encoding/json"
	"fmt"

	"github.com/go-logr/logr"

	monitoringv1beta1 "github.com/javierdlrm/model-monitoring-operator/api/v1beta1"
	"github.com/javierdlrm/model-monitoring-operator/constants"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	sparkv1beta2 "github.com/GoogleCloudPlatform/spark-on-k8s-operator/pkg/apis/sparkoperator.k8s.io/v1beta2"
)

// MonitoringJobBuilder defines the builder for Monitoring job
type MonitoringJobBuilder struct {
	ModelMonitorConfig *monitoringv1beta1.ModelMonitorConfig
	Permissions        *PermissionsBuilder
	Log                logr.Logger
}

// NewMonitoringJobBuilder creates a Monitoring job builder
func NewMonitoringJobBuilder(config *corev1.ConfigMap, log logr.Logger) *MonitoringJobBuilder {
	modelMonitorConfig, err := monitoringv1beta1.NewModelMonitorConfig(config)
	if err != nil {
		fmt.Printf("Failed to get model monitor config %s", err.Error())
		panic("Failed to get model monitor config")
	}
	return &MonitoringJobBuilder{
		ModelMonitorConfig: modelMonitorConfig,
		Permissions:        NewPermissionsBuilder(constants.MonitoringJobAssignee, config, log),
		Log:                log,
	}
}

// CreateMonitoringJobSparkApp creates the Spark Application for Monitoring job
func (b *MonitoringJobBuilder) CreateMonitoringJobSparkApp(monitoringJobName string, modelMonitor *monitoringv1beta1.ModelMonitor) (*sparkv1beta2.SparkApplication, error) {

	// Specs
	metadata := modelMonitor.ObjectMeta
	modelSpec := modelMonitor.Spec.Model
	monitoringSpec := modelMonitor.Spec.Monitoring
	storageSpec := modelMonitor.Spec.Storage
	jobSpec := modelMonitor.Spec.Job

	// Service account
	serviceAccount := constants.DefaultServiceAccountName(b.Permissions.Assignee)

	// Env vars (json format)
	modelSpecBytes, err := json.Marshal(modelSpec)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal %v object to %v ", modelSpec, err)
	}
	monitoringSpecBytes, err := json.Marshal(monitoringSpec)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal %v object to %v ", monitoringSpec, err)
	}
	storageSpecBytes, err := json.Marshal(storageSpec)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal %v object to %v ", storageSpec, err)
	}
	jobSpecBytes, err := json.Marshal(jobSpec)
	if err != nil {
		return nil, fmt.Errorf("Unable to marshal %v object to %v ", jobSpec, err)
	}

	// Spark application
	sparkApp := &sparkv1beta2.SparkApplication{
		ObjectMeta: metav1.ObjectMeta{
			Name:      monitoringJobName,
			Namespace: metadata.Namespace,
			Labels:    metadata.Labels,
		},
		Spec: sparkv1beta2.SparkApplicationSpec{
			Type:                sparkv1beta2.ScalaApplicationType,
			Mode:                sparkv1beta2.ClusterMode,
			Image:               &b.ModelMonitorConfig.Job.ContainerImage,
			ImagePullPolicy:     &constants.MonitoringJobImagePullPolicy,
			MainClass:           &b.ModelMonitorConfig.Job.MainClass,
			MainApplicationFile: &b.ModelMonitorConfig.Job.MainApplicationFile,
			SparkVersion:        constants.MonitoringJobSparkVersion,
			RestartPolicy: sparkv1beta2.RestartPolicy{
				Type: sparkv1beta2.Never,
			},
			Volumes: []corev1.Volume{
				{
					Name: constants.MonitoringJobVolumeName,
					VolumeSource: corev1.VolumeSource{
						HostPath: &corev1.HostPathVolumeSource{
							Path: constants.MonitoringJobVolumeHostPath,
							Type: &constants.MonitoringJobVolumeHostPathType,
						},
					},
				},
			},
			Driver: sparkv1beta2.DriverSpec{
				SparkPodSpec: sparkv1beta2.SparkPodSpec{
					Cores:          &constants.MonitoringJobDriverCores,
					CoreLimit:      &constants.MonitoringJobDriverCoreLimit,
					Memory:         &constants.MonitoringJobDriverMemory,
					Labels:         map[string]string{constants.MonitoringJobSparkVersionLabel: constants.MonitoringJobSparkVersion},
					ServiceAccount: &serviceAccount,
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      constants.MonitoringJobVolumeMountName,
							MountPath: constants.MonitoringJobVolumeMountPath,
						},
					},
					EnvVars: map[string]string{
						constants.MonitoringJobEnvVarModelInfoLabel:        string(modelSpecBytes),
						constants.MonitoringJobEnvVarMonitoringConfigLabel: string(monitoringSpecBytes),
						constants.MonitoringJobEnvVarStorageConfigLabel:    string(storageSpecBytes),
						constants.MonitoringJobEnvVarJobConfigLabel:        string(jobSpecBytes),
					},
				},
			},
			Executor: sparkv1beta2.ExecutorSpec{
				SparkPodSpec: sparkv1beta2.SparkPodSpec{
					Cores:     &constants.MonitoringJobExecutorCores,
					CoreLimit: &constants.MonitoringJobExecutorCoreLimit,
					Memory:    &constants.MonitoringJobExecutorMemory,
					Labels:    map[string]string{constants.MonitoringJobSparkVersionLabel: constants.MonitoringJobSparkVersion},
					VolumeMounts: []corev1.VolumeMount{
						{
							Name:      constants.MonitoringJobVolumeMountName,
							MountPath: constants.MonitoringJobVolumeMountPath,
						},
					},
				},
				Instances: &constants.MonitoringJobExecutorInstances,
			},
		},
	}

	return sparkApp, nil
}
