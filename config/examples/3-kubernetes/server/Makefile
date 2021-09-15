docker_image ?= gardenbed/app-server

build:
	go build

docker:
	docker image build --tag $(docker_image):latest .
	docker image push $(docker_image):latest

k8s-deploy:
	kubectl create -f kubernetes.yaml

k8s-delete:
	kubectl delete service app-server
	kubectl delete deployment app-server
	kubectl delete configmap app-server

.PHONY: build docker k8s-deploy k8s-delete
