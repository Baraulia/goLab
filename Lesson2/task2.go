package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

var ColorsRGB = map[int][]int{
	1: {255, 0, 0},
	2: {255, 165, 0},
	3: {255, 255, 0},
	4: {0, 128, 0},
	5: {0, 0, 255},
	6: {75, 0, 130},
	7: {238, 130, 238},
	0: {255, 255, 255},
}

type Field struct {
	Size   [][]int
	Points []Point
}

type Point struct {
	X     int
	Y     int
	R     int
	Color []int
}

type Parallelogram struct {
	Field  Field
	A      Point
	B      Point
	C      Point
	D      Point
	Area   float64
	Center Point
}

type Circle struct {
	Field  Field
	R      float64
	Color  []int
	Area   float64
	Center Point
}

func CreateField(height, weight int) (Field, error) {
	var e error
	var NewField = Field{}
	if height <= 0 || weight <= 0 {
		e = fmt.Errorf("введите положительные значения размера поля")
	} else {
		NewField = Field{
			Size: make([][]int, height),
		}
		for i := 0; i <= (height - 1); i++ {
			NewField.Size[i] = make([]int, weight)
		}
	}
	return NewField, e
}

func (p *Point) SetColor() {
	m := (p.X + p.Y) / 2
	for m > 7 {
		m = m / 2
	}
	p.Color = ColorsRGB[m]
}

func (f *Field) ThreePoint() *Field {
	for i := 0; i < 3; i++ {
		f.Points = append(f.Points, Point{
			X: rand.Intn(len(f.Size[0])),
			Y: rand.Intn(len(f.Size)),
			R: 1,
		})
	}
	f.Points[0].SetColor()
	for _, p := range f.Points {
		p.Color = f.Points[0].Color
		f.Size[p.Y][p.X] = 1
	}
	return f
}
func (p *Parallelogram) SetFourthPoint() error {
	var err error
	if len(p.Field.Points) != 3 {
		p.Field.ThreePoint()
	} else {
		p.A = p.Field.Points[0]
		p.B = p.Field.Points[1]
		p.C = p.Field.Points[2]
	}
	switch x := (p.A.X + p.B.X) - p.C.X; x <= len(p.Field.Size[0])-1 && x >= 0 {
	case true:
		switch y := (p.A.Y + p.B.Y) - p.C.Y; y <= len(p.Field.Size)-1 && y >= 0 {
		case true:
			p.D.X = x
			p.D.Y = y
			p.Field.Size[y][x] = 1
		case false:
			break
		}
	case false:
		switch x := (p.B.X + p.C.X) - p.A.X; x <= len(p.Field.Size[0])-1 && x >= 0 {
		case true:
			switch y := (p.B.Y + p.C.Y) - p.A.Y; y <= len(p.Field.Size)-1 && y >= 0 {
			case true:
				p.D.X = x
				p.D.Y = y
				p.Field.Size[y][x] = 1
			case false:
				break
			}
		case false:
			switch x := (p.A.X + p.C.X) - p.B.X; x <= len(p.Field.Size[0])-1 && x >= 0 {
			case true:
				switch y := (p.A.Y + p.C.Y) - p.B.Y; y <= len(p.Field.Size)-1 && y >= 0 {
				case true:
					p.D.X = x
					p.D.Y = y
					p.Field.Size[y][x] = 1
				case false:
					err = fmt.Errorf("при заданных координатах трех точек и размера поля невозможно построить четвертую точку параллелограмма")
				}
			}
		}
	}

	return err
}

func (p *Parallelogram) CenterCalc() Point {
	points := []Point{p.A, p.B, p.C, p.D}
	pointsX := points
	sort.Slice(pointsX, func(i, j int) bool {
		return pointsX[i].X < pointsX[j].X
	})
	fmt.Println(pointsX)
	pointsY := points
	sort.Slice(pointsY, func(i, j int) bool {
		return pointsY[i].Y < pointsY[j].Y
	})
	fmt.Println(pointsY)
	p.Center.X = (pointsX[0].X + pointsX[3].X) / 2
	p.Center.Y = (pointsY[0].Y + pointsY[3].Y) / 2
	p.Field.Size[p.Center.Y][p.Center.X] = 2
	return p.Center
}

func (p *Parallelogram) GetArea() float64 {
	Area := math.Abs(float64(p.A.X*(p.B.Y-p.C.Y) + p.B.X*(p.C.Y-p.A.Y) + p.C.X*(p.A.Y-p.B.Y)))
	p.Area = Area
	return Area
}

func (c *Circle) GetArea() float64 {
	c.Area = math.Pi * c.R * c.R
	return c.Area
}

func BuildCircle(p *Parallelogram) (c Circle) {
	if p.Field.Points == nil {
		p.Field.ThreePoint()
	}
	if p.D.X == 0 {
		p.SetFourthPoint()
	}
	c.Field = p.Field
	c.Area = p.GetArea()
	c.Center = p.CenterCalc()
	c.R = math.Sqrt(c.Area / math.Pi)
	return c
}

func (c *Circle) SetColor() {
	if len(c.Color) == 0 {
		s := len(c.Field.Size[0]) * len(c.Field.Size)
		color := c.GetArea() / float64(s)
		for color > 7 {
			color = color / 2
		}
		c.Color = ColorsRGB[int(color)]
	}
}

func main() {
	var v, _ = CreateField(10, 10)
	v.ThreePoint()
	var p = Parallelogram{
		Field: v,
	}
	err := p.SetFourthPoint()
	if err != nil {
		fmt.Println(err)
	}
	p.CenterCalc()
	fmt.Println(p.Center)
	for i := range v.Size {
		fmt.Println(v.Size[i])
	}
	c := BuildCircle(&p)
	c.SetColor()
}
