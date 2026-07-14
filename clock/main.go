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
