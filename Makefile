


PROTOC_2_BIN ?= protoc2
PROTOC_INCLUDE := -I$(GOPATH)/src -I. -Ivendor -Ivendor/github.com/gogo/protobuf/protobuf -Ischeduler/ -Ioperator/agent
PROTOC_INCLUDE := $(PROTOC_INCLUDE) -Ioperator/allocator -Ioperator/maintenance -Ioperator/quota
PROTOC_INCLUDE := $(PROTOC_INCLUDE) -Ioperator/master

PACKAGES ?= $(shell go list ./... | grep -v vendor)
PROJECT_NAMESPACE := github.com/ondrej-smola/mesos-go-http

.PHONY: install
install: install-dependencies binaries

.PHONY: test
test:
	ginkgo -r -skip vendor

.PHONY: vet
vet:
	@go $@ $(PACKAGES)

.PHONY: fmt
fmt:
	@go $@ $(PACKAGES)

.PHONY: install-dependencies
install-dependencies:
	go get github.com/gogo/protobuf/protoc-gen-gogoslick

.PHONY: install-test-dependencies
install-test-dependencies:
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega

.PHONY: proto
proto:
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) mesos.proto --gogoslick_out=.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) scheduler/scheduler.proto --gogoslick_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) operator/agent/agent.proto --gogoslick_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) operator/master/master.proto --gogoslick_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) operator/allocator/allocator.proto --gogoslick_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) operator/maintenance/maintenance.proto --gogoslick_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) operator/quota/quota.proto --gogoslick_out=$(MESOS_PROTO_MAPPING):.

binaries:
	cd examples && go install ./...

clean-proto:
	$(RM) mesos.pb.go
	cd scheduler && $(RM) scheduler.pb.go