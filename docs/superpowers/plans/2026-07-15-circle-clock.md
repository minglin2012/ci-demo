# Circle Clock App Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** Build a cross-platform analog circle clock GUI app using Go + Fyne, with CI builds for Windows/Linux/macOS/Android.

**Architecture:** Single-window Fyne app with a custom `Clock` widget. The widget uses a `time.Ticker` (1s interval) to refresh and Fyne canvas primitives (`Circle`, `Line`) to draw the clock face, tick marks, and three hands.

**Tech Stack:** Go 1.21+, Fyne v2.x, Fyne canvas API

## Global Constraints

- Go version >= 1.21
- Fyne v2 (latest stable)
- Desktop builds via standard Go cross-compilation
- Android build via `fyne package -os android` (requires NDK)
- All angle calculation logic must be in a separately testable function
- Widget must use Fyne canvas primitives (no images or external assets)

---

### Task 1: Scaffold Clock Module

**Files:**
- Create: `clock/go.mod`
- Create: `clock/main.go`
- Create: `clock/clock/` (empty dir, `mkdir -p`)

**Interfaces:**
- Consumes: nothing
- Produces: `go.mod` with module name `github.com/demo/ci-demo-clock`, dependency on `fyne.io/fyne/v2`

- [ ] **Step 1: Create directory structure**

```bash
cd d:/programming/ci-demo
mkdir -p clock/clock
```

- [ ] **Step 2: Initialize Go module**

```bash
cd d:/programming/ci-demo/clock
go mod init github.com/demo/ci-demo-clock
```

- [ ] **Step 3: Add Fyne dependency**

```bash
cd d:/programming/ci-demo/clock
go get fyne.io/fyne/v2@latest
```

Expected: downloads Fyne and its transitive dependencies, updates `go.mod` and `go.sum`.

- [ ] **Step 4: Write minimal main.go to verify imports compile**

Write `clock/main.go`:

```go
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("Circle Clock")
	w.SetContent(widget.NewLabel("Clock placeholder"))
	w.Resize(fyne.NewSize(300, 350))
	w.ShowAndRun()
}
```

- [ ] **Step 5: Verify compilation**

```bash
cd d:/programming/ci-demo/clock
go build ./...
```

Expected: builds successfully with no errors.

- [ ] **Step 6: Commit**

```bash
cd d:/programming/ci-demo
git add clock/go.mod clock/go.sum clock/main.go
git commit -m "feat: scaffold clock Go module with Fyne dependency"
```

---

### Task 2: Write Angle Calculation Logic + Tests

**Files:**
- Create: `clock/clock/hand_angles.go`
- Create: `clock/clock/hand_angles_test.go`

**Interfaces:**
- Consumes: nothing
- Produces:
  - `func CalculateAngles(t time.Time) (hourAngle, minuteAngle, secondAngle float64)`
  - test file with table-driven tests

- [ ] **Step 1: Write hand_angles.go**

Write `clock/clock/hand_angles.go`:

```go
package clock

import "time"

// CalculateAngles returns the angles (in degrees) for hour, minute, and second hands
// based on the given time. Angles follow standard clock convention:
// 0° = 12 o'clock, rotating clockwise.
//
//	secondAngle = second * 6            (360° / 60)
//	minuteAngle = minute * 6 + second * 0.1
//	hourAngle   = (hour % 12) * 30 + minute * 0.5
func CalculateAngles(t time.Time) (hourAngle, minuteAngle, secondAngle float64) {
	hour := t.Hour()
	minute := t.Minute()
	second := t.Second()

	secondAngle = float64(second) * 6.0
	minuteAngle = float64(minute)*6.0 + float64(second)*0.1
	hourAngle = float64(hour%12)*30.0 + float64(minute)*0.5

	return hourAngle, minuteAngle, secondAngle
}
```

- [ ] **Step 2: Write test file**

Write `clock/clock/hand_angles_test.go`:

