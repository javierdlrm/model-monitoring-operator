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
	ModelMonitorControllerName          = ModelMonitorName + "-controller"
	ControllerLabelName                 = ModelMonitorControllerName + "-manager"
	DefaultInferenceLoggerTimeout int64 = 300
	DefaultScalingTarget                = "1"
	DefaultMinReplicas                  = 1
)

// ModelMonitorComponent enum
type ModelMonitorComponent string

// ModelMonitor fields
const (
	Model           ModelMonitorComponent = "model"
	Monitoring      ModelMonitorComponent = "monitoring"
	InferenceLogger ModelMonitorComponent = "inferencelogger"
	KafkaTopic      ModelMonitorComponent = "kafkatopic"
)

// InferenceLogger constants
const (
	InferenceLoggerModelLabel           = "model"
	InferenceLoggerEnvKafkaTopicLabel   = "KAFKA_TOPIC"
	InferenceLoggerEnvKafkaBrokersLabel = "KAFKA_BROKERS"
)

// KafkaTopic constants
const (
	KafkaTopicLabel   = "topic"
	KafkaBrokersLabel = "brokers"
)

// DefaultInferenceLoggerName builds a default name
func DefaultInferenceLoggerName(prefix string) string {
	return prefix + "-" + InferenceLogger.String()
}

// DefaultKafkaTopicName build a default Kafka Topic name
func DefaultKafkaTopicName(modelName string) string {
	return modelName + "-topic"
}

// String return ModelMonitorComponent as string
func (c ModelMonitorComponent) String() string {
	return string(c)
}
