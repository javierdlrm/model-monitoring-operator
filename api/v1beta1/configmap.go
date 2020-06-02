package v1beta1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/javierdlrm/model-monitoring-operator/constants"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ModelMonitorConfig defines the ModelMonitor configuration
// +k8s:openapi-gen=false
type ModelMonitorConfig struct {
	Monitoring      *MonitoringConfig      `json:"monitoring"`
	InferenceLogger *InferenceLoggerConfig `json:"inferencelogger"`
	Kafka           *KafkaConfig           `json:"kafka"`
}

// MonitoringConfig defines the Monitoring configuration
// +k8s:openapi-gen=false
type MonitoringConfig struct {
	Stats    []*StatConfig    `json:"stats"`
	Outliers []*OutlierConfig `json:"outliers"`
	Drift    []*DriftConfig   `json:"drift"`
}

// StatConfig defines a Statistic
// +k8s:openapi-gen=false
type StatConfig struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}

// OutlierConfig defines an Outlier detector
// +k8s:openapi-gen=false
type OutlierConfig struct {
	Name   string            `json:"name"`
	Params map[string]string `json:"params"`
}

// DriftConfig defines a Drift detector
// +k8s:openapi-gen=false
type DriftConfig struct {
	Name      string            `json:"name"`
	Threshold resource.Quantity `json:"threshold"`
	ShowAll   bool              `json:"showall"`
}

// InferenceLoggerConfig defines the configuration for the InferenceLogger service
// +k8s:openapi-gen=false
type InferenceLoggerConfig struct {
	ContainerImage string `json:"containerimage"`
}

// KafkaConfig defines the configuration for Kafka
// +k8s:openapi-gen=false
type KafkaConfig struct {
	Topic *KafkaTopicConfig `json:"topic"`
}

// KafkaTopicConfig defines the configuration for a Kafka topic
// +k8s:openapi-gen=false
type KafkaTopicConfig struct {
	Partitions        int32 `json:"partitions"`
	ReplicationFactor int16 `json:"replicationfactor"`
}

// GetModelMonitorConfig returns the ModelMonitor config
func GetModelMonitorConfig(client client.Client) (*ModelMonitorConfig, error) {
	configMap := &corev1.ConfigMap{}
	err := client.Get(context.TODO(), types.NamespacedName{Name: constants.ModelMonitorConfigMapName, Namespace: constants.ModelMonitoringNamespace}, configMap)
	if err != nil {
		return nil, err
	}

	modelMonitorConfigMap, err := NewModelMonitorConfig(configMap)
	if err != nil {
		return nil, err
	}

	return modelMonitorConfigMap, nil
}

// NewModelMonitorConfig creates a ModelMonitorConfig from a given configmap
func NewModelMonitorConfig(configMap *corev1.ConfigMap) (*ModelMonitorConfig, error) {

	monitoringConfig, err := getMonitoringConfig(configMap)
	if err != nil {
		return nil, err
	}

	inferenceLoggerConfig, err := getInferenceLoggerConfig(configMap)
	if err != nil {
		return nil, err
	}

	kafkaConfig, err := getKafkaConfig(configMap)
	if err != nil {
		return nil, err
	}

	return &ModelMonitorConfig{
		Monitoring:      monitoringConfig,
		InferenceLogger: inferenceLoggerConfig,
		Kafka:           kafkaConfig,
	}, nil
}

func getMonitoringConfig(configMap *corev1.ConfigMap) (*MonitoringConfig, error) {
	monitoringConfig := &MonitoringConfig{}
	key := constants.Monitoring.String()

	if data, ok := configMap.Data[key]; ok {
		err := json.Unmarshal([]byte(data), &monitoringConfig)
		if err != nil {
			return nil, fmt.Errorf("Unable to unmarshall %v json string due to %v ", key, err)
		}
	}
	return monitoringConfig, nil
}

func getInferenceLoggerConfig(configMap *corev1.ConfigMap) (*InferenceLoggerConfig, error) {
	inferenceLoggerConfig := &InferenceLoggerConfig{}
	key := constants.InferenceLogger.String()

	if data, ok := configMap.Data[key]; ok {
		err := json.Unmarshal([]byte(data), &inferenceLoggerConfig)
		if err != nil {
			return nil, fmt.Errorf("Unable to unmarshall %v json string due to %v ", key, err)
		}
	}

	return inferenceLoggerConfig, nil
}

func getKafkaConfig(configMap *corev1.ConfigMap) (*KafkaConfig, error) {
	kafkaConfig := &KafkaConfig{}
	key := constants.Kafka.String()

	if data, ok := configMap.Data[key]; ok {
		err := json.Unmarshal([]byte(data), &kafkaConfig)
		if err != nil {
			return nil, fmt.Errorf("Unable to unmarshall %v json string due to %v ", key, err)
		}
	}

	return kafkaConfig, nil
}
