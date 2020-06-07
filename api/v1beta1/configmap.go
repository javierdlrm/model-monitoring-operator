package v1beta1

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/javierdlrm/model-monitoring-operator/constants"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

// ModelMonitorConfig defines the ModelMonitor configuration
// +k8s:openapi-gen=false
type ModelMonitorConfig struct {
	Job             *JobConfig             `json:"job"`
	InferenceLogger *InferenceLoggerConfig `json:"inferenceLogger"`
}

// InferenceLoggerConfig defines the configuration for the InferenceLogger service
// +k8s:openapi-gen=false
type InferenceLoggerConfig struct {
	ContainerImage string `json:"containerImage"`
}

// JobConfig defines the configuration for the Monitoring job
// +k8s:openapi-gen=false
type JobConfig struct {
	ContainerImage      string `json:"containerImage"`
	MainClass           string `json:"mainClass"`
	MainApplicationFile string `json:"mainApplicationFile"`
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

	jobConfig, err := getJobConfig(configMap)
	if err != nil {
		return nil, err
	}

	inferenceLoggerConfig, err := getInferenceLoggerConfig(configMap)
	if err != nil {
		return nil, err
	}

	return &ModelMonitorConfig{
		Job:             jobConfig,
		InferenceLogger: inferenceLoggerConfig,
	}, nil
}

func getJobConfig(configMap *corev1.ConfigMap) (*JobConfig, error) {
	jobConfig := &JobConfig{}
	key := constants.Job.String()

	if data, ok := configMap.Data[key]; ok {
		err := json.Unmarshal([]byte(data), &jobConfig)
		if err != nil {
			return nil, fmt.Errorf("Unable to unmarshall %v json string due to %v ", key, err)
		}
	}
	return jobConfig, nil
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
