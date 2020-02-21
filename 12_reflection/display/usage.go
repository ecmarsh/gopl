package display

type Movie struct {
	Title, Subtitle string
	Year            int
	Color           bool
	Actor           map[string]string
	Oscars          []string
	Sequel          *string
}

func example() {
	strangelove := Movie{
		Title:    "Dr. Strangelove",
		Subtitle: "How I Learned to Stop Worrying and Love the Bomb",
		Year:     1964,
		Color:    false,
		Actor: map[string]string{
			"Dr.Strangelove":             "Peter Sellers",
			"Grp. Capt. Lionel Mandrake": "Peter Sellers",
			"Pres. Merkin Muffley":       "Peter Sellers",
			"Gen. Buck Turgidson":        "George C. Scott",
			"Brig. Gen. Jack D. Ripper":  "Sterling Hayden",
			`Maj. T.J. "King" Kong`:      "Slim Pickens",
		},
		Oscars: []string{
			"Best Actor (Nomin.)",
			"Best Adapated Screenplay (Nomin.)",
			"Best Director(Nomin.)",
			"Best Picture (Nomin.)",
		},
	}

	// Example Usage
	// Display strangelove (display.Movie)
	/*
		   // Partial output:
		   strangelove.Title = "Dr.Strangelove"
		   strangelove.Year = 1964
		   strangelove.Color = false
		   strangelove.Actor["Grp. Capt. Lionel Mandrake"] = "Peter Sellers"
		   strangelove.Oscars[0] = "Best Actor (nomin.)"
		   strangelove.Sequel = nil
			 ...
	*/

	// can also be used to display internals of library types such as *os.File
	// Display ("os.Stderr", os.Stderr)
	// Output:
	// (*(*os.Stderr).file).fd = 2)
	// (*(*os.Stderr).file).name = "/dev/stderr")
	// (*(*os.Stderr).file).nepipe = 0)

	// Example showing differences of using pointers vs concrete types to Display
	var i interface{} = 3

	// Ex 1: reflect.VAlueOf always returns a concrete type since
	// it extracts contents of interface value.
	Display("i", i)
	// Display i (int):
	i = 3

	// Ex 2: returns kind Ptr which calls `Elem` on value, which returns value
	// representing the _variable_ i itself of kind `Interface.
	Display("&i", &i)
	// Display &i (*interface {}):
	// (*&i).type = int
	// (*&i).value = 3

	/* Example of issues with cycles as currently implemented */
	// a struct that points to itself
	type Cycle struct {
		Value int
		Tail  *Cycle
	}
	var c Cycle
	c = Cycle{42, &c}
	Display("c", c)
	// Would print the ever growing expansion
	// Display c (display.Cycle):
	// c.Value = 42
	// (*c.Tail).Value = 42
	// (*(*c.Tail).Tail).Value = 42
	// (*(*(*c.Tail).Tail).Tail).Value = 42
	// ... ad infinitum ...
}
