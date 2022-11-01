## Copyright Istio Authors
##
## Licensed under the Apache License, Version 2.0 (the "License");
## you may not use this file except in compliance with the License.
## You may obtain a copy of the License at
##
##     http://www.apache.org/licenses/LICENSE-2.0
##
## Unless required by applicable law or agreed to in writing, software
## distributed under the License is distributed on an "AS IS" BASIS,
## WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
## See the License for the specific language governing permissions and
## limitations under the License.

#-----------------------------------------------------------------------------
# Global Variables
#-----------------------------------------------------------------------------
ISTIO_GO := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
export ISTIO_GO
SHELL := /bin/bash -o pipefail

export VERSION ?= 1.14-dev

# Base version of Istio image to use
BASE_VERSION ?= latest

export GO111MODULE ?= on
export GOPROXY ?= https://proxy.golang.org
export GOSUMDB ?= sum.golang.org

# If GOPATH is not set by the env, set it to a sane value
GOPATH ?= $(shell cd ${ISTIO_GO}/../../..; pwd)
export GOPATH

# If GOPATH is made up of several paths, use the first one for our targets in this Makefile
GO_TOP := $(shell echo ${GOPATH} | cut -d ':' -f1)
export GO_TOP

GO ?= go

GOARCH_LOCAL := $(TARGET_ARCH)
GOOS_LOCAL := $(TARGET_OS)

#-----------------------------------------------------------------------------
# Output control
#-----------------------------------------------------------------------------
# Invoke make VERBOSE=1 to enable echoing of the command being executed
export VERBOSE ?= 0
# Place the variable Q in front of a command to control echoing of the command being executed.
Q = $(if $(filter 1,$VERBOSE),,@)
# Use the variable H to add a header (equivalent to =>) to informational output
H = $(shell printf "\033[34;1m=>\033[0m")

ifeq ($(origin DEBUG), undefined)
  BUILDTYPE_DIR:=release
else ifeq ($(DEBUG),0)
  BUILDTYPE_DIR:=release
else
  BUILDTYPE_DIR:=debug
  export GCFLAGS:=all=-N -l
  $(info $(H) Build with debugger information)
endif

# Optional file including user-specific settings (HUB, TAG, etc)
-include .istiorc.mk

# Environment for tests, the directory containing istio and deps binaries.
# Typically same as GOPATH/bin, so tests work seemlessly with IDEs.

export ISTIO_BIN=$(GOBIN)
# Using same package structure as pkg/

# If we are running in the Linux build container on non Linux hosts, we add the
# linux binaries to the build dependencies, BUILD_DEPS, which can be added to other targets
# that would need the Linux binaries (ex. tests).
BUILD_DEPS:=
ifeq ($(IN_BUILD_CONTAINER),1)
  ifneq ($(GOOS_LOCAL),"linux")
    BUILD_DEPS += build-linux
  endif
endif

export ARTIFACTS ?= $(TARGET_OUT)
export JUNIT_OUT ?= $(ARTIFACTS)/junit.xml
export REPO_ROOT := $(shell git rev-parse --show-toplevel)

# Make directories needed by the build system
$(shell mkdir -p $(TARGET_OUT_LINUX))
$(shell mkdir -p $(TARGET_OUT_LINUX)/logs)
$(shell mkdir -p $(dir $(JUNIT_OUT)))

# Need seperate target for init:
$(TARGET_OUT):
	@mkdir -p $@

# scratch dir: this shouldn't be simply 'docker' since that's used for docker.save to store tar.gz files
ISTIO_DOCKER:=${TARGET_OUT_LINUX}/docker_temp

# scratch dir for building isolated images. Please don't remove it again - using
# ISTIO_DOCKER results in slowdown, all files (including multiple copies of envoy) will be
# copied to the docker temp container - even if you add only a tiny file, >1G of data will
# be copied, for each docker image.
DOCKER_BUILD_TOP:=${TARGET_OUT_LINUX}/docker_build
DOCKERX_BUILD_TOP:=${TARGET_OUT_LINUX}/dockerx_build

