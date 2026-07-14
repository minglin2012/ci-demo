# ci-demo — GitHub Actions Multi-Platform CI Demo

演示 GitHub Actions 实现跨平台 CI 自动化构建的项目。

## 项目组成

| 目录 | 内容 | 技术 |
|------|------|------|
| `app/` | Go 命令行工具 | Go 1.21, 标准库 |
| `android/` | Android 演示应用 | Gradle 8.5, Java 17, compileSdk 34 |
| `.github/workflows/` | CI 流水线定义 | GitHub Actions |

## CI 流水线 (`ci-multi-platform.yml`)

```
Push/PR/手动 → Lint (go vet + test) → 并行构建:
                                         ├─ Windows (amd64, arm64)
                                         ├─ Linux (amd64, arm64)
                                         ├─ macOS (amd64, arm64)
                                         └─ Android (Debug + Release APK)
                                       → Artifact 上传
                                       → (可选) GitHub Release
```

- **触发:** push/PR 到 main/master, workflow_dispatch（手动）
- **桌面平台:** Go 交叉编译在 ubuntu-latest 上完成
- **Android:** ubuntu-latest + setup-android action 安装 SDK
- **Release:** 仅 workflow_dispatch + 填写版本号时创建

## 本地开发

```bash
# Go 应用
cd app && go run main.go && go test ./...

# Android
cd android && ./gradlew assembleDebug
```
