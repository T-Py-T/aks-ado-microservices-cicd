# Support

> Tip-cite: main `63300327` + PR #72. Steward resolves after merge; no READY claim.

This repository is an Azure DevOps and AKS microservices delivery example, not a
hosted service. Support is limited to the repository's documented examples,
configuration guidance, tests, and reproducible local checks on `main`.

## Before opening an issue

- Confirm the problem still occurs on the current `main` tip.
- Include the relevant path, commit, command, and a minimal reproduction.
- Redact credentials, tokens, kubeconfigs, personal data, and private endpoint
  details from logs and configuration.
- For an Azure DevOps or AKS result, include only sanitized evidence that you
  are authorized to share; repository artifacts do not substitute for live
  environment evidence.

## Operational boundary

The repository does not provide an operated cluster, Azure DevOps project,
registry, or production support channel. Environment-specific credentials,
service connections, variable groups, and deployment decisions remain with the
authorized operator. Do not commit secrets or request them in an issue.

## Evidence boundary and non-claims

A support response or local validation result does not certify a pipeline,
cluster, image, provider, or deployment. This page asserts no score, production
outcome, security certification, or readiness conclusion. Nothing here is a
`READY` declaration, and this repository does not claim `READY`.

For unresolved evidence gaps and held decisions, see
[Open problems and held decisions](docs/OPEN_PROBLEMS.md). A tip-cite is only a
trace pointer: the Steward resolves the cited short `main` tip and PR after
merge; it is not approval and never implies `READY`.
