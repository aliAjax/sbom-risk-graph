# SBOM Risk Graph

纯 Go 1.23 软件物料清单风险图与可验证合规证据平台。项目接收 SPDX/CycloneDX 文档、漏洞情报和构建证明，规范化组件图并输出可解释风险；它不保存 OCI blob、不提供制品上传/下载，也不是仓库注册中心。

## Run and verify

```text
gofmt -w .
go test ./...
go vet ./...
go test -race ./...
go build ./...
go run ./cmd/riskd -addr :8123 -data-dir ./data
```

The local adapter stores canonical JSON snapshots on disk. `GET /healthz`, `GET /readyz`, `POST /api/v1/products`, `POST /api/v1/sboms`, `POST /api/v1/sboms/{id}/validate`, `POST /api/v1/recompute-jobs`, `GET /api/v1/products/{id}/risk`, and `GET /api/v1/graph/paths` form the smoke path. Production database, object store, advisory feed and signer integrations are explicit ports and return `UNIMPLEMENTED` when not injected.

## Boundaries and security

Only metadata, hashes and evidence are retained. Input size, dependency depth, graph fan-out, version syntax and JSON nesting are bounded. Every imported document receives a canonical digest; evidence is signed through a verifier port and audit records are hash chained. Tenant/product isolation is enforced before graph and risk queries.

## Layout

`internal/domain` holds entities and state machines; `parser`, `graph`, `advisory`, `policy`, `evidence`, and `waiver` are framework-independent; `repository`, `worker`, `transport/http`, and `observability` are adapters. See `docs/operations.md` for recovery, quarantine and feed replay procedures.
