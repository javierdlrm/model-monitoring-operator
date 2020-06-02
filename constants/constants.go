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
	Kafka           ModelMonitorComponent = "kafka"
)

// InferenceLogger constants
const (
	InferenceLoggerModelLabel                          = "model"
	InferenceLoggerEnvKafkaBrokersLabel                = "KAFKA_BROKERS"
	InferenceLoggerEnvKafkaTopicLabel                  = "KAFKA_TOPIC"
	InferenceLoggerEnvKafkaTopicPartitionsLabel        = "KAFKA_TOPIC_PARTITIONS"
	InferenceLoggerEnvKafkaTopicReplicationFactorLabel = "KAFKA_TOPIC_REPLICATION_FACTOR"
)

// KafkaTopic constants
const (
	KafkaTopicLabel                  = "topic"
	KafkaBrokersLabel                = "brokers"
	KafkaTopicPartitionsLabel        = "partitions"
	KafkaTopicReplicationFactorLabel = "replicationfactor"
)

// DefaultInferenceLoggerName builds a default name
func DefaultInferenceLoggerName(prefix string) string {
	return prefix + "-" + InferenceLogger.String()
}

// DefaultKafkaTopicName build a default Kafka Topic name
func DefaultKafkaTopicName(modelName string) string {
	// Don't change. Defaults topic names must match with InferenceLogger
	return modelName + "-inference-topic"
}

// String return ModelMonitorComponent as string
func (c ModelMonitorComponent) String() string {
	return string(c)
}