# dir where tar.gz files from docker.save are stored
ISTIO_DOCKER_TAR:=${TARGET_OUT_LINUX}/release/docker

# Populate the git version for istio/proxy (i.e. Envoy)
#ifeq ($(PROXY_REPO_SHA),)
#  export PROXY_REPO_SHA:=$(shell grep PROXY_REPO_SHA istio.deps  -A 4 | grep lastStableSHA | cut -f 4 -d '"')
#endif

# Envoy binary variables Keep the default URLs up-to-date with the latest push from istio/proxy.

export ISTIO_ENVOY_BASE_URL ?= https://storage.googleapis.com/istio-build/proxy

# Use envoy as the sidecar by default
export SIDECAR ?= envoy

# OS-neutral vars. These currently only work for linux.
export ISTIO_ENVOY_VERSION ?= ${PROXY_REPO_SHA}
export ISTIO_ENVOY_DEBUG_URL ?= $(ISTIO_ENVOY_BASE_URL)/envoy-debug-$(ISTIO_ENVOY_VERSION).tar.gz
export ISTIO_ENVOY_CENTOS_DEBUG_URL ?= $(ISTIO_ENVOY_BASE_URL)/envoy-centos-debug-$(ISTIO_ENVOY_VERSION).tar.gz
export ISTIO_ENVOY_RELEASE_URL ?= $(ISTIO_ENVOY_BASE_URL)/envoy-alpha-$(ISTIO_ENVOY_VERSION).tar.gz
export ISTIO_ENVOY_CENTOS_RELEASE_URL ?= $(ISTIO_ENVOY_BASE_URL)/envoy-centos-alpha-$(ISTIO_ENVOY_VERSION).tar.gz

# Envoy Linux vars.
export ISTIO_ENVOY_LINUX_VERSION ?= ${ISTIO_ENVOY_VERSION}
export ISTIO_ENVOY_LINUX_DEBUG_URL ?= ${ISTIO_ENVOY_DEBUG_URL}
export ISTIO_ENVOY_LINUX_RELEASE_URL ?= ${ISTIO_ENVOY_RELEASE_URL}
# Variables for the extracted debug/release Envoy artifacts.
export ISTIO_ENVOY_LINUX_DEBUG_DIR ?= ${TARGET_OUT_LINUX}/debug
export ISTIO_ENVOY_LINUX_DEBUG_NAME ?= envoy-debug-${ISTIO_ENVOY_LINUX_VERSION}
export ISTIO_ENVOY_LINUX_DEBUG_PATH ?= ${ISTIO_ENVOY_LINUX_DEBUG_DIR}/${ISTIO_ENVOY_LINUX_DEBUG_NAME}
export ISTIO_ENVOY_CENTOS_LINUX_DEBUG_NAME ?= envoy-centos-debug-${ISTIO_ENVOY_LINUX_VERSION}
export ISTIO_ENVOY_CENTOS_LINUX_DEBUG_PATH ?= ${ISTIO_ENVOY_LINUX_DEBUG_DIR}/${ISTIO_ENVOY_CENTOS_LINUX_DEBUG_NAME}

export ISTIO_ENVOY_LINUX_RELEASE_DIR ?= ${TARGET_OUT_LINUX}/release
export ISTIO_ENVOY_LINUX_RELEASE_NAME ?= ${SIDECAR}-${ISTIO_ENVOY_VERSION}
export ISTIO_ENVOY_LINUX_RELEASE_PATH ?= ${ISTIO_ENVOY_LINUX_RELEASE_DIR}/${ISTIO_ENVOY_LINUX_RELEASE_NAME}
export ISTIO_ENVOY_CENTOS_LINUX_RELEASE_NAME ?= envoy-centos-${ISTIO_ENVOY_LINUX_VERSION}
export ISTIO_ENVOY_CENTOS_LINUX_RELEASE_PATH ?= ${ISTIO_ENVOY_LINUX_RELEASE_DIR}/${ISTIO_ENVOY_CENTOS_LINUX_RELEASE_NAME}

