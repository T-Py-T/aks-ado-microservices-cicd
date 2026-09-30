# Hireability and discoverability

> Tip-cite bank: base main `eccde3f7` + ship 239 (hireability lean) — this PR pending Steward;
> provenance only; never `READY`.

This page orients staffing readers to inspectable delivery-engineering choices.
It is not an acceptance gate, scorecard, release declaration, or `READY` signal.

## What / why / how

| Question | Short answer |
| --- | --- |
| **What** | Azure DevOps delivery path around Online Boutique for operator-controlled AKS. |
| **Why** | Show reviewable build → scan → publish → manifest-update → apply ownership without hiding ops steps. |
| **How** | `azure-pipelines.yml` builds/scans eleven images; `deployment-service.yaml` carries versions; promotion stays reviewed before apply. |

## What this proves

This repository is a compact, reviewable example of delivery ownership for a
polyglot microservices application. A reviewer can inspect that the author can:

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
score.

## Stack

- **CI:** Azure Pipelines (`azure-pipelines.yml`) + Trivy filesystem/image scans
- **Deploy:** AKS via `deployment-service.yaml` (operator-controlled apply)
- **App:** eleven-service Online Boutique (Go, C#, Node.js, Python, Java, Locust)
- **Local checks:** pre-commit, language-native tests, kubeconform

## Suggested GitHub topics

`aks`, `azure`, `azure-devops`, `cicd`, `devsecops`, `kubernetes`,
`microservices`, `trivy`, `online-boutique`

Topics aid search only; they do not certify results or readiness.

## License

Repository-specific pipeline, deployment, and documentation work is under the
[MIT License](../LICENSE). Online Boutique source retains Google LLC Apache-2.0
notices; see [THIRD_PARTY_NOTICES.md](../THIRD_PARTY_NOTICES.md).

## Related docs

| Document | Role |
| --- | --- |
| [../README.md](../README.md) | Architecture, pipeline, local validation |
| [../SECURITY.md](../SECURITY.md) | Vulnerability reporting |
| [../CONTRIBUTING.md](../CONTRIBUTING.md) | Contribution and tip-cite rules |
| [OPEN_PROBLEMS.md](OPEN_PROBLEMS.md) | Held evidence boundaries |
| [README.md](README.md) | Documentation index |
