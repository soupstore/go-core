package maths

type Rectangle struct {
	X      int16
	Y      int16
	Width  int16
	Height int16
}

func NewRectangleFromCenter(center Point, radius int16) Rectangle {
	return Rectangle{
		X:      center.X - radius,
		Y:      center.Y - radius,
		Width:  2*radius + 1,
		Height: 2*radius + 1,
	}
}

// Intersects is your classic AABB algorithm
func (r Rectangle) Intersects(r2 Rectangle) bool {
	return r.X < r2.X+r2.Width &&
		r.X+r.Width > r2.X &&
		r.Y < r2.Y+r2.Height &&
		r.Y+r.Height > r2.Y
}

func (r Rectangle) Contains(p Point) bool {
	return r.X <= p.X && p.X < r.X+r.Width &&
		r.Y <= p.Y && p.Y < r.Y+r.Height
}
