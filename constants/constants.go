package constants

import (
	corev1 "k8s.io/api/core/v1"
	"knative.dev/serving/pkg/apis/autoscaling"
)

// ModelMonitoring Operator constants
var (
	ModelMonitoringName         = "model-monitoring"
	ModelMonitoringAPIGroupName = "monitoring.model.dev"
	ModelMonitoringNamespace    = ModelMonitoringName + "-system"
)

// ModelMonitor constants
var (
	ModelMonitorName          = "modelmonitor"
	ModelMonitorAPIName       = "modelmonitors"
	ModelMonitorPodLabelKey   = ModelMonitoringAPIGroupName + "/" + ModelMonitorName
	ModelMonitorConfigMapName = ModelMonitoringName + "-" + ModelMonitorName + "-config"
	ModelMonitorContainerName = ModelMonitorName + "-container"
)

// ModelMonitor Controller Constants
var (
	ModelMonitorControllerName = ModelMonitorName + "-controller"
	ControllerLabelName        = ModelMonitorControllerName + "-manager"
)

// ModelMonitorComponent enum
type ModelMonitorComponent string

// ModelMonitor fields
const (
	Job             ModelMonitorComponent = "job"
	InferenceLogger ModelMonitorComponent = "inferenceLogger"
)

// InferenceLogger constants
const (
	InferenceLoggerNameSuffix = "inferencelogger"
	// Labels
	InferenceLoggerModelLabel                          = "model"
	InferenceLoggerEnvKafkaBrokersLabel                = "KAFKA_BROKERS"
	InferenceLoggerEnvKafkaTopicLabel                  = "KAFKA_TOPIC"
	InferenceLoggerEnvKafkaTopicPartitionsLabel        = "KAFKA_TOPIC_PARTITIONS"
	InferenceLoggerEnvKafkaTopicReplicationFactorLabel = "KAFKA_TOPIC_REPLICATION_FACTOR"
)

// InferenceLogger defaults
var (
	InferenceLoggerDefaultCPU                               = "0.1"
	InferenceLoggerDefaultMemory                            = "128Mi"
	InferenceLoggerDefaultTimeout                     int64 = 300
	InferenceLoggerDefaultScalingClass                      = autoscaling.KPA // kpa or hpa
	InferenceLoggerDefaultScalingMetric                     = "concurrency"   // concurrency, rps or cpu (hpa required)
	InferenceLoggerDefaultScalingTarget                     = 100
	InferenceLoggerDefaultTargetUtilizationPercentage       = "70"
	InferenceLoggerDefaultMinScale                          = 1 // 0 if scale-to-zero is desired
	InferenceLoggerDefaultMaxScale                          = 0 // 0 means limitless
	InferenceLoggerDefaultWindow                            = "60s"
	InferenceLoggerDefaultPanicWindow                       = "10" // percentage of StableWindow
	InferenceLoggerDefaultPanicThreshold                    = "200"
)

// Job constants
const (
	MonitoringJobNameSuffix = "monitoring-job"
	// Labels
	MonitoringJobSparkVersionLabel           = "version"
	MonitoringJobEnvVarModelInfoLabel        = "MODEL_INFO"
	MonitoringJobEnvVarMonitoringConfigLabel = "MONITORING_CONFIG"
	MonitoringJobEnvVarStorageConfigLabel    = "STORAGE_CONFIG"
	MonitoringJobEnvVarJobConfigLabel        = "JOB_CONFIG"
)

// TODO: Add Driver and Executor variables to api
// Job template & defaults
var (
	// Permissions
	MonitoringJobAssignee = "spark"
	// Spark
	MonitoringJobSparkVersion    = "2.4.5"
	MonitoringJobImagePullPolicy = "Always"
	// Driver
	MonitoringJobDriverCores     int32 = 1
	MonitoringJobDriverCoreLimit       = "1000m"
	MonitoringJobDriverMemory          = "512m"
	// Executor
	MonitoringJobExecutorCores     int32 = 1
	MonitoringJobExecutorCoreLimit       = "1000m"
	MonitoringJobExecutorMemory          = "512m"
	MonitoringJobExecutorInstances int32 = 1
	// Volume
	MonitoringJobVolumeName         = "test-volume"
	MonitoringJobVolumeHostPath     = "/tmp"
	MonitoringJobVolumeHostPathType = corev1.HostPathDirectory
	// VolumeMount
	MonitoringJobVolumeMountName = "test-volume"
	MonitoringJobVolumeMountPath = "/tmp"
	// Monitoring
	MonitoringJobPrometheusExportDriverMetrics         = true
	MonitoringJobPrometheusExportExecutorMetrics       = true
	MonitoringJobPrometheusJmxExporterJar              = "/prometheus/jmx_prometheus_javaagent-0.11.0.jar"
	MonitoringJobPrometheusPort                  int32 = 8090
)

// KafkaTopic constants
const (
	KafkaTopicNameSuffix = "inference-topic"
	// Labels
	KafkaTopicLabel                  = "topic"
	KafkaBrokersLabel                = "brokers"
	KafkaTopicPartitionsLabel        = "partitions"
	KafkaTopicReplicationFactorLabel = "replicationFactor"
)

// Permissions
const (
	ServiceAccount = "ServiceAccount"
	Role           = "Role"
	// Suffix
	ServiceAccountNameSuffix = "sa"
	RoleNameSuffix           = "r"
	RoleBindingNameSuffix    = "rb"
)

// DefaultInferenceLoggerName builds a default name
func DefaultInferenceLoggerName(prefix string) string {
	return prefix + "-" + InferenceLoggerNameSuffix
}

// DefaultKafkaTopicName build a default Kafka Topic name
func DefaultKafkaTopicName(modelName string) string {
	// Don't change. Defaults topic names must match with InferenceLogger
	return modelName + "-" + KafkaTopicNameSuffix
}

// DefaultMonitoringJobName build a default Kafka Topic name
func DefaultMonitoringJobName(modelName string) string {
	// Don't change. Defaults topic names must match with InferenceLogger
	return modelName + "-" + MonitoringJobNameSuffix
}

// DefaultServiceAccountName build a default Service Account name
func DefaultServiceAccountName(assignee string) string {
	return assignee + "-" + ServiceAccountNameSuffix
}

// DefaultRoleName build a default Role name
func DefaultRoleName(assignee string) string {
	return assignee + "-" + RoleNameSuffix
}

// DefaultRoleBindingName build a default Role Binding name
func DefaultRoleBindingName(assignee string) string {
	return assignee + "-" + RoleBindingNameSuffix
}

// String return ModelMonitorComponent as string
func (c ModelMonitorComponent) String() string {
	return string(c)
}
