
ifdef SystemRoot
	PATHSEP := \\
	RM := del
else
	PATHSEP := /
	RM := rm
endif

PROTOC_2_BIN ?= protoc2
PROTOC_INCLUDE := -I$(GOPATH)/src/github.com/gogo/protobuf/protobuf -I$(GOPATH)/src -I. -Ilib/scheduler/
PROTOC_INCLUDE := $(PROTOC_INCLUDE) -Ilib/operator/allocator -Ilib/operator/maintenance -Ilib/operator/quota
PROTOC_INCLUDE := $(PROTOC_INCLUDE) -Ilib/operator/agent -Ilib/operator/resource_provider -Ilib/operator/master

PROJECT_NAMESPACE := github.com/ondrej-smola/mesos-go-http

.PHONY: install
install: install-dependencies binaries

.PHONY: pre-commit
pre-commit: fmt vet test

.PHONY: test
test:
	ginkgo -r -skip vendor

.PHONY: vet
vet:
	@echo go vet
	@go tool vet lib
	@go tool vet examples$(PATHSEP)operator
	@go tool vet examples$(PATHSEP)scheduler

.PHONY: fmt
fmt:
	@echo go fmt
	@goimports -w lib
	@go tool vet examples$(PATHSEP)operator
	@go tool vet examples$(PATHSEP)scheduler

.PHONY: install-dependencies
install-dependencies:
	go get golang.org/x/tools/cmd/goimports
	go get github.com/pkg/errors
	go get github.com/gogo/protobuf/protoc-gen-gofast

.PHONY: install-test-dependencies
install-test-dependencies:
	go get github.com/golang/protobuf/proto
	go get github.com/onsi/ginkgo/ginkgo
	go get github.com/onsi/gomega

.PHONY: proto
proto:
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/mesos.proto --gofast_out=.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/scheduler/scheduler.proto --gofast_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/executor/executor.proto --gofast_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/operator/agent/agent.proto --gofast_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/operator/master/master.proto --gofast_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/operator/allocator/allocator.proto --gofast_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/operator/maintenance/maintenance.proto --gofast_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/operator/quota/quota.proto --gofast_out=$(MESOS_PROTO_MAPPING):.
	@$(PROTOC_2_BIN) $(PROTOC_INCLUDE) lib/operator/resource_provider/resource_provider.proto --gofast_out=$(MESOS_PROTO_MAPPING):.

.PHONY: binaries
binaries:
	cd examples$(PATHSEP)scheduler && go install
	cd examples$(PATHSEP)operator && go install

.PHONY: proto-clean
proto-clean:
	$(RM) lib$(PATHSEP)mesos.pb.go
	$(RM) lib$(PATHSEP)scheduler$(PATHSEP)scheduler.pb.go
	$(RM) lib$(PATHSEP)executor$(PATHSEP)executor.pb.go
	$(RM) lib$(PATHSEP)operator$(PATHSEP)agent$(PATHSEP)agent.pb.go
	$(RM) lib$(PATHSEP)operator$(PATHSEP)allocator$(PATHSEP)allocator.pb.go
	$(RM) lib$(PATHSEP)operator$(PATHSEP)maintenance$(PATHSEP)maintenance.pb.go
	$(RM) lib$(PATHSEP)operator$(PATHSEP)master$(PATHSEP)master.pb.go
	$(RM) lib$(PATHSEP)operator$(PATHSEP)quota$(PATHSEP)quota.pb.go
	$(RM) lib$(PATHSEP)operator$(PATHSEP)resource_provider$(PATHSEP)resource_provider.pb.go
