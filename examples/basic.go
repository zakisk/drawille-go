package main

import (
	"fmt"
	"math"
	"time"

	"github.com/zakisk/drawille-go"
)

var t = time.Now()
var rad = -1

func sindata() float64 {
	rad++
	return 2 * math.Sin((math.Pi/9)*float64(rad))
}

func main() {
	s := drawille.NewCanvas(103, 30)
	s.LineColors = []drawille.Color{
		drawille.Red,
		drawille.RoyalBlue,
	}
	s.LabelColor = drawille.Purple
	s.AxisColor = drawille.SeaGreen
	s.NumDataPoints = 50
	s.HorizontalLabels = []string{}

	data := [][]float64{{}, {}}
	for x := 0; x < 8; x++ {
		data[0] = append(data[0], 3)
		data[1] = append(data[1], sindata())
		update(&s.HorizontalLabels)
	}
	fmt.Println(s.Plot(data))
	for x := 0; x < 28; x++ {
		data[0] = append(data[0], 3)
		data[1] = append(data[1], sindata())
		update(&s.HorizontalLabels)
	}
	fmt.Println(s.Plot(data))
	for x := 0; x < 18; x++ {
		data[0] = append(data[0], 3)
		data[1] = append(data[1], sindata())
		update(&s.HorizontalLabels)
	}
	fmt.Println(s.Plot(data))
}

func update(labels *[]string) {
	t = t.Add(time.Second)
	*labels = append(*labels, fmt.Sprintf("%02d:%02d:%02d", t.Hour(), t.Minute(), t.Second()))
}
