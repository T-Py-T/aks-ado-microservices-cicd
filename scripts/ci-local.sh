#!/usr/bin/env bash
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
cd "$REPO_ROOT"

require_cmd() {
  if ! command -v "$1" >/dev/null 2>&1; then
    echo "ci-local: missing required command: $1" >&2
    exit 127
  fi
}

echo "==> Repository policy"
require_cmd python3
python3 -m unittest tests.test_repository_policy

echo "==> Python services (compile + pip-audit)"
require_cmd python3
python3 -m pip install --disable-pip-version-check pip-audit==2.10.1
python3 -m compileall -q src/emailservice src/loadgenerator src/recommendationservice
python3 -m pip_audit -r src/emailservice/requirements.txt
python3 -m pip_audit -r src/loadgenerator/requirements.txt
python3 -m pip_audit -r src/recommendationservice/requirements.txt

echo "==> Node.js services"
require_cmd node
require_cmd npm
for service in currencyservice paymentservice; do
  (
    cd "src/$service"
    npm ci
    npm audit --omit=dev
  )
done
node tests/node_service_smoke.js src/currencyservice server.js "CurrencyService gRPC server started" 17000
node tests/node_service_smoke.js src/paymentservice index.js "PaymentService gRPC server started" 15000

echo "==> Kubernetes manifest (kubeconform)"
require_cmd go
go run github.com/yannh/kubeconform/cmd/kubeconform@v0.8.0 -strict -summary deployment-service.yaml

echo "==> Go services"
for service in checkoutservice frontend productcatalogservice shippingservice; do
  (cd "src/$service" && go test ./...)
done

echo "==> Cart service (.NET)"
require_cmd dotnet
dotnet test src/cartservice/tests/cartservice.tests.csproj --configuration Release

echo "==> Ad service (Java)"
require_cmd java
(
  cd src/adservice
  ./gradlew --no-daemon build
)

echo "ci-local: all checks passed"
