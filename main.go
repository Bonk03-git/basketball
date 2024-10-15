package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"image/color"
)

func main() {

	a := app.New()
	okno := a.NewWindow("Koszykówka")

	circle := canvas.NewCircle(color.RGBA{255, 165, 0, 255})
	circle.StrokeColor = color.RGBA{0, 0, 0, 255}
	circle.StrokeWidth = 2
	circle.Resize(fyne.NewSize(50, 50))

	content := container.NewVBox(widget.NewLabel("Oto koło:"), circle)

	okno.SetContent(content)
	okno.Resize(fyne.NewSize(100, 100))
	okno.ShowAndRun()
}
