package constants

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
	ModelMonitorControllerName           = ModelMonitorName + "-controller"
	ControllerLabelName                  = ModelMonitorControllerName + "-manager"
	DefaultInferenceAdapterTimeout int64 = 300
	DefaultScalingTarget                 = "1"
	DefaultMinReplicas                   = 1
)

// ModelMonitorComponent enum
type ModelMonitorComponent string

// ModelMonitor fields
const (
	Model            ModelMonitorComponent = "model"
	Monitoring       ModelMonitorComponent = "monitoring"
	InferenceAdapter ModelMonitorComponent = "inferenceadapter"
	KafkaTopic       ModelMonitorComponent = "kafkatopic"
)

// InferenceAdapter constants
const (
	InferenceAdapterModelLabel           = "model"
	InferenceAdapterEnvKafkaTopicLabel   = "KAFKA_TOPIC"
	InferenceAdapterEnvKafkaBrokersLabel = "KAFKA_BROKERS"
)

// KafkaTopic constants
const (
	KafkaTopicLabel   = "topic"
	KafkaBrokersLabel = "brokers"
)

// DefaultInferenceAdapterName builds a default name
func DefaultInferenceAdapterName(prefix string) string {
	return prefix + "-" + InferenceAdapter.String()
}

// DefaultKafkaTopicName build a default Kafka Topic name
func DefaultKafkaTopicName(modelName string) string {
	return modelName + "-topic"
}

// String return ModelMonitorComponent as string
func (c ModelMonitorComponent) String() string {
	return string(c)
}
