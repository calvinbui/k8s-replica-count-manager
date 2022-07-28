---
authors: Calvin Bui (3604363+calvinbui@users.noreply.github.com)
state: draft
---

# RFD 180722 - Kubernetes Replica Count Manager (Teleport SRE Challenge)

## What

Provide an API alongside the native Kubernetes deployment controller to manage the pod replica count of deployment resources. The API stores the desired replica count, and attempts to reconcile this value if it does not match the current replica count.

## Why

- Provides a single source of truth for replica counts of all resources
- Deployments and upgrades may reset the replica counts back to a previous number
- Replica counts may be modified by external actors

## Details

The API is intended to be:
  - Kubernetes native
  - run on Amazon Elastic Kubernetes Service (EKS)
  - templated through Helm Charts
  - easy to develop locally
  - communicated with via gRPC and mTLS
  - built, tested and deployed with GitHub Actions

### API

The API will be served via a gRPC endpoint. Clients will have to connect through mutual TLS (mTLS). The following functions will be available to clients:

`ListDeployment`: List of available deployments in the Kubernetes cluster.

```go
message ListDeployment {
  string Name = 1;
  string Namespace = 2;
}
```

`GetDeploymentReplicas`: Get the replica count (desired and current) of the Kubernetes deployment.

```go
message GetDeployment {
  string Name = 1;
  string Namespace = 2;
}
```

`SetDeploymentReplicas`: Set the desired replica count of the Kubernetes deployment resource.

```go
message SetDeploymentReplicas {
  string Name = 1;
  string Namespace = 2;
  int32 Replicas = 3;
}
```

In addition to this, the API will also have a health check endpoint at `/healthz` to verify Kubernetes connectivity.

### State management

The replica count values for each deployment will be persisted when set through the API in a local persistent volume. When no value is present, the deployment's replica count is not managed by this API. Information stored include the name of the deployment, the namespace and replica count. Updates to any values will be stored back into the file.

The file is read once on start up and the values will be stored in memory. The data will be stored in JSON. An example of the data:

```json
[
  {
    "name": "echo-server",
    "namespace": "echo",
    "replicas": 1
  },
  {
    "..."
  }
]
```

### Reconciliation

The API will use the Kubernetes [watch](https://pkg.go.dev/k8s.io/apimachinery/pkg/watch) interface to receive events for replica count values.

Reconciliation actions are performed when the desired replica count inside the API's state does not match the current replica count on the cluster. The API will set the replica count on the deployment to its desired value.

### Developer workflow

#### Local development

When developing the API locally, developers will have access to convenient `makefile` commands to build, run and test the code.

- `make build`: Build the API on your workstation
- `make run`: Run the API on your workstation
- `make run-k8s`: Spins up a minikube cluster and runs the API on it
- `make test`: Run unit tests on your workstation
- `make test-k8s`: Spins up a minikube and run unit tests on it

#### Code Review
The API will use GitHub to manage changes to the code base.

1. Developers will create a GitHub Issue (or find an existing Issue) and claim it by commenting on the issue that they will work on it. This is to help prevent duplicated efforts from developers.

2. A new git branch in the format `<username>/<change>` (i.e. `calvinbui/add-mtls-to-api`) will be created.

3. A GitHub Pull Request will be created from the developer's branch to the primary branch (e.g. `main` or `master`). Contributors will be encouraged to follow the template provided.

4. A GitHub Actions CI/CD pipeline will be triggered on the pull request to check for code quality and perform tests. The tests must pass successfully.

5. The code maintainers will review the changes and comment, approve request for changes.

6. At least two code maintainers have to approve the changes before they can be merged into the primary branch. This merge is performed by a code maintainer.

### Build and Release

The API will be tested, built and deployed using a GitHub Actions CI/CD pipeline. All (and only) changes to the codebase's primary branch will be deployed.

As the API runs in Kubernetes, the GitHub Actions pipeline will produce container image artifacts. These will be stored inside Amazon's Elastic Container Registry (ECR). [With IAM Roles, EKS nodes can be granted permissions to ECR](https://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_on_EKS.html), keeping security and traffic within AWS.

Deployment to Kubernetes will be done with Helm Charts. The chart will be kept inside the code base and have its versioning separate from the API. The chart will include, but is not limited to, the following resources:

- Namespace
- Deployment
- Service Account
- Service
- Role
- Role Binding
- Network Policy
- Pod Disruption Budget
- Config Map

GitHub Actions uses Runners as build agents that perform actions inside a workflow. Runners can be self-hosted. It is preferable to deploy the GitHub Actions Runner inside of the EKS cluster to make use of [IAM Roles for Service Accounts (IRSA)][IRSA] to provide it with permissions to AWS ECR and the EKS cluster to push new container images and deploy the Helm chart templates.

## Possible Future Improvements

- Use a scalable tool for state management
  - Better data structure to avoid looping through all items
  - Push instead of pull when updating state
  - Locking
- High Availability
  - Support running multiple replicas
  - Allow only one replica to be update deployments
    - Primary/secondary system
    - Use Kubernetes service to discover other replicas and achieve quorum
 Notifications (e.g. Slack messages) when reconciliation actions are executed
- Extra features
  - Delete replica count settings
  - Ignore replica count settings
  - Retry logic

## Appendix

### Assumptions

- The Amazon EKS cluster has been deployed
- The GitHub Actions Runner has been deployed
- A storage class is available to the Kubernetes cluster


[IRSA]: https://docs.aws.amazon.com/emr/latest/EMR-on-EKS-DevelopmentGuide/setting-up-enable-IAM.html
