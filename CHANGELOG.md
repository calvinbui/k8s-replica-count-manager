# v0.1.0

## API
- [x] gRPC API
  - [x] ListDeployment
  - [x] GetDeploymentReplicas
    - [x] Current replica count
    - [x] Desired replica count (return null if not set yet?)
  - [x] SetDeploymentReplicas
    - [x] Show the current and new desired state
    - [x] Store replica count
- [x] Reconciling cluster state (i.e. an external actor changed the replica count manually)
- [x] HTTP health check
- [x] mTLS

## Tests
- [x] One or two tests per feature
- [x] Integration tests with Kubernetes cluster
- [x] Deploy and manage a local Kubernetes cluster
- [x] Deploy and upgrade to the local Kubernetes cluster with no service interruption
- [x] README showing everything

## Deployment
- [x] Dockerfile
- [x] makefile to build and run tests
- [x] GitHub Actions
  - [x] Code quality
  - [x] Tests
    - [x] Static analysis
    - [x] Unit tests
  - [x] Build
  - [x] Deployment / Upgrade
- [x] Helm Chart
  - [x] Deployment
  - [x] ServiceAccount
  - [x] Service
  - [x] ClusterRole
  - [x] ClusterRoleBinding
  - [x] PersistentVolumeClaim
  - [x] PodDisruptionBudget
  - [x] README

