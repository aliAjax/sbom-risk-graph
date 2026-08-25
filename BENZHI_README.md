# sbom-risk-graph__014

基于 Go 实现的软件物料清单风险图 API 项目，一款后端服务，完成 SBOM 解析、依赖风险计算与合规证据验证。
## 构建镜像

请从**仓库根目录**执行；`benzhi.Dockerfile`、`build_benzhi_docker.sh`、`BENZHI_README.md` 均固定在该目录：

```bash
./build_benzhi_docker.sh <image-name> [linux/amd64|linux/arm64]
```

## 标准命令

```bash
go build ./...     # 编译
go run ./cmd/riskd   # 启动
go test ./...      # 测试（如有）
```

## 环境

- 基础镜像: golang:1.23
- Go 模块目录: `.`
- 依赖已在镜像构建阶段预下载，容器内离线可用。
- 容器内工作目录: `/app`
