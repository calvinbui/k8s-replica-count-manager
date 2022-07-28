BINARY_NAME=replica-count-manager
DEV_NAMESPACE=default

run:
	go run ./cmd/${BINARY_NAME}/

# this is a bit slow as it builds inside of minikube
# alternative solutions can be:
# 1. building the go binary locally and copying it into the container
# 2. deploying using ko by Google. a bit faster but environment wouldn't match deployed instances
run-k8s: minikube-start docker-build-dev helm-deploy-dev restart-k8s

test:
	go test ./...

test-k8s: minikube-start docker-build-dev helm-deploy-test

build:
	mkdir -p out
	go build -o out/${BINARY_NAME} cmd/${BINARY_NAME}/

minikube-start:
	if ! minikube status >/dev/null; then minikube start; fi

docker-build-dev:
	@eval $$(minikube docker-env) ;\
	docker build -t ${BINARY_NAME}:dev .

helm-deploy-dev:
	kubectl config set-context minikube
	cd charts/${BINARY_NAME} && helm upgrade --install -n ${DEV_NAMESPACE} -f values-dev.yaml replica-count-manager .

# this command is used in GitHub actions, but it is also possible to use it locally
helm-deploy-test:
	ct install --charts charts/${BINARY_NAME}

# part of the challenge asked to upgrade a local cluster with no service interruption
# rolling restart of a deployment will allow this
restart-k8s:
	kubectl config set-context minikube
	kubectl rollout restart deployment replica-count-manager -n ${DEV_NAMESPACE}
