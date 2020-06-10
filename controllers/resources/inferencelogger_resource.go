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
	annotations, err := b.buildAnnotations(metadata, inferenceLoggerSpec)
	if err != nil {
		return nil, err
	}

	// Concurrency (scaling hard limit, default 0 means limitless)
	concurrency := int64(inferenceLoggerSpec.Target)

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

func (b *InferenceLoggerBuilder) buildAnnotations(metadata metav1.ObjectMeta, spec monitoringv1beta1.InferenceLoggerSpec) (map[string]string, error) {

	annotations := metadata.Annotations

	// Autoscaler
	if spec.Autoscaler == "" {
		annotations[autoscaling.ClassAnnotationKey] = constants.InferenceLoggerDefaultScalingClass
	} else {
		annotations[autoscaling.ClassAnnotationKey] = string(spec.Autoscaler)
	}

	// Metric
	if spec.Metric == "" {
		annotations[autoscaling.MetricAnnotationKey] = constants.InferenceLoggerDefaultScalingMetric
	} else {
		annotations[autoscaling.MetricAnnotationKey] = string(spec.Metric)
	}

	// Target
	if spec.Target != 0 {
		annotations[autoscaling.TargetAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultScalingTarget)
	} else {
		annotations[autoscaling.TargetAnnotationKey] = strconv.Itoa(spec.Target)
	}

	// Target utilization
	if spec.TargetUtilization == 0 {
		annotations[autoscaling.TargetUtilizationPercentageKey] = fmt.Sprint(constants.InferenceLoggerDefaultTargetUtilizationPercentage)
	} else {
		annotations[autoscaling.TargetUtilizationPercentageKey] = strconv.Itoa(spec.TargetUtilization)
	}

	// Window
	if spec.Window == "" {
		annotations[autoscaling.WindowAnnotationKey] = constants.InferenceLoggerDefaultWindow
	} else {
		annotations[autoscaling.WindowAnnotationKey] = spec.Window
	}

	// Panic window
	if spec.PanicWindow == 0 {
		annotations[autoscaling.PanicWindowPercentageAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultPanicWindow)
	} else {
		annotations[autoscaling.PanicWindowPercentageAnnotationKey] = strconv.Itoa(spec.PanicWindow)
	}

	// Panic threshold
	if spec.PanicThreshold == 0 {
		annotations[autoscaling.PanicThresholdPercentageAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultPanicThreshold)
	} else {
		annotations[autoscaling.PanicThresholdPercentageAnnotationKey] = strconv.Itoa(spec.PanicThreshold)
	}

	// Min replicas
	if spec.MinReplicas == 0 {
		annotations[autoscaling.MinScaleAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultMinReplicas)
	} else {
		annotations[autoscaling.MinScaleAnnotationKey] = strconv.Itoa(spec.MinReplicas)
	}

	// Max replicas
	if spec.MaxReplicas == 0 {
		annotations[autoscaling.MaxScaleAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultMaxReplicas)
	} else {
		annotations[autoscaling.MaxScaleAnnotationKey] = strconv.Itoa(spec.MaxReplicas)
	}

	return annotations, nil
}
