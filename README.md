# compton-scat — Compton 散射运动学核算命令行工具

compton-scat 是一个 Compton 散射运动学核算工具：给定入射光子能量与散射角（JSON 文件），算出波长移动 Δλ = λc(1−cosθ)、散射光子能量 E' = E/(1+(E/mec²)(1−cosθ)) 与反冲电子动能 Ke = E − E'，并打印入射波长 λ、散射波长 λ'、E' 与 Ke。所有公式共用同一组 h、me、c（CODATA），保证 λc = h/(mec) 与 mec² 自洽、能量守恒 Ke+E'=E 恒成立；同时提供 Klein–Nishina 微分/总截面（总截面随能量升高而下降）、波长–能量互核（λ=hc/E、λ'=λ+Δλ、E'=hc/λ'）与交叉规则自检。能力边界为单次光子–电子碰撞运动学，不做剂量积分、辐射地图或能谱模拟。

## 用法

```text
go run . kinematics example/511keV-90.json
```

打印该算例（511 keV 光子、90° 散射）的 λ = 0.00242631 nm、λ' = 0.00485262 nm、E' ≈ 255.5 keV、Ke ≈ 255.5 keV，能量守恒 Ke+E' 与 E 一致。`kinematics` 换成任意事件 JSON 文件即可。

其他子命令：

```text
go run . wavelength example/511keV-90.json  # 波长与能量的互相核对
go run . section 511                        # Klein-Nishina 截面（能量单位 keV）
go run . trend                              # 总截面随能量变化（1 keV～10 MeV，单调下降）
go run . checks example/511keV-90.json      # 运动学交叉规则自检
go run . section-checks                     # 截面性质自检
go run . constants                          # 打印所用物理常数
go run . help
```

算例：`example/511keV-90.json`（511 keV、90°，E' 满足闭式 E/(1+E/mec²)）、`example/100keV-60.json`、`example/1MeV-180.json`（背散射，Δλ=2λc）、`example/theta0.json`（前向，Δλ=0、E'=E）、`example/20keV-120.json`。事件 JSON 字段为 `name`（可选）、`energy_keV`（keV）、`angle_deg`（度）；未知字段会被拒绝。

## 关键约定

- **常数**：h、me、c、e、ε₀ 全部来自 CODATA 2018；λc、mec²、hc、经典电子半径 r_e 都由这些常数推导，不另存数值。
- **角度**：散射角必须落在 [0, 180] 度；θ=0 时 Δλ=0、E'=E、Ke=0，θ=180° 时 Δλ=2λc、E' 最小。越界直接报错（stderr + 非零退出）。
- **能量**：入射能量必须为正；能量守恒 Ke+E'=E 与动量守恒（余弦定律 + 相对论色散关系）都作为内置核对。
- **闭式**：90° 散射时 E' = E/(1+E/mec²) 的闭式被 checks 命令验证。
- **波长互核**：λ=hc/E、λ'=λ+Δλ、E'=hc/λ' 三条恒等式同时成立。
- **截面**：Klein–Nishina 微分截面用钉死的公式，总截面用复合 Simpson 积分（默认 2000 段），低能趋近 Thomson 截面 8πr_e²/3，随能量升高单调下降。

非法输入（能量非正、角度越界、缺字段、未知字段、文件不存在）都以明文写 stderr 并以非零码退出；不会静默返回错误数值。

## 构建与测试

```text
go build ./...
go test ./...
```

纯标准库，无第三方依赖，无 cgo。

## 许可

MIT，见 [LICENSE](./LICENSE)。
