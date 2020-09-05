VERSION?="0.0.1"
TEST?=./...
GOFMT_FILES?=$$(find . -type f -name '*.go')

default: build deploy

build: fmt
	docker build -f ./deploy/Dockerfile -t elemir/templatemap:latest .
	docker push elemir/templatemap:latest

deploy:
	kubectl delete --ignore-not-found=true -f ./deploy/manifest.yaml
	kubectl apply -f ./deploy/manifest.yaml

fmt:
	gofmt -d $(GOFMT_FILES)
	gofmt -w $(GOFMT_FILES)

test: fmt
	go test ./...

debug: build deploy
	kubectl delete --ignore-not-found=true -f examples/pod.yaml
	kubectl apply -f examples/pod.yaml
	NODE=$$(kubectl get pod nginx -o=json | jq -r '.spec.nodeName'); \
	POD=$$(kubectl get pod -n kube-system -l csi-plugin=templatemap --field-selector spec.nodeName=$${NODE} -o json | jq -r '.items[0].metadata.name'); \
	kubectl logs -f -n kube-system $${POD} --all-containers

.PHONY: default build deploy fmt test
