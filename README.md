# Model Monitoring Operator
Kubernetes operator for ML model monitoring over KFServing, with Kafka, Spark and the [Model Monitoring framework](https://github.com/javierdlrm/model-monitoring)

## How can I start?

1. Install the Model Monitoring operator by choosing one of the versions available in 'install' folder.
`kubectl create -f install/v1beta1/model-monitoring.yaml`

2. Define and apply a Model Monitor resource (check 'config/samples/monitoring_v1beta1_modelmonitor.yaml' as an example)
`kubectl apply -f model-monitor.yaml`

3. Serve your model with KFServing specifying the following logger url:
'http://\<model-monitor-name\>-inferencelogger.\<namespace\>'

4. Check the inference analysis of your model by visiting the sinks specified in your Model Monitor.

## How it works?

The operator deploys an [Inference Logger](https://github.com/javierdlrm/inference-logger) to forward enriched inference logs to Kafka. Then it deploys a Spark job that consumes the corresponding Kafka topics and analysis the logs using a custom implementation of the [Model Monitoring framework](https://github.com/javierdlrm/model-monitoring).

## Monitoring Configuration

In order to see the available statistics, outliers and drift detectors check the [documentation](https://github.com/javierdlrm/model-monitoring) of the framework.

