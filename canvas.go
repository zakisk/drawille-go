package drawille

import (
	"fmt"
	"image"
	"math"
	"strings"
)

// Canvas is a plot of braille characters
type Canvas struct {
	// various settings accessible outside this object
	LineColors []Color
	LabelColor Color
	AxisColor  Color
	ShowAxis   bool

	// a list of labels the canvas will print for the x and y axis
	// horizontal labels must be provided by the caller. too lazy
	// to come up with a good way to print stuff so offloading some
	// of that work to the user. when the horizontal labels arent
	// provided an empty line is printed
	HorizontalLabels []string

	// this value is used to determine the horizontal scale when
	// graphing points in the plot. with braille characters, each
	// cell represents 2 points along the x axis. so if the total
	// graphable area in the canvas is 50 cells then we can plot
	// 100 points of data. if we only want to plot 50 data points
	// then the horizontal scale needs to be 2 in order to utilize
	// the whole graphable area. since image uses ints for points
	// the graphing gets a little weird. ex:
	// graphable area = 88
	// num of data points = 50
	// horizontal scale = 88/50 = 1.76
	// 0, 1.76, 3.52, 5.28, 7.04, 8.8, 10.56, 12.32, 14.08
	// this would map points at 0,1,3,5,7,8,10,12,14
	NumDataPoints int

	// the bounds of the canvas
	area image.Rectangle

	plotWidth   int
	graphHeight int

	// a map of the entire braille grid
	points map[image.Point]Cell

	horizontalOffset int
	horizontalScale  float64
	maxX             int
}

// NewCanvas creates a default canvas
func NewCanvas(width, height int) Canvas {
	c := Canvas{
		AxisColor:  Default,
		LabelColor: Default,
		LineColors: []Color{},
		ShowAxis:   true,
		area:       image.Rect(0, 0, width, height),
		points:     make(map[image.Point]Cell),
	}
	return c
}

// Fill adds values to the Canvas
func (c *Canvas) Fill(data [][]float64) {
	if len(data) == 0 {
		return
	}
	c.points = make(map[image.Point]Cell)
	c.graphHeight = c.area.Dy()
	minDataPoint, maxDataPoint := getMinMaxFloat64From2dSlice(data)
	diff := maxDataPoint
	if minDataPoint < 0 {
		diff -= minDataPoint
	}

	// y axis
	if c.ShowAxis {
		c.graphHeight--
		lenMaxDataPoint := len(fmt.Sprintf("%.2f", maxDataPoint))
		lenMinDataPoint := len(fmt.Sprintf("%.2f", minDataPoint))
		if lenMinDataPoint > lenMaxDataPoint {
			lenMaxDataPoint = lenMinDataPoint
		}
		c.horizontalOffset = lenMaxDataPoint + 2 // y-axis plus space before it
		if len(c.HorizontalLabels) != 0 && len(c.HorizontalLabels) <= c.area.Dx()-c.horizontalOffset {
			c.graphHeight--
		}
		verticalScale := diff / float64(c.graphHeight-1)
		cur := minDataPoint
		for i := c.graphHeight - 1; i >= 0; i-- {
			val := fmt.Sprintf("%.2f", cur)
			c.setText(i, lenMaxDataPoint-len(val), val, c.LabelColor)
			c.setRunes(i, lenMaxDataPoint+1, c.AxisColor, YAXIS)
			cur += verticalScale
		}
	}
	c.plotWidth = (c.area.Dx() - c.horizontalOffset) * 2

	// plot the data
	c.horizontalScale = 1.0
	if c.NumDataPoints > 0 {
		c.horizontalScale = math.Round(float64(c.plotWidth/c.NumDataPoints) + 0.5)
	}
	for i, line := range data {
		if len(line) == 0 {
			continue
		} else if c.NumDataPoints > 0 && len(line) > c.NumDataPoints {
			start := len(line) - c.NumDataPoints
			line = line[start:]
		} else if len(line) > c.plotWidth {
			start := len(line) - c.plotWidth
			line = line[start:]
		}

		// y coordinates are calculated as percentages of the graph height. the percentage
		// is the current point minus the smallest point in the dataset over the diff
		// between largest and smallest points. this means that with small graphs there
		// can be some squashing of the graph due to rounding.
		previousHeight := int(((line[0] - minDataPoint) / diff) * float64(c.graphHeight-1))
		for j, val := range line[1:] {
			height := int(((val - minDataPoint) / diff) * float64(c.graphHeight-1))
			c.setLine(
				image.Pt(
					(c.horizontalOffset*2)+int(float64(j)*c.horizontalScale),
					(c.graphHeight-1-previousHeight)*4,
				),
				image.Pt(
					(c.horizontalOffset*2)+int(float64(j+1)*c.horizontalScale),
					(c.graphHeight-1-height)*4,
				),
				c.lineColor(i),
			)
			previousHeight = height
		}
	}

	// x axis
	if c.ShowAxis {
		axisRunes := []rune{ORIGIN}
		remaining := c.plotWidth / 2
		if len(c.HorizontalLabels) != 0 && len(c.HorizontalLabels) <= remaining {
			start := c.HorizontalLabels[0]
			end := c.HorizontalLabels[len(c.HorizontalLabels)-1]
			c.setRunes(c.graphHeight+1, c.horizontalOffset, c.AxisColor, LABELSTART)
			c.setText(c.graphHeight+1, c.horizontalOffset+1, start, c.LabelColor)
			axisRunes = append(axisRunes, XLABELMARKER)
			remaining--
			minWidth := len(start) + len(end) + 4
			if c.maxX >= minWidth {
				labelPos := c.horizontalOffset + 1 + len(start) + c.maxX - minWidth + 2
				c.setText(c.graphHeight+1, labelPos, end, c.LabelColor)
				c.setRunes(c.graphHeight+1, labelPos+len(end), c.AxisColor, LABELEND)
				axisRunes = append(axisRunes, repeatRune(XAXIS, c.maxX-2)...)
				axisRunes = append(axisRunes, XLABELMARKER)
				remaining -= c.maxX + 1
			}
		}
		axisRunes = append(axisRunes, repeatRune(XAXIS, remaining)...)
		c.setRunes(c.graphHeight, c.horizontalOffset-1, c.AxisColor, axisRunes...)
	}
}

