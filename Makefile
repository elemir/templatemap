VERSION?="0.0.1"
TEST?=./...
GOFMT_FILES?=$$(find . -type f -name '*.go')

default: build deploy

build: fmt
	docker build -f ./deploy/Dockerfile -t elemir/templatemap:latest .

deploy:
	docker push elemir/templatemap:latest
	kubectl delete -f ./deploy/manifest.yaml
	kubectl apply -f ./deploy/manifest.yaml

fmt:
	gofmt -d $(GOFMT_FILES)
	gofmt -w $(GOFMT_FILES)

test:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

.PHONY: default build deploy fmt test
