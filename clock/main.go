package main

import (
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
