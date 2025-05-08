# KubeNap Makefile

.PHONY: deps build run tidy

deps:
	@echo "> Installing Go dependencies"
	go get k8s.io/client-go@v0.29.0
	go get k8s.io/apimachinery@v0.29.0
	go get k8s.io/api@v0.29.0
	go get k8s.io/client-go/informers@v0.29.0

build:
	@echo "> Building KubeNap"
	go build -o bin/kubenap ./cmd/kubenap

run:
	@echo "> Running KubeNap"
	go run ./cmd/kubenap

tidy:
	@echo "> Cleaning up modules"
	go mod tidy
