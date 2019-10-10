// Package embed shows two correct forms of embedding with struct literals.
package embed

import "fmt"

// Point defines X, Y coordinates
type Point struct {
	X, Y int
}

// Circle embeds a point anonymously
type Circle struct {
	Point
	Radius int
}

// Wheel embeds a circle anonymously
type Wheel struct {
	Circle
	Spokes int
}

func embed() {
	var w Wheel
	// Equivalent assignments
	w = Wheel{Circle{Point{8, 8}, 5}, 20}
	w = Wheel{
		Circle: Circle{
			Point:  Point{X: 8, Y: 8},
			Radius: 5,
		},
		Spokes: 20, // Note necessity of trailing commas here and after Radius
	}

	fmt.Printf("%#v\n", w)
	// Output:
	// Wheel{Circle:Circle{Point:Point{X:8, Y:8}, Radius:5}, Spokes:20}
	// Note: "#" adverb causes %v verb to display form in Go syntax

	w.X = 8
	fmt.Printf("%#v\n", w)
	// Output:
	// Wheel{Circle:Circle{Point:Point{X:8, Y:8}, Radius:5}, Spokes:20}
}

// Note regardless of whether Point and Circle are exported (vs point, circle),
// we could still use the shorthand form: w.X = 8 outside of package
// but NOT explicit from (w.circle.point.X = 8) since fields inaccessable.
