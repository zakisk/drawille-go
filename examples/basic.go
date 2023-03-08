package main

import (
	"fmt"
	"math"
	"time"

	"github.com/chriskim06/drawille-go"
)

var t = time.Now()
var rad = -1

func sindata() float64 {
	rad++
	return 2 * math.Sin((math.Pi/9)*float64(rad))
}

func main() {
	s := drawille.NewCanvas(103, 25)
	s.LineColors = []drawille.Color{
		drawille.Red,
		drawille.RoyalBlue,
	}
	s.LabelColor = drawille.Purple
	s.AxisColor = drawille.SeaGreen
	s.NumDataPoints = 50

	i := 0
	labels := []string{}
	data := [][]float64{{}, {}}
	for x := 0; x < 16; x++ {
		data[0] = append(data[0], 3)
		data[1] = append(data[1], sindata())
		update(i, &labels)
		i++
	}
	fmt.Println(data[1])
	s.HorizontalLabels = labels
	fmt.Print(s.Plot(data))
	fmt.Println()
	for x := 0; x < 20; x++ {
		data[0] = append(data[0], 3)
		data[1] = append(data[1], sindata())
		update(i, &labels)
		i++
	}
	s.HorizontalLabels = labels
	fmt.Print(s.Plot(data))
	fmt.Println()
	for x := 0; x < 18; x++ {
		data[0] = append(data[0], 3)
		data[1] = append(data[1], sindata())
		update(i, &labels)
		i++
	}
	s.HorizontalLabels = labels
	fmt.Print(s.Plot(data))
	fmt.Println()
	//	for x := 0; x < 7; x++ {
	//	    data[0] = append(data[0], 256)
	//	    data[1] = append(data[1], 17)
	//	    update(i, &labels)
	//	    i++
	//	}
	//
	// labels = labels[7:]
	// data[0] = data[0][7:]
	// data[1] = data[1][7:]
	//
	//	for x := 0; x < 18; x++ {
	//	    data[0] = append(data[0], 355)
	//	    data[1] = append(data[1], 17)
	//	    update(i, &labels)
	//	    i++
	//	}
	//
	// s.HorizontalLabels = labels
	// fmt.Print(s.Plot(data))
	// fmt.Println()
}

func update(i int, labels *[]string) {
	ti := t.Add(time.Second)
	t = ti
	*labels = append(*labels, fmt.Sprintf("%02d:%02d:%02d", ti.Hour(), ti.Minute(), ti.Second()))
}