```go
package clock

import (
	"math"
	"testing"
	"time"
)

func TestCalculateAngles(t *testing.T) {
	tests := []struct {
		name                     string
		timeStr                  string
		wantHour, wantMin, wantSec float64
	}{
		{
			name:    "12:00:00 — all hands at 12",
			timeStr: "2026-07-15T12:00:00",
			wantHour: 0, wantMin: 0, wantSec: 0,
		},
		{
			name:    "12:15:00 — quarter past",
			timeStr: "2026-07-15T12:15:00",
			wantHour: 7.5, wantMin: 90, wantSec: 0,
		},
		{
			name:    "12:00:30 — 30 seconds",
			timeStr: "2026-07-15T12:00:30",
			wantHour: 0, wantMin: 3, wantSec: 180,
		},
		{
			name:    "6:00:00 — 6 o'clock",
			timeStr: "2026-07-15T06:00:00",
			wantHour: 180, wantMin: 0, wantSec: 0,
		},
		{
			name:    "3:30:00 — half past three",
			timeStr: "2026-07-15T03:30:00",
			wantHour: 105, wantMin: 180, wantSec: 0,
		},
		{
			name:    "00:00:00 — midnight",
			timeStr: "2026-07-15T00:00:00",
			wantHour: 0, wantMin: 0, wantSec: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			parsed, err := time.Parse(time.RFC3339, tt.timeStr)
			if err != nil {
				t.Fatalf("failed to parse time: %v", err)
			}

			h, m, s := CalculateAngles(parsed)

			if !almostEqual(h, tt.wantHour) {
				t.Errorf("hour angle = %.1f, want %.1f", h, tt.wantHour)
			}
			if !almostEqual(m, tt.wantMin) {
				t.Errorf("minute angle = %.1f, want %.1f", m, tt.wantMin)
			}
			if !almostEqual(s, tt.wantSec) {
				t.Errorf("second angle = %.1f, want %.1f", s, tt.wantSec)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	return math.Abs(a-b) < 0.001
}

func TestCalculateAnglesContinuousSecond(t *testing.T) {
	// Verify minute hand moves 0.1° per second
	base, _ := time.Parse(time.RFC3339, "2026-07-15T12:00:00")
	for sec := 0; sec < 60; sec++ {
		tm := base.Add(time.Duration(sec) * time.Second)
		_, m, s := CalculateAngles(tm)

		expectedSec := float64(sec) * 6.0
		expectedMin := float64(sec) * 0.1

		if !almostEqual(s, expectedSec) {
			t.Errorf("second %d: second angle = %.1f, want %.1f", sec, s, expectedSec)
		}
		if !almostEqual(m, expectedMin) {
			t.Errorf("second %d: minute angle = %.1f, want %.1f", sec, m, expectedMin)
		}
	}
}
```

- [ ] **Step 3: Run tests to verify they fail**

```bash
cd d:/programming/ci-demo/clock
go test ./clock/... -v
```

Expected: FAIL — `CalculateAngles` not yet defined.

- [ ] **Step 4: Run tests to verify they pass (after Step 1 code is in place)**

```bash
cd d:/programming/ci-demo/clock
go test ./clock/... -v
```

Expected: all tests PASS.

- [ ] **Step 5: Commit**

```bash
cd d:/programming/ci-demo
git add clock/clock/hand_angles.go clock/clock/hand_angles_test.go
git commit -m "feat: add clock hand angle calculation with tests"
```

---

### Task 3: Implement Clock Widget with Fyne Canvas

**Files:**
- Create: `clock/clock/widget.go`

**Interfaces:**
- Consumes: `CalculateAngles(t time.Time) (hourAngle, minuteAngle, secondAngle float64)` from Task 2
- Produces: `type ClockWidget struct` (implements `fyne.Widget`), renderer type `clockRenderer` (implements `fyne.WidgetRenderer`)

- [ ] **Step 1: Write widget.go**

Write `clock/clock/widget.go`:

