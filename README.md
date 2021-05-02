# GO grpc microservices demonstrator

Demonstrator of creating a full microservices development and deployment solution built on grpc interservice communication, graphql api gateway and skaffold.

This includes the pattern of service stubbing using a custom stubber in `./proto/stub`.

## Local development

Requires:

- [Docker](https://docker.com)
- [Minikube](https://minikube.sigs.k8s.io/docs/start/)
- [Skaffold](https://skaffold.dev/)
- [Kubectl](https://kubectl.docs.kubernetes.io/)
- [Kustomize](https://kubectl.docs.kubernetes.io/)
- [Golang > 1.16](https://golang.org/)
