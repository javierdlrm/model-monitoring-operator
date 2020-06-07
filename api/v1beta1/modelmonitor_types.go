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
	//+required
	Monitoring MonitoringSpec `json:"monitoring"`
	//+required
	Job JobSpec `json:"job"`
	//+optional
	InferenceLogger InferenceLoggerSpec `json:"inferenceLogger,omitempty"`
}

// ModelSpec defines the Model being monitored. It should match with KFserving inferenceservice name
type ModelSpec struct {
	//+required
	Name string `json:"name"`
	//+optional
	ID string `json:"id,omitempty"`
	//+optional
	Version int `json:"version,omitempty"`
	//+required
	Schemas ModelSchemasSpec `json:"schemas"`
}

// ModelSchemasSpec defines the inference schema of a model
type ModelSchemasSpec struct {
	//+required
	Request string `json:"request"`
	//+required
	Response string `json:"response"`
}

// MonitoringSpec defines the Monitoring settings
type MonitoringSpec struct {
	//+required
	Trigger TriggerSpec `json:"trigger"`
	//+required
	Stats []StatSpec `json:"stats"`
	//+optional
	Outliers []OutlierSpec `json:"outliers,omitempty"`
	//+optional
	Drift []DriftSpec `json:"drift,omitempty"`
}

// TriggerSpec defines the Monitoring trigger setting
type TriggerSpec struct {
	//+required
	Window WindowSpec `json:"window"`
}

// WindowSpec defines a Window as Monitoring job trigger
type WindowSpec struct {
	//+required
	Duration int `json:"duration"`
	//+required
	Slide int `json:"slide"`
	//+required
	WatermarkDelay int `json:"watermarkDelay"`
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
	Threshold resource.Quantity `json:"threshold,omitempty"`
	//+optional
	ShowAll bool `json:"showall,omitempty"`
}

// InferenceLoggerSpec defines the configuration for InferenceLogger Knative Service.
type InferenceLoggerSpec struct {
	// +optional
	MinReplicas int `json:"minReplicas,omitempty"`
	// +optional
	MaxReplicas int `json:"maxReplicas,omitempty"`
	// +optional
	Parallelism int `json:"parallelism,omitempty"`
}

//JobSpec defines the configuration for Monitoring job
type JobSpec struct {
	//+required
	Source SourceSpec `json:"source"`
	//+optional
	Sink []SinkSpec `json:"sink,omitempty"`
	//+optional
	Timeout int `json:"timeout,omitempty"`
}

//SourceSpec defines the configuration of the source for the Monitoring job
type SourceSpec struct {
	//+required
	Kafka KafkaSpec `json:"kafka"`
}

//SinkSpec defines the configuration of the sink for the Monitoring job
type SinkSpec struct {
	//+required
	//+kubebuilder:validation:Enum=stats;outliers;drift
	Pipe string `json:"pipe"`
	//+required
	Kafka KafkaSpec `json:"kafka"`
}

// KafkaSpec defines the KafkaTopic used for inference logging.
type KafkaSpec struct {
	//+required
	Brokers string `json:"brokers"`
	//+required
	Topic KafkaTopicSpec `json:"topic"`
}

// KafkaTopicSpec defines a Kafka topic
type KafkaTopicSpec struct {
	//+required
	Name string `json:"name"`
	//+optional
	Partitions int32 `json:"partitions"`
	//+optional
	ReplicationFactor int16 `json:"replicationFactor"`
}

// ModelMonitorStatus defines the observed state of ModelMonitor
type ModelMonitorStatus struct {
	// TODO: Add InferenceLogger service name
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
