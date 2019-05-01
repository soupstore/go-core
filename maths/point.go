package maths

import "math"

// Point is a position inside a chunk or the position of a chunk itself.
type Point struct {
	X int16
	Y int16
}

func (c Point) Add(c2 Point) Point {
	return Point{c.X + c2.X, c.Y + c2.Y}
}

func (c Point) Sub(c2 Point) Point {
	return Point{c.X - c2.X, c.Y - c2.Y}
}

func (c Point) Abs() Point {
	if c.X < 0 {
		c.X = -c.X
	}
	if c.Y < 0 {
		c.Y = -c.Y
	}

	return c
}

func (c Point) Euclidean() float64 {
	return math.Sqrt(math.Pow(float64(c.X), 2) + math.Pow(float64(c.Y), 2))
}
