package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestCreateField(t *testing.T) {
	Convey("Negative size of field", t, func() {
		_, err := CreateField(-5, 10)
		So(err, ShouldBeError, "введите положительные значения размера поля")
	})
	Convey("Zero size of field", t, func() {
		_, err := CreateField(0, 10)
		So(err, ShouldBeError, "введите положительные значения размера поля")
	})
	Convey("Test of creating field", t, func() {
		f, _ := CreateField(2, 2)
		So(f, ShouldResemble, Field{Size: [][]int{{0, 0}, {0, 0}}, Points: nil})
	})
}

func TestPoint_SetColor(t *testing.T) {
	Convey("Big coordinates of the point", t, func() {
		p := Point{X: 100, Y: 100}
		p.SetColor()
		bools := false
		for _, i := range ColorsRGB {
			if p.Color[0] == i[0] && p.Color[1] == i[1] && p.Color[2] == i[2] {
				bools = true
				break
			}
		}
		So(bools, ShouldEqual, true)
	})
}

func TestField_ThreePoint(t *testing.T) {
	Convey("Number of created points", t, func() {
		f, _ := CreateField(10, 10)
		f.ThreePoint()
		So(len(f.Points), ShouldEqual, 3)

	})
}

func TestParallelogram_SetFourthPoint(t *testing.T) {
	Convey("Test fourth point", t, func() {
		f, _ := CreateField(10, 10)
		f1, _ := CreateField(1000, 1000)
		f2, _ := CreateField(10, 10)
		f.ThreePoint()
		p := Parallelogram{Field: f}
		p1 := Parallelogram{Field: f1}
		p2 := Parallelogram{Field: f2}
		p2.A, p2.B, p2.C = Point{0, 0, 1, ColorsRGB[0]}, Point{2, 4, 1, ColorsRGB[0]}, Point{7, 1, 1, ColorsRGB[0]}
		err := p.SetFourthPoint()
		if err != nil {
			So(err, ShouldBeError, "при заданных координатах трех точек и размера поля невозможно построить четвертую точку параллелограмма")
		}
		So(p.D.X >= 0 && p.D.X <= len(f.Size[0])-1 && p.D.Y >= 0 && p.D.Y <= len(f.Size)-1, ShouldEqual, true)

		err = p1.SetFourthPoint()
		if err != nil {
			So(err, ShouldBeError, "при заданных координатах трех точек и размера поля невозможно построить четвертую точку параллелограмма")
		}
		So(p.D.X >= 0 && p.D.X <= len(f.Size[0])-1 && p.D.Y >= 0 && p.D.Y <= len(f.Size)-1, ShouldEqual, true)

		err = p2.SetFourthPoint()
		if err != nil {
			So(err, ShouldBeError, "при заданных координатах трех точек и размера поля невозможно построить четвертую точку параллелограмма")
		}
		So(p.D.X >= 0 && p.D.X <= len(f.Size[0])-1 && p.D.Y >= 0 && p.D.Y <= len(f.Size)-1, ShouldEqual, true)

		p2.Field.Points = []Point{{X: 2, Y: 5, R: 1, Color: ColorsRGB[0]}, {X: 6, Y: 8, R: 1, Color: ColorsRGB[0]}, {X: 9, Y: 6, R: 1, Color: ColorsRGB[0]}}
		err = p2.SetFourthPoint()
		if err != nil {
			So(err, ShouldBeError, "при заданных координатах трех точек и размера поля невозможно построить четвертую точку параллелограмма")
		}
		So(p.D.X >= 0 && p.D.X <= len(f.Size[0])-1 && p.D.Y >= 0 && p.D.Y <= len(f.Size)-1, ShouldEqual, true)

		p2.Field.Points = []Point{{0, 9, 1, ColorsRGB[0]}, {5, 0, 1, ColorsRGB[0]}, {9, 9, 1, ColorsRGB[0]}}
		err = p2.SetFourthPoint()
		if err != nil {
			So(err, ShouldBeError, "при заданных координатах трех точек и размера поля невозможно построить четвертую точку параллелограмма")
		}
		So(p.D.X >= 0 && p.D.X <= len(f.Size[0])-1 && p.D.Y >= 0 && p.D.Y <= len(f.Size)-1, ShouldEqual, true)
	})
}

func TestParallelogram_CenterCalc(t *testing.T) {
	Convey("Check Parallelogram_CenterCalc", t, func() {
		f, _ := CreateField(9, 9)
		p := Parallelogram{Field: f}
		p.A, p.B, p.C, p.D = Point{0, 0, 1, ColorsRGB[0]},
			Point{len(p.Field.Size[0]) - 1, len(p.Field.Size) - 1, 1, ColorsRGB[0]},
			Point{len(p.Field.Size[0]) - 1, 0, 1, ColorsRGB[0]},
			Point{0, len(p.Field.Size) - 1, 1, ColorsRGB[0]}
		p.CenterCalc()
		So(p.Center.X, ShouldEqual, 4)
		So(p.Center.Y, ShouldEqual, 4)
	})
}

func TestParallelogram_GetArea(t *testing.T) {
	Convey("Test GetArea for parallelogram", t, func() {
		f, _ := CreateField(10, 10)
		p := Parallelogram{Field: f}
		p.A, p.B, p.C = Point{2, 2, 1, ColorsRGB[0]},
			Point{5, 0, 1, ColorsRGB[0]},
			Point{8, 5, 1, ColorsRGB[0]}
		So(p.GetArea(), ShouldEqual, 21)
	})
}

func TestCircle_GetArea(t *testing.T) {
	Convey("Test GetArea for Circle", t, func() {
		f, _ := CreateField(10, 10)
		c := Circle{
			Field: f,
			R:     3,
			Color: ColorsRGB[2],
		}
		So(c.GetArea(), ShouldAlmostEqual, 28.27, 0.1)
	})
}

func TestBuildCircle(t *testing.T) {
	Convey("Test Build Circle", t, func() {
		f, _ := CreateField(10, 10)
		p := Parallelogram{Field: f}
		p2 := Parallelogram{Field: f}
		p.Field.Points = []Point{{2, 5, 1, ColorsRGB[0]}, {6, 8, 1, ColorsRGB[0]}, {9, 6, 1, ColorsRGB[0]}}
		c := BuildCircle(&p)
		c2 := BuildCircle(&p2)
		So(p.GetArea(), ShouldAlmostEqual, c.GetArea(), 0.1)
		So(p2.GetArea(), ShouldAlmostEqual, c2.GetArea(), 0.1)
	})
}

func TestCircle_SetColor(t *testing.T) {
	Convey("Test SetColor for Circle", t, func() {
		f, _ := CreateField(10, 10)
		c := Circle{Field: f}
		c.R = 3.0
		c.GetArea()
		c.SetColor()
		bools := false
		for _, i := range ColorsRGB {
			if c.Color[0] == i[0] && c.Color[1] == i[1] && c.Color[2] == i[2] {
				bools = true
				break
			}
		}
		So(bools, ShouldEqual, true)
	})
}
