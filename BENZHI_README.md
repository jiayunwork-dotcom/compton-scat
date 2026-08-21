# compton-scat — Compton 散射运动学核算命令行工具

compton-scat 是 Compton 散射运动学核算工具：给定入射光子能量与散射角的 JSON 文件，按 Δλ = λc(1−cosθ)、E' = E/(1+(E/mec²)(1−cosθ)) 计算波长移动、散射光子能量与反冲电子动能，并附 Klein–Nishina 截面与交叉规则自检。纯标准库，无网络依赖，无 cgo。

## 构建 / 运行 / 测试

```text
go build ./...
go run . kinematics example/511keV-90.json   # CLI：打印 λ、λ'、E'、Ke
go run . checks example/511keV-90.json       # 运动学交叉规则自检
go test ./...                                # 单元测试（constants / kinematics / section / cli）
```

## 评测镜像

本目录评测专用文件（勿覆盖项目自带 Dockerfile/README）：

- `benzhi.Dockerfile`
- `build_benzhi_docker.sh`
- `BENZHI_README.md`（本文件）

两种架构都要构建并进容器验证：

```bash
chmod +x build_benzhi_docker.sh
./build_benzhi_docker.sh <image-name> linux/arm64
./build_benzhi_docker.sh <image-name> linux/amd64
docker run -it <image-name>:latest
```

进容器后运行 `go build ./... && go test ./...`，再用 `go run . kinematics example/511keV-90.json` 验证 CLI 输出 λ、λ'、E'、Ke。