```go
package clock

import (
	"image/color"
	"math"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/widget"
)

// Default styling constants.
const (
	faceRadius   = 120.0
	centerRadius = 6.0
	tickLength   = 14.0

	hourHandLength   = 60.0
	minuteHandLength = 90.0
	secondHandLength = 100.0

	hourHandWidth   = 5.0
	minuteHandWidth = 3.0
	secondHandWidth = 1.5
)

// ClockWidget is a custom Fyne widget that renders an analog clock.
type ClockWidget struct {
	widget.BaseWidget
	ticker  *time.Ticker
	current time.Time
}

// NewClockWidget creates a new analog clock widget.
func NewClockWidget() *ClockWidget {
	c := &ClockWidget{
		current: time.Now(),
	}
	c.ExtendBaseWidget(c)

	// Start the ticker to refresh every second.
	c.ticker = time.NewTicker(time.Second)
	go func() {
		for t := range c.ticker.C {
			c.current = t
			c.Refresh()
		}
	}()

	return c
}

// CreateRenderer implements fyne.Widget.
func (c *ClockWidget) CreateRenderer() fyne.WidgetRenderer {
	r := &clockRenderer{
		clock: c,
	}
	r.buildObjects()
	return r
}

// Stop releases the ticker. Call when the widget is no longer needed.
func (c *ClockWidget) Stop() {
	if c.ticker != nil {
		c.ticker.Stop()
	}
}

// MinSize returns a fixed minimum size big enough for the clock face.
func (c *ClockWidget) MinSize() fyne.Size {
	diameter := float32(faceRadius*2 + 20)
	return fyne.NewSize(diameter, diameter)
}

// clockRenderer draws the analog clock face and hands.
type clockRenderer struct {
	clock *ClockWidget

	// canvas objects — rebuilt on each refresh
	faceCircle  *canvas.Circle
	ticks       []*canvas.Line
	hourHand    *canvas.Line
	minuteHand  *canvas.Line
	secondHand  *canvas.Line
	centerDot   *canvas.Circle
	objects     []fyne.CanvasObject
}

func (r *clockRenderer) buildObjects() {
	cx := float32(faceRadius + 10)
	cy := float32(faceRadius + 10)

	// --- face background ---
	r.faceCircle = canvas.NewCircle(color.White)
	r.faceCircle.StrokeColor = color.Black
	r.faceCircle.StrokeWidth = 3
	r.faceCircle.Resize(fyne.NewSize(float32(faceRadius*2), float32(faceRadius*2)))
	r.faceCircle.Move(fyne.NewPos(10, 10))

	// --- hour tick marks (12 lines around the perimeter) ---
	r.ticks = make([]*canvas.Line, 12)
	for i := 0; i < 12; i++ {
		angle := float64(i) * 30.0 * (math.Pi / 180.0)
		innerR := faceRadius - tickLength
		outerR := faceRadius

		x1 := cx + float32(innerR*math.Sin(angle))
		y1 := cy - float32(innerR*math.Cos(angle))
		x2 := cx + float32(outerR*math.Sin(angle))
		y2 := cy - float32(outerR*math.Cos(angle))

		line := canvas.NewLine(color.Black)
		line.StrokeWidth = 3
		line.Position1 = fyne.NewPos(x1, y1)
		line.Position2 = fyne.NewPos(x2, y2)
		r.ticks[i] = line
	}

	// --- hands (positions will be updated each frame) ---
	r.hourHand = canvas.NewLine(color.Black)
	r.hourHand.StrokeWidth = float32(hourHandWidth)

	r.minuteHand = canvas.NewLine(color.Black)
	r.minuteHand.StrokeWidth = float32(minuteHandWidth)

	r.secondHand = canvas.NewLine(color.Red)
	r.secondHand.StrokeWidth = float32(secondHandWidth)

	// --- center dot ---
	r.centerDot = canvas.NewCircle(color.Black)
	r.centerDot.Resize(fyne.NewSize(float32(centerRadius*2), float32(centerRadius*2)))
	r.centerDot.Move(fyne.NewPos(cx-float32(centerRadius), cy-float32(centerRadius)))

	// collect all objects (order = z-order, last = top)
	r.objects = []fyne.CanvasObject{r.faceCircle}
	for _, t := range r.ticks {
		r.objects = append(r.objects, t)
	}
	r.objects = append(r.objects, r.hourHand, r.minuteHand, r.secondHand, r.centerDot)
}

func (r *clockRenderer) updateHands() {
	cx := float32(faceRadius + 10)
	cy := float32(faceRadius + 10)

	h, m, s := CalculateAngles(r.clock.current)
	hRad := h * (math.Pi / 180.0)
	mRad := m * (math.Pi / 180.0)
	sRad := s * (math.Pi / 180.0)

	// hour hand
	hx := cx + float32(hourHandLength*math.Sin(hRad))
	hy := cy - float32(hourHandLength*math.Cos(hRad))
	r.hourHand.Position1 = fyne.NewPos(cx, cy)
	r.hourHand.Position2 = fyne.NewPos(hx, hy)

	// minute hand
	mx := cx + float32(minuteHandLength*math.Sin(mRad))
	my := cy - float32(minuteHandLength*math.Cos(mRad))
	r.minuteHand.Position1 = fyne.NewPos(cx, cy)
	r.minuteHand.Position2 = fyne.NewPos(mx, my)

	// second hand
	sx := cx + float32(secondHandLength*math.Sin(sRad))
	sy := cy - float32(secondHandLength*math.Cos(sRad))
	r.secondHand.Position1 = fyne.NewPos(cx, cy)
	r.secondHand.Position2 = fyne.NewPos(sx, sy)
}

func (r *clockRenderer) Layout(size fyne.Size) {
	r.updateHands()
}

func (r *clockRenderer) MinSize() fyne.Size {
	return r.clock.MinSize()
}

func (r *clockRenderer) Refresh() {
	r.updateHands()
}

func (r *clockRenderer) Objects() []fyne.CanvasObject {
	return r.objects
}

func (r *clockRenderer) Destroy() {}
```