# Envoy macOS vars.
# TODO Change url when official envoy release for macOS is available
export ISTIO_ENVOY_MACOS_VERSION ?= 1.0.2
export ISTIO_ENVOY_MACOS_RELEASE_URL ?= https://github.com/istio/proxy/releases/download/${ISTIO_ENVOY_MACOS_VERSION}/istio-proxy-${ISTIO_ENVOY_MACOS_VERSION}-macos.tar.gz
# Variables for the extracted debug/release Envoy artifacts.
export ISTIO_ENVOY_MACOS_RELEASE_DIR ?= ${TARGET_OUT}/release
export ISTIO_ENVOY_MACOS_RELEASE_NAME ?= envoy-${ISTIO_ENVOY_MACOS_VERSION}
export ISTIO_ENVOY_MACOS_RELEASE_PATH ?= ${ISTIO_ENVOY_MACOS_RELEASE_DIR}/${ISTIO_ENVOY_MACOS_RELEASE_NAME}

# Allow user-override envoy bootstrap config path.
export ISTIO_ENVOY_BOOTSTRAP_CONFIG_PATH ?= ${ISTIO_GO}/tools/packaging/common/envoy_bootstrap.json
export ISTIO_ENVOY_BOOTSTRAP_CONFIG_DIR = $(dir ${ISTIO_ENVOY_BOOTSTRAP_CONFIG_PATH})

# If the hub is not explicitly set, use default to istio.
HUB ?=istio
ifeq ($(HUB),)
  $(error "HUB cannot be empty")
endif

# For dockerx builds, allow HUBS which is a space seperated list of hubs. Default to HUB.
HUBS ?= $(HUB)

# If tag not explicitly set in users' .istiorc.mk or command line, default to the git sha.
TAG ?= $(shell git rev-parse --verify HEAD)
ifeq ($(TAG),)
  $(error "TAG cannot be empty")
endif

VARIANT :=
ifeq ($(VARIANT),)
  TAG_VARIANT:=${TAG}
else
  TAG_VARIANT:=${TAG}-${VARIANT}
endif

PULL_POLICY ?= IfNotPresent
ifeq ($(TAG),latest)
  PULL_POLICY = Always
endif
ifeq ($(PULL_POLICY),)
  $(error "PULL_POLICY cannot be empty")
endif

include tools/proto/proto.mk

.PHONY: default
default: init build test

.PHONY: init
# Downloads envoy, based on the SHA defined in the base pilot Dockerfile
#init: $(TARGET_OUT)/istio_is_init
#	@mkdir -p ${TARGET_OUT}/logs
#	@mkdir -p ${TARGET_OUT}/release

# I tried to make this dependent on what I thought was the appropriate
# lock file, but it caused the rule for that file to get run (which
# seems to be about obtaining a new version of the 3rd party libraries).
#$(TARGET_OUT)/istio_is_init: bin/init.sh istio.deps | $(TARGET_OUT)
#	@# Add a retry, as occasionally we see transient connection failures to GCS
#	@# Like `curl: (56) OpenSSL SSL_read: SSL_ERROR_SYSCALL, errno 104`
#	TARGET_OUT=$(TARGET_OUT) ISTIO_BIN=$(ISTIO_BIN) GOOS_LOCAL=$(GOOS_LOCAL) bin/retry.sh SSL_ERROR_SYSCALL bin/init.sh
#	touch $(TARGET_OUT)/istio_is_init

# init.sh downloads envoy and webassembly plugins
${TARGET_OUT}/${SIDECAR}: init
${ISTIO_ENVOY_LINUX_DEBUG_PATH}: init
${ISTIO_ENVOY_LINUX_RELEASE_PATH}: init
${ISTIO_ENVOY_MACOS_RELEASE_PATH}: init

# Pull dependencies, based on the checked in Gopkg.lock file.
# Developers must manually run `dep ensure` if adding new deps
depend: init | $(TARGET_OUT)

DIRS_TO_CLEAN := $(TARGET_OUT)
DIRS_TO_CLEAN += $(TARGET_OUT_LINUX)

$(OUTPUT_DIRS):
	@mkdir -p $@

