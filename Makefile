.PHONY: default
default: build

VER=$(shell git describe --always --tags --dirty)
IMAGE:=octet/operator-nodejs
CRD:=deploy/crds/opnodejs.appsody.dev_nodejsreports_crd.yaml


.PHONY: build
build:
	operator-sdk build $(IMAGE)

.PHONY: check-build
check-build:
	go build -o build/_output/bin/operator-nodejs-local ./cmd/manager

.PHONY: push
push: build
	docker push $(IMAGE)

kube-init:
	kubectl create -f $(CRD) || kubectl replace -f $(CRD)

local-run: check-build
	WATCH_NAMESPACE=default ./build/_output/bin/operator-nodejs-local

example-report:
	./bin/create-report