- [ ] **Step 2: Verify compilation**

```bash
cd d:/programming/ci-demo/clock
go build ./...
```

Expected: builds successfully.

- [ ] **Step 3: Commit**

```bash
cd d:/programming/ci-demo
git add clock/clock/widget.go
git commit -m "feat: implement analog clock widget with Fyne canvas"
```

---

### Task 4: Wire Up main.go and Verify Full App Builds

**Files:**
- Modify: `clock/main.go`

**Interfaces:**
- Consumes: `clock.NewClockWidget()` from Task 3
- Produces: runnable desktop app

- [ ] **Step 1: Replace main.go with clock app**

Write `clock/main.go` (overwrite placeholder):

```go
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/theme"

	"github.com/demo/ci-demo-clock/clock"
)

func main() {
	a := app.New()
	a.Settings().SetTheme(theme.LightTheme())

	w := a.NewWindow("Circle Clock")
	w.SetPadded(false)

	clockWidget := clock.NewClockWidget()
	defer clockWidget.Stop()

	// Center the clock widget without extra padding.
	size := clockWidget.MinSize()
	w.Resize(size)
	w.SetFixedSize(true)
	w.SetContent(clockWidget)

	w.ShowAndRun()
}
```

- [ ] **Step 2: Verify full build**

```bash
cd d:/programming/ci-demo/clock
go build ./...
go vet ./...
```

Expected: builds + passes vet.

- [ ] **Step 3: Run all tests**

```bash
cd d:/programming/ci-demo/clock
go test ./... -v
```

Expected: all tests PASS.

- [ ] **Step 4: Commit**

```bash
cd d:/programming/ci-demo
git add clock/main.go
git commit -m "feat: wire up main.go with clock widget"
```

---

### Task 5: Add CI Build Job for Clock

**Files:**
- Modify: `.github/workflows/ci-multi-platform.yml`
- After the `build-android` job and before `release`, insert a new `build-clock` block.

**Interfaces:**
- Consumes: `clock/` module structure from Task 4
- Produces: artifact `clock-windows-amd64`, `clock-linux-amd64`, `clock-macos-amd64`, `clock-macos-arm64`

Note: Android APK for Fyne requires NDK + `fyne` CLI. This plan adds a start but it may need debugging like `build-android` did. Desktop builds work identically to existing `build-windows/linux/macos` jobs.

- [ ] **Step 1: Add clock build jobs to CI workflow**

Insert after the existing `build-android` job (after line ~221) and before `release`:

```yaml
  # =========================================================================
  # 7. Clock GUI App — Cross-Platform Build
  # =========================================================================
  build-clock-windows:
    name: 🪟 Clock — Windows
    needs: lint
    runs-on: ubuntu-latest
    strategy:
      matrix:
        arch: [amd64]
    defaults:
      run:
        working-directory: ./clock
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go ${{ env.GO_VERSION }}
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Build Clock (Windows ${{ matrix.arch }})
        run: |
          GOOS=windows GOARCH=${{ matrix.arch }} go build \
            -ldflags="-s -w" \
            -o ../dist/ci-demo-clock-windows-${{ matrix.arch }}.exe \
            ./main.go

      - name: Upload artifact (Clock Windows ${{ matrix.arch }})
        uses: actions/upload-artifact@v4
        with:
          name: clock-windows-${{ matrix.arch }}
          path: dist/ci-demo-clock-windows-${{ matrix.arch }}.exe
          retention-days: 7

  build-clock-linux:
    name: 🐧 Clock — Linux
    needs: lint
    runs-on: ubuntu-latest
    strategy:
      matrix:
        arch: [amd64]
    defaults:
      run:
        working-directory: ./clock
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go ${{ env.GO_VERSION }}
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Build Clock (Linux ${{ matrix.arch }})
        run: |
          GOOS=linux GOARCH=${{ matrix.arch }} go build \
            -ldflags="-s -w" \
            -o ../dist/ci-demo-clock-linux-${{ matrix.arch }} \
            ./main.go

      - name: Upload artifact (Clock Linux ${{ matrix.arch }})
        uses: actions/upload-artifact@v4
        with:
          name: clock-linux-${{ matrix.arch }}
          path: dist/ci-demo-clock-linux-${{ matrix.arch }}
          retention-days: 7

  build-clock-macos:
    name: 🍎 Clock — macOS
    needs: lint
    runs-on: ubuntu-latest
    strategy:
      matrix:
        arch: [amd64, arm64]
    defaults:
      run:
        working-directory: ./clock
    steps:
      - uses: actions/checkout@v4

      - name: Setup Go ${{ env.GO_VERSION }}
        uses: actions/setup-go@v5
        with:
          go-version: ${{ env.GO_VERSION }}

      - name: Build Clock (macOS ${{ matrix.arch }})
        run: |
          GOOS=darwin GOARCH=${{ matrix.arch }} go build \
            -ldflags="-s -w" \
            -o ../dist/ci-demo-clock-macos-${{ matrix.arch }} \
            ./main.go

      - name: Upload artifact (Clock macOS ${{ matrix.arch }})
        uses: actions/upload-artifact@v4
        with:
          name: clock-macos-${{ matrix.arch }}
          path: dist/ci-demo-clock-macos-${{ matrix.arch }}
          retention-days: 7
```

- [ ] **Step 2: Update release job needs list**

In the `release` job, add the new clock jobs to `needs:`:

```yaml
    needs:
      - build-windows
      - build-linux
      - build-macos
      - build-android
      - build-clock-windows
      - build-clock-linux
      - build-clock-macos
```

And update the release body table to include clock artifacts:

```yaml
          body: |
            ## 🚀 Release ${{ github.event.inputs.release_version }}

            Automated multi-platform build via GitHub Actions.

            ### CLI Artifacts
            | Platform | Architectures |
            |----------|--------------|
            | 🪟 Windows | amd64, arm64 |
            | 🐧 Linux | amd64, arm64 |
            | 🍎 macOS | amd64, arm64 (Apple Silicon) |
            | 🤖 Android | APK (debug + release) |

            ### Clock GUI App
            | 🪟 Windows | amd64 |
            | 🐧 Linux | amd64 |
            | 🍎 macOS | amd64, arm64 |

            **Commit:** ${{ github.sha }}
```

- [ ] **Step 3: Verify CI syntax locally (optional)**

```bash
# No local CI validator available, but review the YAML structure
cd d:/programming/ci-demo
cat .github/workflows/ci-multi-platform.yml | head -30
```

- [ ] **Step 4: Commit**

```bash
cd d:/programming/ci-demo
git add .github/workflows/ci-multi-platform.yml
git commit -m "ci: add clock GUI app build jobs (Windows/Linux/macOS)"
```

---

### Task 6: Verify End-to-End — Run Clock Locally (manual)

**No file changes.** This is a manual verification task.

- [ ] **Step 1: Build and run the clock app locally**

```bash
cd d:/programming/ci-demo/clock
go run ./main.go
```

Expected: a window opens showing a circular analog clock with moving hands.

- [ ] **Step 2: Verify tests still pass**

```bash
cd d:/programming/ci-demo/clock
go test ./... -v
```

Expected: all PASS.

- [ ] **Step 3: Push all commits and verify CI on GitHub**

```bash
cd d:/programming/ci-demo
git push origin master
```

Go to GitHub Actions → 🚀 Multi-Platform CI Build → verify `build-clock-*` jobs pass.
```