.PHONY: ${GEN_CERT}
GEN_CERT := ${ISTIO_BIN}/generate_cert
${GEN_CERT}:
	GOOS=$(GOOS_LOCAL) && GOARCH=$(GOARCH_LOCAL) && common/scripts/gobuild.sh $@ ./security/tools/generate_cert

#-----------------------------------------------------------------------------
# Target: precommit
#-----------------------------------------------------------------------------
.PHONY: precommit format lint

# Target run by the pre-commit script, to automate formatting and lint
# If pre-commit script is not used, please run this manually.
precommit: format lint

format: fmt ## Auto formats all code. This should be run before sending a PR.

fmt: format-go format-python tidy-go

ifeq ($(DEBUG),1)
# gobuild script uses custom linker flag to set the variables.
RELEASE_LDFLAGS=''
else
RELEASE_LDFLAGS='-extldflags -static -s -w'
endif

# List of all binaries to build
# We split the binaries into "agent" binaries and standard ones. This corresponds to build "agent".
# This allows conditional compilation to avoid pulling in costly dependencies to the agent, such as XDS and k8s.
AGENT_BINARIES:=./pilot/cmd/pilot-agent
STANDARD_BINARIES:=./istioctl/cmd/istioctl \
  ./pilot/cmd/pilot-discovery 
#  ./pkg/test/echo/cmd/client \
#  ./pkg/test/echo/cmd/server \
#  ./samples/extauthz/cmd/extauthz \
#  ./operator/cmd/operator \
#  ./cni/cmd/istio-cni \
#  ./cni/cmd/istio-cni-taint \
#  ./cni/cmd/install-cni \
#  ./tools/istio-iptables \
#  ./tools/bug-report
BINARIES:=$(STANDARD_BINARIES) $(AGENT_BINARIES)

# List of binaries included in releases
RELEASE_BINARIES:=pilot-discovery pilot-agent istioctl bug-report

.PHONY: build
build: depend ## Builds all go binaries.
	GOPROXY=https://goproxy.cn GOOS=$(GOOS_LOCAL) GOARCH=$(GOARCH_LOCAL) LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $(TARGET_OUT)/ $(STANDARD_BINARIES)
	GOPROXY=https://goproxy.cn GOOS=$(GOOS_LOCAL) GOARCH=$(GOARCH_LOCAL) LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $(TARGET_OUT)/ -tags=agent $(AGENT_BINARIES)

# The build-linux target is responsible for building binaries used within containers.
# This target should be expanded upon as we add more Linux architectures: i.e. build-arm64.
# Then a new build target can be created such as build-container-bin that builds these
# various platform images.
.PHONY: build-linux
build-linux: depend
	GOOS=linux GOARCH=$(GOARCH_LOCAL) LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $(TARGET_OUT_LINUX)/ $(STANDARD_BINARIES)
	GOOS=linux GOARCH=$(GOARCH_LOCAL) LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $(TARGET_OUT_LINUX)/ -tags=agent $(AGENT_BINARIES)

# Create targets for TARGET_OUT_LINUX/binary
# There are two use cases here:
# * Building all docker images (generally in CI). In this case we want to build everything at once, so they share work
# * Building a single docker image (generally during dev). In this case we just want to build the single binary alone
BUILD_ALL ?= true
define build-linux
.PHONY: $(TARGET_OUT_LINUX)/$(shell basename $(1))
ifeq ($(BUILD_ALL),true)
$(TARGET_OUT_LINUX)/$(shell basename $(1)): build-linux
	@:
else
$(TARGET_OUT_LINUX)/$(shell basename $(1)): $(TARGET_OUT_LINUX)
	GOOS=linux GOARCH=$(GOARCH_LOCAL) LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $(TARGET_OUT_LINUX)/ -tags=$(2) $(1)
endif
endef

$(foreach bin,$(STANDARD_BINARIES),$(eval $(call build-linux,$(bin),"")))
$(foreach bin,$(AGENT_BINARIES),$(eval $(call build-linux,$(bin),"agent")))

# Create helper targets for each binary, like "pilot-discovery"
# As an optimization, these still build everything
$(foreach bin,$(BINARIES),$(shell basename $(bin))): build
ifneq ($(TARGET_OUT_LINUX),$(LOCAL_OUT))
# if we are on linux already, then this rule is handled by build-linux above, which handles BUILD_ALL variable
$(foreach bin,$(BINARIES),${LOCAL_OUT}/$(shell basename $(bin))): build
endif

