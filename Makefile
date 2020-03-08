.PHONY: default
default: build

VER=$(shell git describe --always --tags --dirty)
IMAGE:=octet/operator-nodejs
CRD:=deploy/crds/opnodejs.octetcloud.com_nodejsdiagnosticreports_crd.yaml


.PHONY: build
build:
	operator-sdk build $(IMAGE)

.PHONY: push
push: build
	docker push $(IMAGE)

kube-init:
	kubectl create -f $(CRD)

kube-run:
	operator-sdk run --local --namespace=default
