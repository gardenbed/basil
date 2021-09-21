docker_image ?= gardenbed/app-client

build:
	go build

docker:
	docker image build --tag $(docker_image):latest .
	docker image push $(docker_image):latest

k8s-deploy:
	kubectl create -f kubernetes.yaml

k8s-delete:
	kubectl delete deployment app-client
	kubectl delete configmap app-client

.PHONY: build docker k8s-deploy k8s-delete
