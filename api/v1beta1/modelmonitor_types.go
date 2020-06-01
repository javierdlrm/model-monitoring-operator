/*
Copyright 2020 Javier de la Rúa Martínez.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1beta1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ModelMonitorSpec defines the desired state of ModelMonitor
type ModelMonitorSpec struct {
	//+required
	Model ModelSpec `json:"model"`
	//+optional
	Monitoring MonitoringSpec `json:"monitoring"`
	//+optional
	InferenceAdapter InferenceAdapterSpec `json:"inferenceadapter"`
	//+required
	Kafka KafkaSpec `json:"kafka"`
}

// ModelSpec defines the Model being monitored. It should match with KFserving inferenceservice name
type ModelSpec struct {
	//+required
	Name string `json:"name"`
	//+optional
	ID string `json:"id,omitempty"`
	//+optional
	Version int `json:"version,omitempty"`
}

// MonitoringSpec defines the Monitoring settings
type MonitoringSpec struct {
	//+optional
	Stats []StatSpec `json:"stats,omitempty"`
	//+optional
	Outliers []OutlierSpec `json:"outliers,omitempty"`
	//+optional
	Drift []DriftSpec `json:"drift,omitempty"`
}

// StatSpec defines a Statistic
type StatSpec struct {
	//+required
	Name string `json:"name"`
	//+optional
	Params map[string]string `json:"params,omitempty"`
}

// OutlierSpec defines an Outlier detector
type OutlierSpec struct {
	//+required
	Name string `json:"name"`
	//+optional
	Params map[string]string `json:"params,omitempty"`
}

// DriftSpec defines a Drift detector
type DriftSpec struct {
	//+required
	Name string `json:"name"`
	//+optional
	Threshold resource.Quantity `json:"threshold"`
	//+optional
	ShowAll bool `json:"showall"`
}

// InferenceAdapterSpec defines the configuration for InferenceAdapter Knative Service.
type InferenceAdapterSpec struct {
	// +optional
	MinReplicas int `json:"minReplicas,omitempty"`
	// +optional
	MaxReplicas int `json:"maxReplicas,omitempty"`
	// +optional
	Parallelism int `json:"parallelism,omitempty"`
}

// KafkaSpec defines the KafkaTopic used for inference logging.
type KafkaSpec struct {
	//+required
	Brokers string `json:"brokers"`
	// TODO: Get brokers using Namespace and ClusterName (needed for Kafka Topic creation)
}

// ModelMonitorStatus defines the observed state of ModelMonitor
type ModelMonitorStatus struct {
	// TODO: Add InferenceAdapter service name
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:resource:path=modelmonitors,shortName=modelmonitor

// ModelMonitor is the Schema for the modelmonitors API
type ModelMonitor struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ModelMonitorSpec   `json:"spec,omitempty"`
	Status ModelMonitorStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ModelMonitorList contains a list of ModelMonitor
type ModelMonitorList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ModelMonitor `json:"items"`
}

func init() {
	SchemeBuilder.Register(&ModelMonitor{}, &ModelMonitorList{})
}
