// +build !ignore_autogenerated

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

// Code generated by controller-gen. DO NOT EDIT.

package v1beta1

import (
	runtime "k8s.io/apimachinery/pkg/runtime"
)

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DriftConfig) DeepCopyInto(out *DriftConfig) {
	*out = *in
	out.Threshold = in.Threshold.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DriftConfig.
func (in *DriftConfig) DeepCopy() *DriftConfig {
	if in == nil {
		return nil
	}
	out := new(DriftConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *DriftSpec) DeepCopyInto(out *DriftSpec) {
	*out = *in
	out.Threshold = in.Threshold.DeepCopy()
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new DriftSpec.
func (in *DriftSpec) DeepCopy() *DriftSpec {
	if in == nil {
		return nil
	}
	out := new(DriftSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InferenceLoggerConfig) DeepCopyInto(out *InferenceLoggerConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InferenceLoggerConfig.
func (in *InferenceLoggerConfig) DeepCopy() *InferenceLoggerConfig {
	if in == nil {
		return nil
	}
	out := new(InferenceLoggerConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *InferenceLoggerSpec) DeepCopyInto(out *InferenceLoggerSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new InferenceLoggerSpec.
func (in *InferenceLoggerSpec) DeepCopy() *InferenceLoggerSpec {
	if in == nil {
		return nil
	}
	out := new(InferenceLoggerSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaConfig) DeepCopyInto(out *KafkaConfig) {
	*out = *in
	if in.Topic != nil {
		in, out := &in.Topic, &out.Topic
		*out = new(KafkaTopicConfig)
		**out = **in
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaConfig.
func (in *KafkaConfig) DeepCopy() *KafkaConfig {
	if in == nil {
		return nil
	}
	out := new(KafkaConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaSpec) DeepCopyInto(out *KafkaSpec) {
	*out = *in
	out.Topic = in.Topic
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaSpec.
func (in *KafkaSpec) DeepCopy() *KafkaSpec {
	if in == nil {
		return nil
	}
	out := new(KafkaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaTopicConfig) DeepCopyInto(out *KafkaTopicConfig) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaTopicConfig.
func (in *KafkaTopicConfig) DeepCopy() *KafkaTopicConfig {
	if in == nil {
		return nil
	}
	out := new(KafkaTopicConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *KafkaTopicSpec) DeepCopyInto(out *KafkaTopicSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new KafkaTopicSpec.
func (in *KafkaTopicSpec) DeepCopy() *KafkaTopicSpec {
	if in == nil {
		return nil
	}
	out := new(KafkaTopicSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelMonitor) DeepCopyInto(out *ModelMonitor) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ObjectMeta.DeepCopyInto(&out.ObjectMeta)
	in.Spec.DeepCopyInto(&out.Spec)
	out.Status = in.Status
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelMonitor.
func (in *ModelMonitor) DeepCopy() *ModelMonitor {
	if in == nil {
		return nil
	}
	out := new(ModelMonitor)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelMonitor) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelMonitorConfig) DeepCopyInto(out *ModelMonitorConfig) {
	*out = *in
	if in.Monitoring != nil {
		in, out := &in.Monitoring, &out.Monitoring
		*out = new(MonitoringConfig)
		(*in).DeepCopyInto(*out)
	}
	if in.InferenceLogger != nil {
		in, out := &in.InferenceLogger, &out.InferenceLogger
		*out = new(InferenceLoggerConfig)
		**out = **in
	}
	if in.Kafka != nil {
		in, out := &in.Kafka, &out.Kafka
		*out = new(KafkaConfig)
		(*in).DeepCopyInto(*out)
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelMonitorConfig.
func (in *ModelMonitorConfig) DeepCopy() *ModelMonitorConfig {
	if in == nil {
		return nil
	}
	out := new(ModelMonitorConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelMonitorList) DeepCopyInto(out *ModelMonitorList) {
	*out = *in
	out.TypeMeta = in.TypeMeta
	in.ListMeta.DeepCopyInto(&out.ListMeta)
	if in.Items != nil {
		in, out := &in.Items, &out.Items
		*out = make([]ModelMonitor, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelMonitorList.
func (in *ModelMonitorList) DeepCopy() *ModelMonitorList {
	if in == nil {
		return nil
	}
	out := new(ModelMonitorList)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyObject is an autogenerated deepcopy function, copying the receiver, creating a new runtime.Object.
func (in *ModelMonitorList) DeepCopyObject() runtime.Object {
	if c := in.DeepCopy(); c != nil {
		return c
	}
	return nil
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelMonitorSpec) DeepCopyInto(out *ModelMonitorSpec) {
	*out = *in
	out.Model = in.Model
	in.Monitoring.DeepCopyInto(&out.Monitoring)
	out.InferenceLogger = in.InferenceLogger
	out.Kafka = in.Kafka
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelMonitorSpec.
func (in *ModelMonitorSpec) DeepCopy() *ModelMonitorSpec {
	if in == nil {
		return nil
	}
	out := new(ModelMonitorSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelMonitorStatus) DeepCopyInto(out *ModelMonitorStatus) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelMonitorStatus.
func (in *ModelMonitorStatus) DeepCopy() *ModelMonitorStatus {
	if in == nil {
		return nil
	}
	out := new(ModelMonitorStatus)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelSchemaSpec) DeepCopyInto(out *ModelSchemaSpec) {
	*out = *in
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelSchemaSpec.
func (in *ModelSchemaSpec) DeepCopy() *ModelSchemaSpec {
	if in == nil {
		return nil
	}
	out := new(ModelSchemaSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *ModelSpec) DeepCopyInto(out *ModelSpec) {
	*out = *in
	out.Schema = in.Schema
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new ModelSpec.
func (in *ModelSpec) DeepCopy() *ModelSpec {
	if in == nil {
		return nil
	}
	out := new(ModelSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringConfig) DeepCopyInto(out *MonitoringConfig) {
	*out = *in
	if in.Stats != nil {
		in, out := &in.Stats, &out.Stats
		*out = make([]*StatConfig, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(StatConfig)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Outliers != nil {
		in, out := &in.Outliers, &out.Outliers
		*out = make([]*OutlierConfig, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(OutlierConfig)
				(*in).DeepCopyInto(*out)
			}
		}
	}
	if in.Drift != nil {
		in, out := &in.Drift, &out.Drift
		*out = make([]*DriftConfig, len(*in))
		for i := range *in {
			if (*in)[i] != nil {
				in, out := &(*in)[i], &(*out)[i]
				*out = new(DriftConfig)
				(*in).DeepCopyInto(*out)
			}
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringConfig.
func (in *MonitoringConfig) DeepCopy() *MonitoringConfig {
	if in == nil {
		return nil
	}
	out := new(MonitoringConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *MonitoringSpec) DeepCopyInto(out *MonitoringSpec) {
	*out = *in
	if in.Stats != nil {
		in, out := &in.Stats, &out.Stats
		*out = make([]StatSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Outliers != nil {
		in, out := &in.Outliers, &out.Outliers
		*out = make([]OutlierSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if in.Drift != nil {
		in, out := &in.Drift, &out.Drift
		*out = make([]DriftSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new MonitoringSpec.
func (in *MonitoringSpec) DeepCopy() *MonitoringSpec {
	if in == nil {
		return nil
	}
	out := new(MonitoringSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutlierConfig) DeepCopyInto(out *OutlierConfig) {
	*out = *in
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutlierConfig.
func (in *OutlierConfig) DeepCopy() *OutlierConfig {
	if in == nil {
		return nil
	}
	out := new(OutlierConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *OutlierSpec) DeepCopyInto(out *OutlierSpec) {
	*out = *in
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new OutlierSpec.
func (in *OutlierSpec) DeepCopy() *OutlierSpec {
	if in == nil {
		return nil
	}
	out := new(OutlierSpec)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatConfig) DeepCopyInto(out *StatConfig) {
	*out = *in
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatConfig.
func (in *StatConfig) DeepCopy() *StatConfig {
	if in == nil {
		return nil
	}
	out := new(StatConfig)
	in.DeepCopyInto(out)
	return out
}

// DeepCopyInto is an autogenerated deepcopy function, copying the receiver, writing into out. in must be non-nil.
func (in *StatSpec) DeepCopyInto(out *StatSpec) {
	*out = *in
	if in.Params != nil {
		in, out := &in.Params, &out.Params
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
}

// DeepCopy is an autogenerated deepcopy function, copying the receiver, creating a new StatSpec.
func (in *StatSpec) DeepCopy() *StatSpec {
	if in == nil {
		return nil
	}
	out := new(StatSpec)
	in.DeepCopyInto(out)
	return out
}
