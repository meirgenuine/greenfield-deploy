# greenfield-deploy
This project supports greenfield's projects deployment in Kubernetes using Docker images + Github Actions as workflow.

## Table of Contents
- [Usage](#deploy-projects-with-greenfield-deploy)
- [How it works](#how-it-works)
- [How to use GitHub Actions workflows](#github-actions-workflows)
- [Container Orchestration](#container-orchestration)
- [How to write k8s manifests](#kubernetes-manifests)
- [Telegram Bot](#telegram-bot)
- [Local setup](#local-setup)
- [License](#license)


### [Deploy projects with greenfield-deploy](#deploy-projects-with-greenfield-deploy)
To deploy specific project in kubernetes cluster do the following steps:
1. Be sure that all tests are passed and docker image is pushed into container registry. (See [How to use GitHub Actions workflows](#github-actions-workflows))
2. Send message to [Telegram Bot](#telegram-bot) with required parameters.
3. Receive message from bot about deployment status.
4. Double check that deployment succeeded by:
```
$ kubectl get pods -n <namespace>
```

### [How it works](#how-it-works)
On each PR to the main branch of the specific project Github Actions starts running all checks (tests, linters, etc...). After successfull checks Github Actions builds and pushes new docker image to docker container registry. Afterwards, developer who is reponsible for that commit can deploy this server via the Telegram Bot by sending message to it. Greenfield-deploy service receives the message from Telegram Bot and downloads k8s manifests from the repo. As the last step, Greenfield-deploy service applies this k8s manifests with specific version to Kubernetes Cluster.

As an example greenfield-deploy project contains workflow for greenfield project.

### [GitHub Actions Workflows](#github-actions-workflows)

[This repository](https://github.com/meirgenuine/greenfield) utilizes several GitHub Actions workflows to automate common tasks. Below is a brief description of the workflows and how they operate.

### 1. Unit Tests, Gosec and Lint Workflow

This workflow comprises three main jobs: **build-test**, **gosec**, and **golangci-lint**.

The workflow is triggered when there's a new commit pushed to the **master** or **develop** branches, or when a new pull request is opened.

Environment Variables and Secrets
- GH_ACCESS_TOKEN: GitHub access token used for cloning and downloading Go modules. Stored as a GitHub Secret.

#### Build and Unit Test Job

This job is responsible for building the code and running unit tests. Here's what happens during this job:

1. **Setup Go:** The appropriate Go version is installed.
2. **Checkout Code:** The latest code is checked out from the repository.
3. **Setup Caching:** A cache is set up for Go modules and build cache.
4. **Setup GitHub Token:** The GitHub token is configured to allow fetching of private Go modules.
5. **Build:** The code is built using the make build command.
6. **Unit Tests:** Unit tests are run using the make test command.

#### Gosec Job

This job runs the Gosec security scanner on the codebase. Here's what happens during this job:

7. **Setup Go and Checkout Code:** Similar to the build and test job, Go is setup, code is checked out, and GitHub token is configured.
8. **Cache setup:** A cache is set up for Go modules and build cache.
9. **Download Dependencies:** Go modules are downloaded using go mod tidy and go mod download commands.
10. **Run Gosec:** The Gosec security scanner is run on the entire codebase.

#### Golangci-lint Job

This job runs the Golangci-lint linter on the codebase. Here's what happens during this job:

11. **Setup Go and Checkout Code:** Similar to previous jobs, Go is setup, code is checked out, and GitHub token is configured.
12. **Cache setup:** A cache is set up for Go modules and build cache.
13. **Download Dependencies:** Go modules are downloaded using go mod tidy and go mod download commands.
14. **Run Golangci-lint:** The Golangci-lint tool is run on the entire codebase.

*This workflow ensures that the codebase is secure, correctly linted, and passes all unit tests on every commit or pull request, thereby maintaining the quality and reliability of the codebase.*

### 2. End to End Test Workflow

This workflow comprises a single job named end-to-end-test. It aims to run end-to-end tests on the codebase.

The workflow is triggered on every push to the **master** or **develop** branches and on every new pull request.

Environment Variables and Secrets
- GH_ACCESS_TOKEN: GitHub access token used for cloning and downloading Go modules. Stored as a GitHub Secret.

#### End to End Test Job

This job is responsible for running end-to-end tests on the application. Here's what happens during this job:
1. **Install Go:** The appropriate Go version, specified in the matrix, is installed on the runner.
2. **Checkout Code:** The latest code is checked out from the repository.
3. **Setup Caching:** A cache is set up for Go modules and build cache. This speeds up the setup process by avoiding the re-download of dependencies that have not changed.
4. **Setup GitHub Token:** The GitHub token is configured to allow fetching of private Go modules.
5. **Build:** The code is built using the make build command. This command compiles the Go code and generates executable binaries.
6. **Start E2E Local Chain:** A local blockchain is started using the make e2e_start_localchain command. This local chain is used in the subsequent end-to-end tests. A delay of 5 seconds is introduced to ensure that the local chain has enough time to start up properly.
7. **Run E2E Test:** End-to-end tests are run using the make e2e_test command. These tests interact with the local blockchain started in the previous step, simulating real-world user interactions.

*This workflow ensures that the application functions as expected from an end-to-end perspective. It verifies the integration between different parts of the application and helps in detecting issues that unit tests might miss.*


### 3. Docker Release Workflow

The Docker Release workflow is responsible for building and pushing Docker images to Docker Hub.

The workflow is triggered when there's a new commit pushed to the **master** or **develop** branches, or when a new pull request is opened.

Here's an overview of the steps in the workflow:

1. **Check out code:** The workflow starts by checking out the latest code from the repository.
2. **Build Docker image:** Then it builds a Docker image from the checked-out code using the Dockerfile located in the repository root. During this step, a few OCI Image Format Specification labels are added to the Docker image.
3. **Log into Docker Hub registry:** After the image has been built, the workflow logs into Docker Hub using provided secrets.
4. **Tag and push Docker image to Docker Hub:** The Docker image is tagged using the commit SHA that triggered the workflow. Then, the Docker image is pushed to Docker Hub.

Environment Variables and Secrets
- DOCKER_USERNAME: Your Docker Hub username. Stored as a GitHub Secret.
- DOCKER_PASSWORD: Your Docker Hub password or token. Stored as a GitHub Secret.
- IMAGE_NAME: The name of the Docker image, including the Docker Hub username and the repository (e.g., meirgenuine/greenfield).
- IMAGE_SOURCE: The source code repository for the Docker image (e.g., https://github.com/bnb-chain/greenfield).

#### Using the Docker Image

After the workflow completes, the Docker image will be available on Docker Hub tagged with the commit SHA of the commit that triggered the workflow. The image can be pulled and run with Docker using the following commands (replace __COMMIT_SHA__ with the actual commit SHA):

```(bash)
docker pull meirgenuine/greenfield:COMMIT_SHA
```

*This workflow allows us to automatically build and publish Docker images for every commit, providing immutable and traceable deployments.*


### 4. Release Workflow

This workflow is designed to automate the process of creating software releases on GitHub. It gets triggered whenever a tag prefixed with v is pushed to the repository, such as v1.2.3.

The workflow comprises two main jobs: build and release.

#### Build Job

This job runs on both Linux (Ubuntu) and MacOS environments. It performs the following tasks:
- Checks out the repository code.
- Installs the appropriate Go version specified in the matrix.
- Caches Go dependencies to speed up subsequent workflow runs.
- Sets up the GitHub token for fetching private Go modules.
- Compiles the Go code to generate executable binaries specific to the operating system. The compiled binary is stored under ./build/bin/gnfd.
- Uploads the generated binaries as workflow artifacts. The Linux binary is uploaded under the name linux, and the MacOS binary is uploaded under the name macos.

#### Release Job

This job runs only on a Linux (Ubuntu) environment and depends on the completion of the build job. It performs the following tasks:
- Sets an environment variable RELEASE_VERSION which captures the version from the Git tag.
- Checks out the repository code.
- Downloads the linux and macos artifacts generated in the build job. The downloaded artifacts are placed in respective ./linux and ./macos directories.
- Prepares the release assets. This includes renaming the downloaded binaries and preparing a zip archive containing testnet configuration files.
- Generates a change log for the new release using a script.
- Uses the softprops/action-gh-release action to create a new GitHub release. The release includes the version tag, change log, and attached assets (Linux binary, MacOS binary, and testnet configuration zip archive).

*This workflow is essential for the consistent and efficient creation of new software releases. It ensures that each release comes with compiled binaries for supported operating systems, ready for users to download and use.*

### [Container Orchestration](#container-orchestration)
[Kubernetes](https://kubernetes.io/) would be used as solution for container orchestration.
To interact with kubernetes Greenfield-deploy uses [k8s](https://github.com/din-mukhammed/greenfield-deploy/blob/main/pkg/k8s/k8s.go) package. As MVP, Greenfield-deploy uses config file from default location (~/.kube/config) but this project also supports approach with credentials usage (see more in [k8s.go](https://github.com/din-mukhammed/greenfield-deploy/blob/main/pkg/k8s/k8s.go#L25-L32)).

### [Kubernetes manifests](#kubernetes-manifests)
Each project must have directory with kubernetes configs in `deployments/` folder. See more details in [official doc](https://kubernetes.io/docs/concepts/overview/working-with-objects/#:~:text=Understanding%20Kubernetes%20objects-,Kubernetes%20objects%20are%20persistent%20entities%20in%20the%20Kubernetes%20system.,running%20(and%20on%20which%20nodes)). The greenfield-deploy project currently supports the following k8s objects: Deployment, Service, Pod, Job, CronJob.
So, for adding new project developer has to do the following steps:
1. Describe k8s manifests in folder: `deployments/<new-project-name>/` for each environment
2. Each k8s config should start with prefix: `k8s_<environment>_`


### [Telegram Bot](#telegram-bot)

Project has a telegram bot that provides deploying operations. It allows you to run applications from created images.
The bot sends a request to the deployment service. The request contains all information about image that you want to deploy.
Before using the bot, make sure that you have permissions to perform the deployment. New users are added via a pull request. To add a new user you should create a pull request and write his telegram username to config.yaml `PROJECT_DIR/bot/config/config.yaml` \<username\> : all. 


### [Local setup](#local-setup)
1. Create kubernetes cluster and namespace only once
    - Install [kubectl](https://kubernetes.io/docs/tasks/tools/#kubectl)
    - Install [kind](https://kind.sigs.k8s.io/) tool
    - Create k8s cluster: `kind create cluster --name greenfield`
    - Create k8s namespace: `kubectl create namespace prod`
    - Check that k8s config generated in default home path: `~/.kube/config`
2. Fork greenfield-deploy. Get your personal github [access token](https://docs.github.com/en/authentication/keeping-your-account-and-data-secure/creating-a-personal-access-token)
3. Run deploy server and telegram bot.
```
make build
GITHUB_REPO="YOUR GITHUB LOGIN" GITHUB_TOKEN="YOUR TOKEN" ./greenfield-deploy web
# run telegram bot in another terminal
make bot
```
4. Go ahead to deploy via [telegram](https://t.me/test_greenfield_deployment_bot).
- Send message to bot: `/deploy greenfield latest greenfield prod production`


### [License](#license)

The Greenfield-deploy project is licensed under the [Apache License 2.0](https://github.com/deedy/Apache-License/blob/master/LICENSE)
