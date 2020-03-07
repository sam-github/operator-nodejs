.PHONY: default
default: build

VER=$(shell git describe --always --tags --dirty)
IMAGE:=octet/operator-nodejs


.PHONY: build
build:
	operator-sdk build $(IMAGE)

.PHONY: push
push: build
	docker push $(IMAGE)
