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
			fyne.Do(c.Refresh)
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
	faceCircle *canvas.Circle
	ticks      []*canvas.Line
	hourHand   *canvas.Line
	minuteHand *canvas.Line
	secondHand *canvas.Line
	centerDot  *canvas.Circle
	objects    []fyne.CanvasObject
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

	r.secondHand = canvas.NewLine(color.RGBA{255, 0, 0, 255})
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

// Destroy implements fyne.WidgetRenderer. No resources to release.
func (r *clockRenderer) Destroy() {}
