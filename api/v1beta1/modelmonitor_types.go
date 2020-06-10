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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ModelMonitorSpec defines the desired state of ModelMonitor
type ModelMonitorSpec struct {
	//+required
	Model ModelSpec `json:"model"`
	//+required
	Monitoring MonitoringSpec `json:"monitoring"`
	//+required
	Storage StorageSpec `json:"storage"`
	//+optional
	Job JobSpec `json:"job,omitempty"`
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
	Version *int `json:"version,omitempty"`
	//+required
	Schemas ModelSchemasSpec `json:"schemas"`
}

// ModelSchemasSpec defines the inference schema of a model
type ModelSchemasSpec struct {
	//+required
	Request string `json:"request"`
	//+required
	Response string `json:"response"`

	//+required
	Instance string `json:"instance"`
	//+required
	Prediction string `json:"prediction"`
}

// MonitoringSpec defines the Monitoring settings
type MonitoringSpec struct {
	//+required
	Trigger TriggerSpec `json:"trigger"`
	//+required
	Stats StatSpec `json:"stats"`
	//+optional
	Baseline *BaselineSpec `json:"baseline,omitempty"`
	//+optional
	Outliers *OutlierSpec `json:"outliers,omitempty"`
	//+optional
	Drift *DriftSpec `json:"drift,omitempty"`
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
	//+optional
	Max *MaxSpec `json:"max,omitempty"`
	//+optional
	Min *MinSpec `json:"min,omitempty"`
	//+optional
	Count *CountSpec `json:"count,omitempty"`
	//+optional
	Sum *SumSpec `json:"sum,omitempty"`
	//+optional
	Pow2Sum *Pow2SumSpec `json:"pow2Sum,omitempty"`
	//+optional
	Distr *DistrSpec `json:"distr,omitempty"`

	//+optional
	Avg *AvgSpec `json:"avg,omitempty"`
	//+optional
	Mean *MeanSpec `json:"mean,omitempty"`
	//+optional
	Stddev *StddevSpec `json:"stddev,omitempty"`
	//+optional
	Perc *PercSpec `json:"perc,omitempty"`

	//+optional
	Cov *CovSpec `json:"cov,omitempty"`
	//+optional
	Corr *CorrSpec `json:"corr,omitempty"`
}

// MaxSpec defines a Max stat
type MaxSpec struct{}

// MinSpec defines a Min stat
type MinSpec struct{}

// CountSpec defines a Count stat
type CountSpec struct{}

// SumSpec defines a Sum stat
type SumSpec struct{}

// Pow2SumSpec defines a Pow2Sum stat
type Pow2SumSpec struct{}

// DistrSpec defines a Distribution
type DistrSpec struct {
	//+optional
	Bounds map[string][]string `json:"bounds,omitempty"`
	//+optional
	Binning Binning `json:"binning,omitempty"`
}

// Binning defines the Distribution binning algorithm
//+kubebuilder:validation:Enum=sturge
type Binning string

// AvgSpec defines an Avg
type AvgSpec struct{}

// MeanSpec defines a Mean
type MeanSpec struct{}

// StddevSpec defines a Standard deviation
type StddevSpec struct {
	//+optional
	//+kubebuilder:validation:Enum=sample;population
	Type string `json:"type,omitempty"`
}

// PercSpec defines Percentiles
type PercSpec struct {
	//+required
	Percentiles []string `json:"percentiles"`
	//+optional
	IQR bool `json:"iqr,omitempty"`
}

// CovSpec defines a Covariance
type CovSpec struct {
	//+optional
	//+kubebuilder:validation:Enum=sample;population
	Type string `json:"type,omitempty"`
}

// CorrSpec defines a Correlation
type CorrSpec struct {
	//+optional
	//+kubebuilder:validation:Enum=sample;population
	Type string `json:"type,omitempty"`
}

// BaselineSpec defines Baseline stats
type BaselineSpec struct {
	//+optional
	Descriptive string `json:"descriptive,omitempty"`
	//+optional
	Distributions string `json:"distributions,omitempty"`
}

// OutlierSpec defines an Outlier detector
type OutlierSpec struct {
	//+optional
	Descriptive []string `json:"descriptive,omitempty"`
}

// DriftSpec defines a Drift detector
type DriftSpec struct {
	//+optional
	Wasserstein *ThresholdBasedDriftSpec `json:"wasserstein,omitempty"`
	//+optional
	KullbackLeibler *ThresholdBasedDriftSpec `json:"kullbackLeibler,omitempty"`
	//+optional
	JensenShannon *ThresholdBasedDriftSpec `json:"jensenShannon,omitempty"`
}

// ThresholdBasedDriftSpec defines a threshold-based Drift detector
type ThresholdBasedDriftSpec struct {
	//+required
	Threshold string `json:"threshold"`
	//+optional
	ShowAll bool `json:"showAll,omitempty"`
}

// StorageSpec defines the Storage settings
type StorageSpec struct {
	//+required
	Inference SinkSpec `json:"inference"`
	//+required
	Analysis AnalysisSpec `json:"analysis"`
}

// AnalysisSpec defines the Analysis storage
type AnalysisSpec struct {
	//+required
	Stats SinkSpec `json:"stats"`
	//+optional
	Outliers *SinkSpec `json:"outliers,omitempty"`
	//+optional
	Drift *SinkSpec `json:"drift,omitempty"`
}

//SinkSpec defines the configuration of a Sink
type SinkSpec struct {
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
	Partitions int32 `json:"partitions,omitempty"`
	//+optional
	ReplicationFactor int16 `json:"replicationFactor,omitempty"`
}

//JobSpec defines the configuration for Monitoring job
type JobSpec struct {
	//+optional
	Timeout int `json:"timeout,omitempty"`

	// TODO: Add Driver and Executor specs
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

// ModelMonitorStatus defines the observed state of ModelMonitor
type ModelMonitorStatus struct {
	// TODO: Add Model Monitor status
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
