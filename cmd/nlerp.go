//go:build ignore

package main

import (
	"fmt"

	g2d "github.com/jphsd/graphics2d"
	"github.com/jphsd/graphics2d/color"
	"github.com/jphsd/graphics2d/image"
	"github.com/jphsd/nonlinear"
)

var (
	table = []struct {
		nl   nonlinear.NonLinear
		name string
	}{
		{&nonlinear.NLLinear{}, "Linear"},
		{&nonlinear.NLSquare{}, "Square"},
		{&nonlinear.NLCube{}, "Cube"},
		{&nonlinear.NLCircle1{}, "Circle1"},
		{&nonlinear.NLCircle2{}, "Circle2"},
		{&nonlinear.NLSin{}, "Sin"},
		{&nonlinear.NLSin1{}, "Sin1"},
		{&nonlinear.NLSin2{}, "Sin2"},
		{nonlinear.NewNLLame(2, 2), "Lame 1"},
		{nonlinear.NewNLLame(4, 4), "Lame 2"},
		{nonlinear.NewNLLame(8, 8), "Lame 3"},
		{nonlinear.NewNLLame(.5, .5), "Lame 4"},
		{nonlinear.NewNLLame(.25, .25), "Lame 5"},
		{nonlinear.NewNLLame(.125, .125), "Lame 6"},
		{nonlinear.NewNLLame(.5, 2), "Lame 7"},
		{nonlinear.NewNLLame(2, .5), "Lame 8"},
		{nonlinear.NewNLExponential(1), "Exponential 1"},
		{nonlinear.NewNLExponential(10), "Exponential 2"},
		{nonlinear.NewNLExponential(100), "Exponential 3"},
		{nonlinear.NewNLLogarithmic(1), "Logarithmic 1"},
		{nonlinear.NewNLLogarithmic(10), "Logarithmic 2"},
		{nonlinear.NewNLLogarithmic(100), "Logarithmic 3"},
		{nonlinear.NewNLGauss(1), "Gauss 1"},
		{nonlinear.NewNLGauss(3), "Gauss 2"},
		{nonlinear.NewNLGauss(6), "Gauss 3"},
		{nonlinear.NewNLLogistic(1, 0.5), "Logistic 1"},
		{nonlinear.NewNLLogistic(12, 0.5), "Logistic 2"},
		{nonlinear.NewNLLogistic(60, 0.5), "Logistic 3"},
		{nonlinear.NewNLLogistic(1, 0.2), "Logistic 4"},
		{nonlinear.NewNLLogistic(12, 0.2), "Logistic 5"},
		{nonlinear.NewNLLogistic(38, 0.2), "Logistic 6"},
		{nonlinear.NewNLLogistic(1, 0.8), "Logistic 7"},
		{nonlinear.NewNLLogistic(12, 0.8), "Logistic 8"},
		{nonlinear.NewNLLogistic(100, 0.8), "Logistic 9"},
		{&nonlinear.NLP3{}, "P3"},
		{&nonlinear.NLP5{}, "P5"},
		{nonlinear.NewNLStopped([][]float64{
			{0.24, 0.1},
			{0.25, 0.3},
			{0.49, 0.4},
			{0.5, 0.6},
			{0.74, 0.7},
			{0.75, 0.9},
		}), "NLStopped"},
	}
)

// Graph each non-linear func with 100 steps
func main() {
	for i, row := range table {
		graph(row.nl, i, row.name)
	}
}

func graph(nl nonlinear.NonLinear, n int, s string) {
	// Create image to write into
	width, height := 1000, 1000
	img := image.NewRGBA(width, height, color.White)
	g2d.FillPath(img, g2d.Square([]float64{500, 500}, 800), g2d.LightGrayPen)

	// Plot graph in 800x800 box

	xfm := g2d.Translate(100, 900)
	xfm.Scale(800, -800)

	path := g2d.NewPath([]float64{0, 0})
	dt := 1.0 / 100
	t := dt
	for i := 0; i < 100; i++ {
		v := nl.Transform(t)
		path.AddSteps([][]float64{{t, v}})
		//fmt.Printf("%f => %f\n", t, v)
		t += dt
	}

	g2d.DrawPath(img, path.Transform(xfm), g2d.BlackPen)

	// Capture image output
	image.SaveImage(img, fmt.Sprintf("nlerp-%d %s", n, s))
}
