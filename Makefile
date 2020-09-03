VERSION?="0.0.1"
TEST?=./...
GOFMT_FILES?=$$(find . -type f -name '*.go')

default: build deploy

build: fmt
	docker build -f ./deploy/Dockerfile -t elemir/templatemap:latest .
	docker push elemir/templatemap:latest

deploy:
	kubectl delete -f ./deploy/manifest.yaml
	kubectl apply -f ./deploy/manifest.yaml

fmt:
	gofmt -d $(GOFMT_FILES)
	gofmt -w $(GOFMT_FILES)

test: fmt
	go test ./...

.PHONY: default build deploy fmt test
