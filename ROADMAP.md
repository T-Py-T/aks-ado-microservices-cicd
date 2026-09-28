# Roadmap

> Tip-cite: base main `864c3389` + Ship 227. Steward resolves after merge; this pointer is not approval and never `READY`.

**Status:** planning document. This page lists intended, evidence-backed portfolio hygiene
work. It is not a scorecard, acceptance record, release declaration, or `READY` claim.

This repository documents an Azure DevOps delivery path for an AKS-hosted microservices
application. The items below are realistic next steps derived from
[`README.md`](README.md), [`docs/OPEN_PROBLEMS.md`](docs/OPEN_PROBLEMS.md), and the
existing CI documentation. Completing an item updates repository evidence only; it does
not establish live Azure DevOps, registry, or cluster proof without separate operator
receipts.

## Baseline provenance

| Field | Value |
| --- | --- |
| Base `main` tip | `864c3389` |
| Merge PR | Ship 227 — README/roadmap tip-cite cross-link (pending) |
| Wayfinder | [#90](https://github.com/T-Py-T/aks-ado-microservices-cicd/issues/90) — next steps after Ship 219 |
| Steward | Resolves the 8-character tip against `main`; no `READY` claim |

Ship 219 landed [`docs/README.md`](docs/README.md) at tip `864c3389`. Use
[`CHANGELOG.md`](CHANGELOG.md) for landed docs-ship history; use this roadmap for
planned work that has not yet shipped. The root README
[Evidence status](README.md#evidence-status) section carries the ship-specific
tip-cite bank; this page cross-links that bank without duplicating it.

## Planned work

Each item carries an honest status. `UNTESTED` means the work is not yet implemented or
not yet re-validated on `main`. `GAP` means repository evidence is incomplete or
stale. `BLOCKED-AUTH` means progress depends on authorized Azure DevOps, registry, or
AKS credentials outside this repository.

### Docs and tip-cite hygiene

| Item | Status | Notes |
| --- | --- | --- |
| Reconcile README tip-cite bank with current doc headers | **GAP** | Ship 227 adds README↔roadmap cross-links and recent ship pointers; remaining doc headers may still need alignment without inventing live-run evidence. |
| Record future docs-only ships in `CHANGELOG.md` | **UNTESTED** | Maintain the Ship 171 pattern: factual entries, tip-cite per merge, no readiness language. |
| Refresh `docs/OPEN_PROBLEMS.md` when held boundaries change | **GAP** | Inventory should track new gaps; closing a roadmap item here is not the same as closing an open problem. |

### CI and offline validation evidence

| Item | Status | Notes |
| --- | --- | --- |
| Document GitHub PR-check scope vs Azure Pipeline scope | **GAP** | [`.github/workflows/pr-checks.yml`](.github/workflows/pr-checks.yml) proves offline policy, language tests, and kubeconform; [`azure-pipelines.yml`](azure-pipelines.yml) build/scan/publish steps require authorized Azure DevOps execution. |
| Extend repository policy tests when pipeline contracts change | **UNTESTED** | [`tests/test_repository_policy.py`](tests/test_repository_policy.py) guards digest pins, workflow triggers, and grype exceptions; new delivery constraints should add focused tests, not narrative claims. |
| Keep pre-commit and PR-check commands aligned with README | **GAP** | Local validation commands in README should stay copy-paste accurate as checks evolve. |

### Live delivery evidence (held outside the repo)

| Item | Status | Notes |
| --- | --- | --- |
| Operator-retained Azure DevOps run receipt template | **BLOCKED-AUTH** | Document what an authorized run should record (build ID, scan summary, manifest commit) without asserting current pipeline success from source. |
| AKS promotion and rollout checklist separate from build evidence | **BLOCKED-AUTH** | Per open problems: manifest review, apply authorization, and rollout observation stay operator steps; screenshots in README are illustrative, not current-state proof. |
| 2A lane telemetry capture | **BLOCKED-AUTH / Telemetry GAP** | No invented metrics or score; telemetry remains held until an authorized environment supplies and retains receipts. |

### Portfolio alignment

| Item | Status | Notes |
| --- | --- | --- |
| Mirror portable roadmap structure for sibling EKS/Jenkins portfolio repo | **UNTESTED** | Keep section layout simple so eks-jenkins can reuse the same planning pattern without shared live-cluster claims. |

## Explicit non-claims

- No roadmap item is `READY`.
- No score, metric, production outcome, security certification, or bake-off result is asserted.
- No secrets, service connections, or live cluster access are supplied or unlocked by this page.
- No live Azure DevOps run, registry scan, or AKS rollout is inferred from planned work or repository artifacts.
- Completing a docs or CI hygiene item does not close authentication or live-evidence gaps listed in [`docs/OPEN_PROBLEMS.md`](docs/OPEN_PROBLEMS.md).
