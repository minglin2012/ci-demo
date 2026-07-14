# Circle Clock — Cross-Platform GUI App Design

## Summary

用 Go + [Fyne](https://fyne.io/) 实现一个圆形模拟时钟，覆盖 Windows / Linux / macOS / Android，通过现有 GitHub Actions CI 管线构建发布。

## Architecture

```
main.go → clock/widget.go → Fyne Canvas API
```

单窗口应用，核心是一个自定义 `fyne.CanvasObject`，每秒通过 `time.Ticker` 触发重绘。

## Drawing Layers (bottom → top)

| Layer | Content | Fyne API |
|-------|---------|----------|
| 1 | White circular clock face | `canvas.NewCircle()` |
| 2 | 12 hour tick marks (rotated 30°/ea) | `canvas.NewLine()` × 12 |
| 3 | Hour / Minute / Second hands | `canvas.NewLine()` × 3 |
| 4 | Center dot | `canvas.NewCircle()` |

## Hand Angle Calculation

```
secondAngle = second * 6            (360° / 60)
minuteAngle = minute * 6 + second * 0.1
hourAngle   = (hour % 12) * 30 + minute * 0.5
```

## Project Structure

```
clock/
├── main.go
├── clock/
│   └── widget.go
├── go.mod
└── clock_test.go
```

## CI Integration

在 `.github/workflows/ci-multi-platform.yml` 中新增 `build-clock` job：
- 桌面平台：Go 交叉编译（同现有 app 模式）
- Android：通过 `fyne package` 或 NDK 交叉编译生成 APK

## Platforms

- Windows (amd64)
- Linux (amd64)
- macOS (amd64, arm64)
- Android (APK)
