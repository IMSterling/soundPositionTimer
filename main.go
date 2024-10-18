package main

import (
	"gioui.org/f32"
	"gioui.org/text"
	"gioui.org/widget/material"
	"image"
	"image/color"
	"log"
	"math"
	"os"
	"strconv"
	"strings"
	"time"

	"gioui.org/app"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
	"gioui.org/widget"
)

var progressIncrementer chan bool
var rotationsPerSecond float64

func main() {

	progressIncrementer = make(chan bool)
	rotationsPerSecond = 0.25 // default to 4 seconds per rotation
	// Animate continuously
	go func() {
		for {
			time.Sleep(time.Second / 25)
			progressIncrementer <- true
		}
	}()

	go func() {
		w := new(app.Window)
		w.Option(app.Title("Speech timer"))
		w.Option(app.Size(unit.Dp(600), unit.Dp(600)))
		if err := draw(w); err != nil {
			log.Fatal(err)
		}
		os.Exit(0)
	}()

	app.Main()
}

type C = layout.Context
type D = layout.Dimensions

func draw(w *app.Window) error {

	var ops op.Ops
	var speedButton widget.Clickable
	var secondsPerRotationInput widget.Editor
	radius := 200.0 // Radius of our dial
	th := material.NewTheme()

	// Invalidate to animate
	go func() {
		for range progressIncrementer {
			w.Invalidate()
		}
	}()
	startTime := time.Now()
	for {
		switch e := w.Event().(type) {

		case app.FrameEvent:
			gtx := app.NewContext(&ops, e)
			if speedButton.Clicked(gtx) {
				// Read from the input box
				inputString := secondsPerRotationInput.Text()
				inputString = strings.TrimSpace(inputString)
				inputFloat, _ := strconv.ParseFloat(inputString, 32)
				rotationsPerSecond = 1 / inputFloat
			}

			layout.Flex{
				Axis:    layout.Vertical,
				Spacing: layout.SpaceStart,
			}.Layout(gtx,
				layout.Rigid(
					func(gtx C) D {
						centerX, centerY := gtx.Constraints.Max.X/2, gtx.Constraints.Max.Y/8
						center := f32.Pt(float32(centerX), float32(centerY))
						elapsed := time.Since(startTime).Seconds()

						// Colored circle
						circle := clip.Ellipse{
							Min: image.Pt(centerX-int(radius), centerY-int(radius)),
							Max: image.Pt(centerX+int(radius), centerY+int(radius)),
						}.Op(gtx.Ops)
						colorPurple := color.NRGBA{R: 200, G: 162, B: 200, A: 255}
						paint.FillShape(gtx.Ops, colorPurple, circle)

						// Inner white circle
						scaleFactor := 0.97
						innerCircle := clip.Ellipse{
							Min: image.Pt(centerX-int(radius*scaleFactor), centerY-int(radius*scaleFactor)),
							Max: image.Pt(centerX+int(radius*scaleFactor), centerY+int(radius*scaleFactor)),
						}.Op(gtx.Ops)
						colorWhite := color.NRGBA{R: 255, G: 255, B: 255, A: 255}
						paint.FillShape(gtx.Ops, colorWhite, innerCircle)

						// Spinning hand
						var path clip.Path
						path.Begin(gtx.Ops)
						angle := elapsed * rotationsPerSecond * 2 * math.Pi
						end := f32.Pt(center.X-float32(radius*math.Cos(angle)), center.Y-float32(radius*math.Sin(angle)))
						path.MoveTo(center)
						path.LineTo(end)
						path.Close()

						colorBlack := color.NRGBA{R: 0, G: 0, B: 0, A: 255}
						paint.FillShape(&ops, colorBlack,
							clip.Stroke{
								Path:  path.End(),
								Width: 4,
							}.Op(),
						)
						// Cardinal hashes
						var cardinalPath clip.Path
						cardinalPath.Begin(gtx.Ops)

						type Cardinal struct {
							radius float64
							angles []float64
						}

						cardinals := []Cardinal{
							{radius / 1.2, []float64{0, math.Pi / 2, math.Pi, 3 * math.Pi / 2}},
							{radius / 1.1, []float64{math.Pi / 4, 3 * math.Pi / 4, 5 * math.Pi / 4, 7 * math.Pi / 4}},
							{radius / 1.05, []float64{math.Pi / 8, 3 * math.Pi / 8, 5 * math.Pi / 8, 7 * math.Pi / 8, 9 * math.Pi / 8, 11 * math.Pi / 8, 13 * math.Pi / 8, 15 * math.Pi / 8}},
						}

						for _, cardinal := range cardinals {
							for _, angle := range cardinal.angles {
								startLine := f32.Pt(center.X+float32(cardinal.radius*math.Cos(angle)), center.Y+float32(cardinal.radius*math.Sin(angle)))
								endLine := f32.Pt(center.X+float32(radius*math.Cos(angle)), center.Y+float32(radius*math.Sin(angle)))
								cardinalPath.MoveTo(startLine)
								cardinalPath.LineTo(endLine)
							}
						}

						paint.FillShape(&ops, colorPurple,
							clip.Stroke{
								Path:  cardinalPath.End(),
								Width: 4,
							}.Op(),
						)

						d := image.Point{X: gtx.Constraints.Max.X / 2, Y: gtx.Constraints.Max.Y / 2}
						return layout.Dimensions{Size: d}
					},
				),

				layout.Rigid(
					func(gtx C) D {

						secondsPerRotationInput.SingleLine = true
						secondsPerRotationInput.Alignment = text.Middle

						margins := layout.Inset{
							Top:    unit.Dp(0),
							Right:  unit.Dp(170),
							Bottom: unit.Dp(40),
							Left:   unit.Dp(170),
						}
						border := widget.Border{
							Color:        color.NRGBA{R: 204, G: 204, B: 204, A: 255},
							CornerRadius: unit.Dp(3),
							Width:        unit.Dp(2),
						}
						ed := material.Editor(th, &secondsPerRotationInput, "Seconds per rotation")
						return margins.Layout(gtx,
							func(gtx C) D {
								return border.Layout(gtx, ed.Layout)
							},
						)
					},
				),

				layout.Rigid(
					func(gtx C) D {
						margins := layout.Inset{
							Top:    unit.Dp(25),
							Bottom: unit.Dp(25),
							Right:  unit.Dp(35),
							Left:   unit.Dp(35),
						}
						return margins.Layout(gtx,
							func(gtx C) D {
								btn := material.Button(th, &speedButton, "Change speed")
								return btn.Layout(gtx)
							},
						)
					},
				),
			)
			e.Frame(gtx.Ops)
			
		case app.DestroyEvent:
			return e.Err
		}

	}
}
