# Contributing

> Tip-cite: main `f1aff2a9` + PR #55.

Thanks for helping improve this Azure DevOps and AKS delivery example. Keep
changes focused on the repository-specific pipeline, deployment, tests, and
documentation around the upstream Online Boutique application.

## Before opening a pull request

- Start from the current `main` branch and use a short branch name that
  describes the change.
- Keep one concern per pull request and explain which pipeline stage, service,
  or Kubernetes resource is affected.
- Do not commit credentials, registry passwords, generated images, local
  environments, or editor settings.
- When changing upstream-derived source or assets, preserve their notices and
  update `THIRD_PARTY_NOTICES.md` when provenance changes.

## Local validation

Install the repository hooks and run the policy checks:

```bash
python -m pip install pre-commit
pre-commit run --all-files
```

Run the focused service checks relevant to your change. The README documents
these examples:

```bash
for service in checkoutservice frontend productcatalogservice shippingservice; do
  (cd "src/$service" && go test ./...)
done
dotnet test src/cartservice/tests/cartservice.tests.csproj --configuration Release
```

For Node.js and Python service changes, run the documented dependency audits
and smoke checks as well. Validate Kubernetes changes offline with:

```bash
go run github.com/yannh/kubeconform/cmd/kubeconform@v0.8.0 -strict -summary deployment-service.yaml
```

## Pull requests

1. Describe the change, its operational assumptions, and the validation you
   ran.
2. Include tests or documentation updates for changed behavior.
3. For pipeline or manifest changes, call out image tags, registry settings,
   and any environment-specific values that reviewers must verify.
4. Keep secrets in Azure DevOps variable groups or service connections, never
   in the repository.
5. Wait for the pull-request checks to pass before requesting merge.
