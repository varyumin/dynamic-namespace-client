VERSION ?= $(shell git rev-parse --short HEAD || echo "GitNotFound")

.PHONY: build
build:
	dep ensure --vendor-only
	go build -ldflags "-X main.VERSION=$(VERSION)" .

.PHONY: clean
clean:
	rm -fr vendor
	rm -f kubectl-dynamicns

.PHONY: docker-build
docker-build:
	docker build -t registry.tcsbank.ru:5050/k8s/dynamic-namespace-client:$(VERSION) .

.PHONY: docker-push
docker-push:
	docker push registry.tcsbank.ru:5050/k8s/dynamic-namespace-client:$(VERSION)