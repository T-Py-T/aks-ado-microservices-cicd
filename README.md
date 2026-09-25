# AKS Microservices Delivery with Azure DevOps

[![PR Checks](https://github.com/T-Py-T/aks-ado-microservices-cicd/actions/workflows/pr-checks.yml/badge.svg?branch=main)](https://github.com/T-Py-T/aks-ado-microservices-cicd/actions/workflows/pr-checks.yml)


An Azure DevOps delivery pipeline for
[Google Cloud's Online Boutique](https://github.com/GoogleCloudPlatform/microservices-demo),
a polyglot e-commerce application made of eleven gRPC services.

The repository adds an Azure delivery path around the upstream application:
source scanning, per-service container builds, image scanning, manifest version
updates, and a Kubernetes definition for an operator-controlled AKS deployment.

![Azure delivery architecture](docs/img/CICD-Architechture.png)

## Architecture and evidence path

Use the architecture image as the map, then follow the repository paths that
carry each delivery step:

1. [`azure-pipelines.yml`](azure-pipelines.yml) shows source scanning, the
   eleven image builds, image scans, and manifest-version update.
2. [`deployment-service.yaml`](deployment-service.yaml) shows the AKS-facing
   workloads, probes, resources, ports, and service-to-service addresses.
3. [`docs/OPEN_PROBLEMS.md`](docs/OPEN_PROBLEMS.md) records evidence gaps and
   held decisions; it keeps repository evidence separate from authorized Azure
   DevOps, registry, and live-cluster evidence.

This path makes the implementation reviewable without treating diagrams,
source files, or screenshots as proof of current live deployment state.

## How the pipeline works

```text
source change
    │
    ▼
Trivy filesystem scan
    │
    ▼
build and publish 11 service images
    │
    ▼
pull and scan each image
    │
    ▼
update deployment-service.yaml with the build version
    │
    ▼
review and apply the manifest to an AKS environment
```

[`azure-pipelines.yml`](azure-pipelines.yml) contains the build and scan jobs.
[`deployment-service.yaml`](deployment-service.yaml) defines the service
deployments, ports, probes, resource requests, and service-to-service addresses.

## Application services

| Service | Language | Responsibility |
| --- | --- | --- |
| `frontend` | Go | Browser-facing store and session handling |
| `cartservice` | C# | Redis-backed cart storage |
| `productcatalogservice` | Go | Product listing and search |
| `currencyservice` | Node.js | Currency conversion |
| `paymentservice` | Node.js | Mock payment processing |
| `shippingservice` | Go | Shipping estimates |
| `emailservice` | Python | Mock order-confirmation email |
| `checkoutservice` | Go | Checkout workflow orchestration |
| `recommendationservice` | Python | Product recommendations |
| `adservice` | Java | Contextual text ads |
| `loadgenerator` | Python/Locust | Synthetic browsing and checkout traffic |

[![Online Boutique service architecture](docs/img/architecture-diagram.png)](docs/img/architecture-diagram.png)

The application code comes from Online Boutique. The Azure Pipeline,
environment promotion flow, image-version updates, and AKS integration are the
repository-specific work.

## Keep exploring

- [Open problems and held decisions](docs/OPEN_PROBLEMS.md) — active status
  inventory, not a scorecard or `READY` gate.

## Prerequisites

- an Azure DevOps project and pipeline;
- an Azure Container Registry or Docker Hub service connection;
- a target AKS cluster;
- `kubectl` access to the target cluster; and
- Trivy available in the build agent, or permission for the pipeline to install it.

The pipeline expects an Azure DevOps variable group named `Docker` and a
container registry service connection named `Docker Hub`. Update those names
and the image repository values to match your environment.

Keep registry passwords and Azure credentials in variable groups or service
connections. Do not commit them to this repository.

## Local validation

Install pre-commit, then run the repository policy checks:

```bash
python -m pip install pre-commit
pre-commit run --all-files
```

Service-specific tests live under `src/<service>/` and use each service's
native toolchain. Run the Go and .NET suites with:

```bash
for service in checkoutservice frontend productcatalogservice shippingservice; do
  (cd "src/$service" && go test ./...)
done
dotnet test src/cartservice/tests/cartservice.tests.csproj --configuration Release
```

Audit the resolved Node and Python dependency sets with:

```bash
(cd src/currencyservice && npm ci && npm audit --omit=dev)
(cd src/paymentservice && npm ci && npm audit --omit=dev)
node tests/node_service_smoke.js src/currencyservice server.js "CurrencyService gRPC server started" 17000
node tests/node_service_smoke.js src/paymentservice index.js "PaymentService gRPC server started" 15000
python -m pip install pip-audit
pip-audit -r src/emailservice/requirements.txt
pip-audit -r src/recommendationservice/requirements.txt
pip-audit -r src/loadgenerator/requirements.txt
```

The pull-request gate runs the same language-level checks before a branch can
merge. Container definitions pin every external parent image by digest. Grype
exceptions in [`.grype.yaml`](.grype.yaml) apply only to advisories whose fixes
require a prerelease runtime.

Validate the Kubernetes resources offline with:

```bash
go run github.com/yannh/kubeconform/cmd/kubeconform@v0.8.0 -strict -summary deployment-service.yaml
```

## Azure DevOps setup

1. Import or connect this repository to Azure Repos/Pipelines.
2. Create the container registry service connection.
3. Create the `Docker` variable group with `dockerUsername` and a secret
   `dockerPassword`.
4. Update the image repository names in `azure-pipelines.yml`.
5. Start the pipeline manually and review the Trivy results before promoting
   the generated image versions.
6. Review the updated manifest, then apply it to the intended AKS environment.

The manifest update job records the pipeline build ID in
`deployment-service.yaml`, giving each deployment an explicit set of service
versions. Automatic Azure Pipeline triggers are disabled because the pipeline
publishes images and writes the selected build version back to the manifest.

## Deployment

Apply the generated manifest to a configured cluster:

```bash
az aks get-credentials --resource-group <resource-group> --name <cluster>
kubectl apply --dry-run=server -f deployment-service.yaml
kubectl apply -f deployment-service.yaml
kubectl rollout status deployment/frontend
kubectl get pods,svc
```

The frontend is exposed through a Kubernetes service; the remaining services
communicate over gRPC using the DNS names declared in the manifest.

## Hireability / what this proves

This repository is a compact, reviewable example of delivery ownership for a
polyglot microservices application. It demonstrates that I can:

- map source changes to repeatable build, image-scan, publish, and
  manifest-update stages in Azure DevOps;
- describe Kubernetes delivery details—service-to-service addresses, probes,
  resource requests, image versions, and rollout checks—without hiding the
  operational steps; and
- keep delivery evidence bounded: credentials stay in service connections or
  variable groups, promotion is reviewed before apply, and the documented
  checks are reproducible from the repository.

These are implementation and documentation signals, not a claim of production
uptime, measured delivery improvement, security certification, or a readiness
score. See [`docs/OPEN_PROBLEMS.md`](docs/OPEN_PROBLEMS.md) for the active
evidence gaps and held decisions that bound this narrative; it is not a
scorecard or `READY` gate. The repository does not assert `READY`; the Steward
resolves cited tips against `main`.

## Screenshots

| Stage | Capture |
| --- | --- |
| Azure Pipeline | ![Azure Pipeline](docs/img/azure-pipelines.png) |
| Manifest update | ![Updated image versions](docs/img/yaml-updates.png) |
| Development AKS | ![Development deployment](docs/img/dev-kube.png) |
| Production AKS | ![Production deployment](docs/img/prod-kube.png) |
| Trivy filesystem scan | ![Trivy filesystem scan](docs/img/trivy-file-scan.png) |
| Trivy image scan | ![Trivy image scan](docs/img/trivy-iamge-scan.png) |

## License

Repository-specific pipeline, deployment, and documentation work is available
under the [MIT License](LICENSE). Online Boutique source files retain Google
LLC's Apache License 2.0 notices. See
[THIRD_PARTY_NOTICES.md](THIRD_PARTY_NOTICES.md).

## Evidence status

> Tip-cite bank:
> - Ship 31 / PR #66: T-Py-T/aks-ado-microservices-cicd `207755f9`
> - Ship 37 / PR #67: T-Py-T/aks-ado-microservices-cicd `d1812187`
> - Ship 49 / PR #68: T-Py-T/aks-ado-microservices-cicd `880adbd4`
> - Ship 91 / PR #73: T-Py-T/aks-ado-microservices-cicd `a75013c4`

The 2A lane remains **BLOCKED-AUTH / Telemetry GAP**. This repository makes no
`READY` claim and reports no invented score; no readiness conclusion is asserted.
The Steward resolves each 8-character tip against `main` when a full SHA is needed.