// Plot sets the Canvas and return the string representation of it
func (c *Canvas) Plot(data [][]float64) string {
	if len(data) == 0 {
		return ""
	}
	c.Fill(data)
	return c.String()
}

// String allows the Canvas to implement the Stringer interface
func (c Canvas) String() string {
	var b strings.Builder

	// go through each row of the canvas and print the lines
	for row := 0; row < c.area.Dy(); row++ {
		for col := 0; col < c.area.Dx(); col++ {
			b.WriteString(c.points[image.Pt(col, row)].String())
		}
		if row < c.area.Dy()-1 {
			b.WriteString("\n")
		}
	}

	return b.String()
}

func (c *Canvas) setText(row, col int, text string, color Color) {
	for i, letter := range text {
		if p := image.Pt(col+i, row); p.In(c.area) {
			c.points[p] = Cell{
				Rune:   letter,
				color:  color,
				letter: true,
			}
		}
	}
}

func (c *Canvas) setRunes(row, col int, color Color, runes ...rune) {
	for i, r := range runes {
		if p := image.Pt(col+i, row); p.In(c.area) {
			c.points[p] = Cell{
				Rune:  r,
				color: color,
				axis:  true,
			}
		}
	}
}

func (c *Canvas) setPoint(p image.Point, color Color) {
	point := image.Pt(p.X/2, p.Y/4)
	if !point.In(c.area) {
		return
	}
	if x := point.X - c.horizontalOffset + 1; x > c.maxX {
		c.maxX = x
	}
	c.points[point] = Cell{
		Rune:  c.points[point].Rune | BRAILLE[p.Y%4][p.X%2],
		color: color,
	}
}

func (c *Canvas) setLine(p0, p1 image.Point, color Color) {
	for _, p := range line(p0, p1) {
		c.setPoint(p, color)
	}
}

func (c Canvas) lineColor(i int) Color {
	if len(c.LineColors) == 0 || i > len(c.LineColors)-1 {
		return Default
	}
	return c.LineColors[i]
}