MARKDOWN_LINT_ALLOWLIST=localhost:8080,storage.googleapis.com/istio-artifacts/pilot/,http://ratings.default.svc.cluster.local:9080/ratings

lint-helm-global:
	find manifests -name 'Chart.yaml' -print0 | ${XARGS} -L 1 dirname | xargs -r helm lint

lint: lint-python lint-copyright-banner lint-scripts lint-go lint-dockerfiles lint-markdown lint-yaml lint-licenses lint-helm-global ## Runs all linters.
	@bin/check_samples.sh
	@testlinter
	@envvarlinter istioctl pilot security

go-gen:
	@mkdir -p /tmp/bin
	@PATH="${PATH}":/tmp/bin go generate ./...

refresh-goldens:
	@REFRESH_GOLDEN=true go test ${GOBUILDFLAGS} ./operator/...
	@REFRESH_GOLDEN=true go test ${GOBUILDFLAGS} ./pkg/kube/inject/...
	@REFRESH_GOLDEN=true go test ${GOBUILDFLAGS} ./pilot/pkg/security/authz/builder/...
	# we needn't cni here
	@#REFRESH_GOLDEN=true go test ${GOBUILDFLAGS} ./cni/pkg/plugin/...

update-golden: refresh-goldens

# Keep dummy target since some build pipelines depend on this
gen-charts:
	@echo "This target is no longer required and will be removed in the future"

gen: \
	mod-download-go \
	go-gen \
	mirror-licenses \
	format \
	update-crds \
	proto \
	copy-templates \
	gen-kustomize \
	update-golden ## Update all generated code.

gen-check: gen check-clean-repo

