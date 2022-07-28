# Contributing

We use GitHub to host code, build code, create releases, track issues and feature requests, and accept pull requests.

## MIT Software License

In short, when you submit code changes, your submissions are understood to be under the same [MIT License](https://github.com/calvinbui/teleport-sre-challenge/blob/master/LICENSE) that covers the project. Feel free to contact the maintainers if that's a concern.

## Report Bugs or Request Features using GitHub Issues

We use GitHub issues to raise any bug or request new features. GitHub will present templates for both.

## Contribution Process

Pull requests are the best way to propose changes to the codebase (we use [GitHub Flow](https://guides.github.com/introduction/flow/index.html)). We actively welcome your pull requests:

1. Fork this repo and create a branch from master.
2. If you've added code that should be tested, add tests.
3. If you've changed the config, update the README.md file.
4. Ensure the test suite passes.
5. Make sure your code lints.
6. Issue that pull request!

## Local Development

The following are required:

- go
- make
- Docker
- minikube
- Helm
- protobuf and protoc-gen-go-grpc

Certificates can be generated with `./certificates/gen.sh`.

When developing the API locally, developers will have access to convenient makefile commands to build, run and test the code.

- `make build`: Build the API on your workstation
- `make run`: Run the API on your workstation
- `make run-k8s`: Spins up a minikube cluster and runs the API on it
- `make test`: Run unit tests on your workstation
- `make test-k8s`: Spins up a minikube and run unit tests on it
- `make helm-deploy-dev`: Deploy the Helm chart to minikube

## Code Review

The API will use GitHub to manage changes to the code base.

1. Developers will create a GitHub Issue (or find an existing Issue) and claim it by commenting on the issue that they will work on it. This is to help prevent duplicated efforts from developers.

2. A new git branch in the format `<username>/<change>` (i.e. `calvinbui/add-mtls-to-api`) will be created.

3. A GitHub Pull Request will be created from the developer's branch to the primary branch (e.g. `main` or `master`). Contributors will be encouraged to follow the template provided.

4. A GitHub Actions CI/CD pipeline will be triggered on the pull request to check for code quality and perform tests. The tests must pass successfully.

5. The code maintainers will review the changes and comment, approve request for changes.

6. At least two code maintainers have to approve the changes before they can be merged into the primary branch. This merge is performed by a code maintainer.

## Release Process

The API will be tested, built and deployed using a GitHub Actions CI/CD pipeline. All (and only) changes to the codebase's primary branch will be deployed.

As the API runs in Kubernetes, the GitHub Actions pipeline will produce container image artifacts. These will be stored inside Amazon's Elastic Container Registry (ECR). [With IAM Roles, EKS nodes can be granted permissions to ECR](https://docs.aws.amazon.com/AmazonECR/latest/userguide/ECR_on_EKS.html), keeping security and traffic within AWS.

Deployment to Kubernetes will be done with Helm Charts. The chart will be kept inside the code base and have its versioning separate from the API. The chart will include, but is not limited to, the following resources:

## Dependency Management

This project uses [Go modules](https://golang.org/cmd/go/#hdr-Modules__module_versions__and_more) to manage dependencies on external packages

To add or update a new dependency, use the `go get` command:

```bash
# Pick the latest tagged release.
go get example.com/some/module/pkg

# Pick a specific version.
go get example.com/some/module/pkg@vX.Y.Z
```
You have to commit the changes to `go.mod` and `go.sum` before submitting the pull request.
