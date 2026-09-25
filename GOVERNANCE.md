# Governance

> Tip-cite: base `main` `6b17dd48` + Ship 139. Steward resolves after merge; this pointer is not approval and never `READY`.

This repository is an Azure DevOps and AKS delivery lab. Repository files,
checks, manifests, and pull requests are evidence of intended implementation;
they are not live-environment approval, deployment evidence, or a `READY`
declaration.

## Tip-cite protocol

1. Start work from the current `main` tip and record at least eight hexadecimal
   characters of that base commit.
2. Keep the base tip and pull request number in the change's tip-cite.
3. After merge, the Steward resolves the cited tip against `main`. A tip-cite
   is a trace pointer, not approval, certification, or `READY` evidence.

Azure DevOps runs, registry provenance, AKS context, rollout health, and
operator authorization must be verified in the environment being changed.
Keep credentials and environment-specific values in approved Azure DevOps
variable groups or service connections, not in this repository.
