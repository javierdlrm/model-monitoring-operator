
# Image details
IMG_NAME ?= javierdlrm/model-monitoring-operator
VERSION ?= v1beta1
IMG_V ?= ${IMG_NAME}:${VERSION}
IMG_L ?= ${IMG_NAME}:latest
# Mode
MODE ?= version
ifeq ($(MODE),latest)
IMG=${IMG_L}
else
IMG=${IMG_V}
endif
# Produce CRDs that work back to Kubernetes 1.11 (no version conversion)
CRD_OPTIONS ?= "crd:trivialVersions=true"

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

all: manager

# Run tests
test: generate fmt vet manifests
	# go test ./... -coverprofile cover.out
	# TODO: Pending. Currently kubebuilder/bin/etcd binaries are not compatible with WSL. (https://github.com/kubernetes-sigs/kubebuilder/issues/300)

# Build manager binary
manager: generate fmt vet
	go build -o bin/manager main.go

# Run against the configured Kubernetes cluster in ~/.kube/config
run: generate fmt vet manifests
	go run ./main.go

# Generates a yaml file for installing the operator
installer: manifests generate
	cd config/default/manager && kustomize edit set image controller=${IMG}
	kustomize build config/overlays/dev > install/${VERSION}/model-monitoring.yaml

# Install CRDs into a cluster
install: manifests
	kubectl create ns model-monitoring-system
	kustomize build config/default/crd | kubectl apply -f -
	kustomize build config/overlays/dev/configmap | kubectl apply -f -

# Uninstall CRDs from a cluster
uninstall: manifests
	kustomize build config/default/crd | kubectl delete -f -
	kustomize build config/overlays/dev/configmap | kubectl delete -f -
	kubectl delete ns model-monitoring-system

# Deploy controller in the configured Kubernetes cluster in ~/.kube/config
deploy: manifests installer
	kubectl apply -f install/${VERSION}/model-monitoring.yaml

# Remove deployment in the configured Kubernetes cluster in ~/.kube/config
undeploy: installer
	kubectl delete -f install/${VERSION}/model-monitoring.yaml

# Generate manifests e.g. CRD, RBAC etc.
manifests: controller-gen
	$(CONTROLLER_GEN) $(CRD_OPTIONS) rbac:roleName=manager-role webhook paths="./..." \
	output:crd:artifacts:config=config/default/crd/bases output:rbac:artifacts:config=config/default/rbac

# Run go fmt against code
fmt:
	go fmt ./...

# Run go vet against code
vet:
	go vet ./...

# Generate code
generate: controller-gen
	$(CONTROLLER_GEN) object:headerFile="hack/boilerplate.go.txt" paths="./..."

# Build and publish as docker image
docker: fmt vet docker-build docker-push

# Build the docker image
docker-build: test
	docker build . -t ${IMG_V} -t ${IMG_L}

# Push the docker image
docker-push:
	docker push ${IMG_V}
	docker push ${IMG_L}

# find or download controller-gen
# download controller-gen if necessary
controller-gen:
ifeq (, $(shell which controller-gen))
	@{ \
	set -e ;\
	CONTROLLER_GEN_TMP_DIR=$$(mktemp -d) ;\
	cd $$CONTROLLER_GEN_TMP_DIR ;\
	go mod init tmp ;\
	go get sigs.k8s.io/controller-tools/cmd/controller-gen@v0.2.5 ;\
	rm -rf $$CONTROLLER_GEN_TMP_DIR ;\
	}
CONTROLLER_GEN=$(GOBIN)/controller-gen
else
CONTROLLER_GEN=$(shell which controller-gen)
endif
