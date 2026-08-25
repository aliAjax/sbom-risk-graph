#!/usr/bin/env sh
set -eu
base="${SBOM_BASE_URL:-http://127.0.0.1:8123}"
curl -fsS "$base/healthz" >/dev/null
curl -fsS -X POST "$base/api/v1/products" -H 'content-type: application/json' -d '{"id":"p1","tenant_id":"t1","name":"demo","version":"1.0.0"}' >/dev/null
curl -fsS -X POST "$base/api/v1/sboms" -H 'content-type: application/json' -d '{"id":"s1","product_id":"p1","format":"cyclonedx","document":{"bom-ref":"root","components":[{"bom-ref":"a","name":"lib-a","version":"1.2.3","purl":"pkg:generic/lib-a"}],"dependencies":[]}}' >/dev/null
curl -fsS "$base/api/v1/products/p1/risk" >/dev/null
echo "smoke ok"
