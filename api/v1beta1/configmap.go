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
	Monitoring       *MonitoringConfig       `json:"monitoring"`
	InferenceAdapter *InferenceAdapterConfig `json:"inferenceadapter"`
}

// MonitoringConfig defines the Monitoring configuration
// +k8s:openapi-gen=false
type MonitoringConfig struct {
	Stats    []StatConfig    `json:"stats"`
	Outliers []OutlierConfig `json:"outliers"`
	Drift    []DriftConfig   `json:"drift"`
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

// InferenceAdapterConfig defines the configuration for the InferenceAdapter service
// +k8s:openapi-gen=false
type InferenceAdapterConfig struct {
	ContainerImage string `json:"containerimage"`
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

	inferenceAdapterConfig, err := getInferenceAdapterConfig(configMap)
	if err != nil {
		return nil, err
	}

	return &ModelMonitorConfig{
		Monitoring:       monitoringConfig,
		InferenceAdapter: inferenceAdapterConfig,
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

func getInferenceAdapterConfig(configMap *corev1.ConfigMap) (*InferenceAdapterConfig, error) {
	inferenceAdapterConfig := &InferenceAdapterConfig{}
	key := constants.InferenceAdapter.String()

	if data, ok := configMap.Data[key]; ok {
		err := json.Unmarshal([]byte(data), &inferenceAdapterConfig)
		if err != nil {
			return nil, fmt.Errorf("Unable to unmarshall %v json string due to %v ", key, err)
		}
	}
	return inferenceAdapterConfig, nil
}
