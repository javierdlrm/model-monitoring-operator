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
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	"knative.dev/serving/pkg/apis/autoscaling"
	knservingv1 "knative.dev/serving/pkg/apis/serving/v1"
)

var customizableServiceAnnotations = []string{
	autoscaling.MinScaleAnnotationKey,
	autoscaling.MaxScaleAnnotationKey,
	autoscaling.ClassAnnotationKey,
	autoscaling.MetricAnnotationKey,
	autoscaling.TargetAnnotationKey,
	autoscaling.TargetUtilizationPercentageKey,
	autoscaling.WindowAnnotationKey,
	autoscaling.PanicWindowPercentageAnnotationKey,
	autoscaling.PanicThresholdPercentageAnnotationKey,
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
	annotations, err := b.buildAnnotations(metadata, inferenceLoggerSpec)
	if err != nil {
		return nil, err
	}

	// Resources
	resources, err := b.buildResources(metadata, inferenceLoggerSpec)
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
									Resources: resources,
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

	annotations := utils.Filter(metadata.Annotations, func(key string) bool {
		return !utils.Includes(customizableServiceAnnotations, key)
	})

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
	if spec.Target == 0 {
		annotations[autoscaling.TargetAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultScalingTarget)
	} else {
		annotations[autoscaling.TargetAnnotationKey] = strconv.Itoa(spec.Target)
	}

	// Target utilization
	if spec.TargetUtilization == "" {
		annotations[autoscaling.TargetUtilizationPercentageKey] = constants.InferenceLoggerDefaultTargetUtilizationPercentage
	} else {
		annotations[autoscaling.TargetUtilizationPercentageKey] = spec.TargetUtilization
	}

	// Window
	if spec.Window == "" {
		annotations[autoscaling.WindowAnnotationKey] = constants.InferenceLoggerDefaultWindow
	} else {
		annotations[autoscaling.WindowAnnotationKey] = spec.Window
	}

	// Panic window
	if spec.PanicWindow == "" {
		annotations[autoscaling.PanicWindowPercentageAnnotationKey] = constants.InferenceLoggerDefaultPanicWindow
	} else {
		annotations[autoscaling.PanicWindowPercentageAnnotationKey] = spec.PanicWindow
	}

	// Panic threshold
	if spec.PanicThreshold == "" {
		annotations[autoscaling.PanicThresholdPercentageAnnotationKey] = constants.InferenceLoggerDefaultPanicThreshold
	} else {
		annotations[autoscaling.PanicThresholdPercentageAnnotationKey] = spec.PanicThreshold
	}

	// Min replicas
	if spec.MinScale == 0 {
		annotations[autoscaling.MinScaleAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultMinScale)
	} else {
		annotations[autoscaling.MinScaleAnnotationKey] = strconv.Itoa(spec.MinScale)
	}

	// Max replicas
	if spec.MaxScale == 0 {
		annotations[autoscaling.MaxScaleAnnotationKey] = fmt.Sprint(constants.InferenceLoggerDefaultMaxScale)
	} else {
		annotations[autoscaling.MaxScaleAnnotationKey] = strconv.Itoa(spec.MaxScale)
	}

	return annotations, nil
}

func (b *InferenceLoggerBuilder) buildResources(metadata metav1.ObjectMeta, spec monitoringv1beta1.InferenceLoggerSpec) (corev1.ResourceRequirements, error) {

	defaultResources := corev1.ResourceList{
		corev1.ResourceCPU:    resource.MustParse(constants.InferenceLoggerDefaultCPU),
		corev1.ResourceMemory: resource.MustParse(constants.InferenceLoggerDefaultMemory),
	}

	if spec.Resources.Requests == nil {
		spec.Resources.Requests = defaultResources
	} else {
		for name, value := range defaultResources {
			if _, ok := spec.Resources.Requests[name]; !ok {
				spec.Resources.Requests[name] = value
			}
		}
	}

	if spec.Resources.Limits == nil {
		spec.Resources.Limits = defaultResources
	} else {
		for name, value := range defaultResources {
			if _, ok := spec.Resources.Limits[name]; !ok {
				spec.Resources.Limits[name] = value
			}
		}
	}

	return spec.Resources, nil
}