copy-templates:
	rm manifests/charts/istiod-remote/templates/*
	rm manifests/charts/gateways/istio-egress/templates/*

	# gateway charts
	cp -r manifests/charts/gateways/istio-ingress/templates/* manifests/charts/gateways/istio-egress/templates
	find ./manifests/charts/gateways/istio-egress/templates -type f -exec sed -i -e 's/ingress/egress/g' {} \;
	find ./manifests/charts/gateways/istio-egress/templates -type f -exec sed -i -e 's/Ingress/Egress/g' {} \;

	# external istiod remote cluster charts
	cp manifests/charts/base/templates/services.yaml manifests/charts/istiod-remote/templates
	cp manifests/charts/base/templates/endpoints.yaml manifests/charts/istiod-remote/templates
	cp manifests/charts/base/templates/reader-serviceaccount.yaml manifests/charts/istiod-remote/templates
	cp manifests/charts/istio-control/istio-discovery/templates/mutatingwebhook.yaml manifests/charts/istiod-remote/templates
	cp manifests/charts/istio-control/istio-discovery/templates/reader-clusterrole.yaml manifests/charts/istiod-remote/templates
	cp manifests/charts/istio-control/istio-discovery/templates/reader-clusterrolebinding.yaml manifests/charts/istiod-remote/templates

	# external istiod config cluster charts
	cp manifests/charts/istio-control/istio-discovery/files/injection-template.yaml manifests/charts/istiod-remote/files
	cp manifests/charts/istio-control/istio-discovery/files/gateway-injection-template.yaml manifests/charts/istiod-remote/files
	cp manifests/charts/istio-control/istio-discovery/templates/istiod-injector-configmap.yaml manifests/charts/istiod-remote/templates
	cp manifests/charts/istio-control/istio-discovery/templates/configmap.yaml manifests/charts/istiod-remote/templates
	cp manifests/charts/istio-control/istio-discovery/templates/telemetryv2_*.yaml manifests/charts/istiod-remote/templates
	sed -e '1 i {{- if .Values.global.configCluster }}' -e '$$ a {{- end }}' manifests/charts/base/crds/crd-all.gen.yaml > manifests/charts/istiod-remote/templates/crd-all.gen.yaml
	sed -e '1 i {{- if .Values.global.configCluster }}' -e '$$ a {{- end }}' manifests/charts/base/crds/crd-operator.yaml > manifests/charts/istiod-remote/templates/crd-operator.yaml
	sed -e '1 i {{- if .Values.global.configCluster }}' -e '$$ a {{- end }}' manifests/charts/istio-control/istio-discovery/templates/validatingwebhookconfiguration.yaml > manifests/charts/istiod-remote/templates/validatingwebhookconfiguration.yaml
	sed -e '1 i {{- if .Values.global.configCluster }}' -e '$$ a {{- end }}' manifests/charts/istio-control/istio-discovery/templates/serviceaccount.yaml > manifests/charts/istiod-remote/templates/serviceaccount.yaml
	sed -e '1 i {{- if .Values.global.configCluster }}' -e '$$ a {{- end }}' manifests/charts/istio-control/istio-discovery/templates/role.yaml > manifests/charts/istiod-remote/templates/role.yaml
	sed -e '1 i {{- if .Values.global.configCluster }}' -e '$$ a {{- end }}' manifests/charts/istio-control/istio-discovery/templates/rolebinding.yaml > manifests/charts/istiod-remote/templates/rolebinding.yaml
	sed -e '1 i {{- if .Values.global.configCluster }}' -e '$$ a {{- end }}' manifests/charts/istio-control/istio-discovery/templates/clusterrole.yaml > manifests/charts/istiod-remote/templates/clusterrole.yaml
	sed -e '1 i {{- if .Values.global.configCluster }}' -e '$$ a {{- end }}' manifests/charts/istio-control/istio-discovery/templates/clusterrolebinding.yaml > manifests/charts/istiod-remote/templates/clusterrolebinding.yaml

	# copy istio-discovery values, but apply some local customizations
	cp manifests/charts/istio-control/istio-discovery/values.yaml manifests/charts/istiod-remote/
	yq -i '.telemetry.enabled=false | .global.externalIstiod=true | .global.omitSidecarInjectorConfigMap=true | .pilot.configMap=false' manifests/charts/istiod-remote/values.yaml
# Generate kustomize templates.
gen-kustomize:
	helm3 template istio --namespace istio-system --include-crds manifests/charts/base > manifests/charts/base/files/gen-istio-cluster.yaml
	helm3 template istio --namespace istio-system manifests/charts/istio-control/istio-discovery \
		> manifests/charts/istio-control/istio-discovery/files/gen-istio.yaml
	helm3 template operator --namespace istio-operator manifests/charts/istio-operator \
		--set hub=gcr.io/istio-testing --set tag=${VERSION} > manifests/charts/istio-operator/files/gen-operator.yaml

#-----------------------------------------------------------------------------
# Target: go build
#-----------------------------------------------------------------------------

# Non-static istioctl targets. These are typically a build artifact.
${TARGET_OUT}/release/istioctl-linux-amd64: depend
	GOOS=linux GOARCH=amd64 LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $@ ./istioctl/cmd/istioctl
${TARGET_OUT}/release/istioctl-linux-armv7: depend
	GOOS=linux GOARCH=arm GOARM=7 LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $@ ./istioctl/cmd/istioctl
${TARGET_OUT}/release/istioctl-linux-arm64: depend
	GOOS=linux GOARCH=arm64 LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $@ ./istioctl/cmd/istioctl
${TARGET_OUT}/release/istioctl-osx: depend
	GOOS=darwin GOARCH=amd64 LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $@ ./istioctl/cmd/istioctl
${TARGET_OUT}/release/istioctl-osx-arm64: depend
	GOOS=darwin GOARCH=arm64 LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $@ ./istioctl/cmd/istioctl
${TARGET_OUT}/release/istioctl-win.exe: depend
	GOOS=windows LDFLAGS=$(RELEASE_LDFLAGS) common/scripts/gobuild.sh $@ ./istioctl/cmd/istioctl

# generate the istioctl completion files
${TARGET_OUT}/release/istioctl.bash: ${LOCAL_OUT}/istioctl
	${LOCAL_OUT}/istioctl completion bash > ${TARGET_OUT}/release/istioctl.bash

${TARGET_OUT}/release/_istioctl: ${LOCAL_OUT}/istioctl
	${LOCAL_OUT}/istioctl completion zsh > ${TARGET_OUT}/release/_istioctl

.PHONY: binaries-test
binaries-test:
	go test ${GOBUILDFLAGS} ./tests/binary/... -v --base-dir ${TARGET_OUT} --binaries="$(RELEASE_BINARIES)"

# istioctl-all makes all of the non-static istioctl executables for each supported OS
.PHONY: istioctl-all
istioctl-all: ${TARGET_OUT}/release/istioctl-linux-amd64 ${TARGET_OUT}/release/istioctl-linux-armv7 ${TARGET_OUT}/release/istioctl-linux-arm64 \
	${TARGET_OUT}/release/istioctl-osx \
	${TARGET_OUT}/release/istioctl-osx-arm64 \
	${TARGET_OUT}/release/istioctl-win.exe

.PHONY: istioctl.completion
istioctl.completion: ${TARGET_OUT}/release/istioctl.bash ${TARGET_OUT}/release/_istioctl

# istioctl-install builds then installs istioctl into $GOPATH/BIN
# Used for debugging istioctl during dev work
.PHONY: istioctl-install-container
istioctl-install-container: istioctl

#-----------------------------------------------------------------------------
# Target: test
#-----------------------------------------------------------------------------

.PHONY: test

# This target sets JUNIT_REPORT to the location of the  go-junit-report binary.
# This binary is provided in the build container. If it is not found, the build
# container is not being used, so ask the user to install go-junit-report.
JUNIT_REPORT := $(shell which go-junit-report 2> /dev/null || echo "${ISTIO_BIN}/go-junit-report")

${ISTIO_BIN}/go-junit-report:
	@echo "go-junit-report was not found in the build environment."
	@echo "Please install go-junit-report (ex. go install github.com/jstemmer/go-junit-report@latest)"
	@exit 1

# This is just an alias for racetest now
test: racetest ## Runs all unit tests

# For now, keep a minimal subset. This can be expanded in the future.
BENCH_TARGETS ?= ./pilot/...

.PHONY: racetest
racetest: $(JUNIT_REPORT)
	go test ${GOBUILDFLAGS} ${T} -race ./... 2>&1 | tee >($(JUNIT_REPORT) > $(JUNIT_OUT))

.PHONY: benchtest
benchtest: $(JUNIT_REPORT) ## Runs all benchmarks
	prow/benchtest.sh run $(BENCH_TARGETS)
	prow/benchtest.sh compare

report-benchtest:
	prow/benchtest.sh report

#-----------------------------------------------------------------------------
# Target: clean
#-----------------------------------------------------------------------------
.PHONY: clean

clean: ## Cleans all the intermediate files and folders previously generated.
	rm -rf $(DIRS_TO_CLEAN)

#-----------------------------------------------------------------------------
# Target: docker
#-----------------------------------------------------------------------------
.PHONY: push

# for now docker is limited to Linux compiles - why ?
include tools/istio-docker.mk

push: docker.push ## Build and push docker images to registry defined by $HUB and $TAG

#-----------------------------------------------------------------------------
# Target: environment and tools
#-----------------------------------------------------------------------------
.PHONY: show.env show.goenv

show.env: ; $(info $(H) environment variables...)
	$(Q) printenv

show.goenv: ; $(info $(H) go environment...)
	$(Q) $(GO) version
	$(Q) $(GO) env

# show makefile variables. Usage: make show.<variable-name>
show.%: ; $(info $* $(H) $($*))
	$(Q) true

#-----------------------------------------------------------------------------
# Target: custom resource definitions
#-----------------------------------------------------------------------------

update-crds:
	bin/update_crds.sh

#-----------------------------------------------------------------------------
# Target: artifacts and distribution
#-----------------------------------------------------------------------------
# deb, rpm, etc packages
include tools/packaging/packaging.mk

#-----------------------------------------------------------------------------
# Target: integration tests
#-----------------------------------------------------------------------------
include tests/integration/tests.mk

include common/Makefile.common.mk
