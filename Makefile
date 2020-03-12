.PHONY: default
default: build

VER=$(shell git describe --always --tags --dirty)
IMAGE:=octet/operator-nodejs
CRD:=deploy/crds/opnodejs.octetcloud.com_nodejsdiagnosticreports_crd.yaml


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
	kubectl create -f $(CRD)

local-run:
	operator-sdk run --local --namespace=default

example-report:
	kubectl apply -f deploy/crds/opnodejs.octetcloud.com_v1alpha1_nodejsdiagnosticreport_cr.yaml
