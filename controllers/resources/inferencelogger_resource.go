package resources

import (
	"fmt"
	"strconv"

	monitoringv1beta1 "github.com/javierdlrm/model-monitoring-operator/api/v1beta1"
	"github.com/javierdlrm/model-monitoring-operator/constants"
	typesutils "github.com/javierdlrm/model-monitoring-operator/utils"

	"github.com/kubeflow/kfserving/pkg/utils"

	"github.com/go-logr/logr"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"knative.dev/serving/pkg/apis/autoscaling"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

var serviceAnnotationDisallowedList = []string{
	autoscaling.MinScaleAnnotationKey,
	autoscaling.MaxScaleAnnotationKey,
	"kubectl.kubernetes.io/last-applied-configuration",
}

// InferenceLoggerBuilder defines the builder for InferenceLogger
type InferenceLoggerBuilder struct {
	ModelMonitorConfig *monitoringv1beta1.ModelMonitorConfig
	Log                logr.Logger
}

// NewInferenceLoggerBuilder creates an InferenceLogger builder
func NewInferenceLoggerBuilder(config *corev1.ConfigMap, log logr.Logger) *InferenceLoggerBuilder {
	modelMonitorConfig, err := monitoringv1beta1.NewModelMonitorConfig(config)
	if err != nil {
		fmt.Printf("Failed to get model monitor config %s", err.Error())
		panic("Failed to get model monitor config")
	}
	return &InferenceLoggerBuilder{
		ModelMonitorConfig: modelMonitorConfig,
		Log:                log,
	}
}

// CreateInferenceLoggerService creates the Knative Service for InferenceLogger
func (b *InferenceLoggerBuilder) CreateInferenceLoggerService(serviceName string, modelMonitor *monitoringv1beta1.ModelMonitor) (*knservingv1.Service, error) {

	// Specs
	metadata := modelMonitor.ObjectMeta
	inferenceLoggerSpec := modelMonitor.Spec.InferenceLogger
	inferenceSpec := modelMonitor.Spec.Storage.Inference

	// Autoscaling annotations
	annotations, err := b.buildAnnotations(metadata, inferenceLoggerSpec.MinReplicas, inferenceLoggerSpec.MaxReplicas, inferenceLoggerSpec.Parallelism)
	if err != nil {
		return nil, err
	}

	// Concurrency
	concurrency := int64(inferenceLoggerSpec.Parallelism)

	// Knative Service
	service := &knservingv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName,
			Namespace: metadata.Namespace,
			Labels:    metadata.Labels,
		},
		Spec: knservingv1.ServiceSpec{
			ConfigurationSpec: knservingv1.ConfigurationSpec{
				Template: knservingv1.RevisionTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: utils.Union(metadata.Labels, map[string]string{
							constants.InferenceLoggerModelLabel: metadata.Name,
						}),
						Annotations: annotations,
					},
					Spec: knservingv1.RevisionSpec{
						TimeoutSeconds:       &constants.InferenceLoggerDefaultTimeout,
						ContainerConcurrency: &concurrency,
						PodSpec: corev1.PodSpec{
							Containers: []corev1.Container{
								{
									Image:           b.ModelMonitorConfig.InferenceLogger.ContainerImage,
									Name:            constants.ModelMonitorContainerName,
									ImagePullPolicy: corev1.PullAlways,
									Env: []corev1.EnvVar{
										corev1.EnvVar{
											Name:  constants.InferenceLoggerEnvKafkaBrokersLabel,
											Value: inferenceSpec.Kafka.Brokers,
										},
										corev1.EnvVar{
											Name:  constants.InferenceLoggerEnvKafkaTopicLabel,
											Value: inferenceSpec.Kafka.Topic.Name,
										},
										corev1.EnvVar{
											Name:  constants.InferenceLoggerEnvKafkaTopicPartitionsLabel,
											Value: typesutils.String32(inferenceSpec.Kafka.Topic.Partitions),
										},
										corev1.EnvVar{
											Name:  constants.InferenceLoggerEnvKafkaTopicReplicationFactorLabel,
											Value: typesutils.String16(inferenceSpec.Kafka.Topic.ReplicationFactor),
										},
									},
									ReadinessProbe: &corev1.Probe{
										Handler: corev1.Handler{
											TCPSocket: &corev1.TCPSocketAction{
												Port: intstr.FromInt(0),
											},
										},
										SuccessThreshold: 1,
									},
								},
							},
						},
					},
				},
			},
		},
	}

	return service, nil
}

func (b *InferenceLoggerBuilder) buildAnnotations(metadata metav1.ObjectMeta, minReplicas int, maxReplicas int, parallelism int) (map[string]string, error) {

	annotations := utils.Filter(metadata.Annotations, func(key string) bool {
		return !utils.Includes(serviceAnnotationDisallowedList, key)
	})

	if minReplicas == 0 {
		annotations[autoscaling.MinScaleAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultMinReplicas)
	} else if minReplicas != 0 {
		annotations[autoscaling.MinScaleAnnotationKey] = fmt.Sprint(minReplicas)
	}

	if maxReplicas != 0 {
		annotations[autoscaling.MaxScaleAnnotationKey] = fmt.Sprint(maxReplicas)
	}

	// User can pass down scaling target annotation to overwrite the target default 1
	if _, ok := annotations[autoscaling.TargetAnnotationKey]; !ok {
		if parallelism == 0 {
			annotations[autoscaling.TargetAnnotationKey] = constants.InferenceLoggerDefaultScalingTarget
		} else {
			annotations[autoscaling.TargetAnnotationKey] = strconv.Itoa(parallelism)
		}
	}
	// User can pass down scaling class annotation to overwrite the default scaling KPA
	if _, ok := annotations[autoscaling.ClassAnnotationKey]; !ok {
		annotations[autoscaling.ClassAnnotationKey] = autoscaling.KPA
	}
	return annotations, nil
}
