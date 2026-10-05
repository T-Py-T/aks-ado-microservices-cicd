# Governance

This repository is an Azure DevOps and AKS delivery example with a single
owner, [@T-Py-T](https://github.com/T-Py-T) (see [MAINTAINERS.md](MAINTAINERS.md)).

## How changes are made

1. Changes arrive as pull requests against `main`, one concern per pull
   request.
2. The PR Checks workflow must pass before merge.
3. The owner reviews and decides whether to merge.
4. Notable changes are recorded in [CHANGELOG.md](CHANGELOG.md). Planned work
   lives in [ROADMAP.md](ROADMAP.md), and known gaps in
   [docs/OPEN_PROBLEMS.md](docs/OPEN_PROBLEMS.md).

## What the repository can and can't show

Repository files, checks, manifests and pull requests show the intended
implementation. They aren't approval for, or proof of, a live deployment.
Azure DevOps runs, registry provenance, AKS context, rollout health and
operator authorization must be verified in the environment being changed.

Keep credentials and environment-specific values in approved Azure DevOps
variable groups or service connections, not in this repository.
