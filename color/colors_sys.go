package color

var Sys = struct {
	Black         color_sys
	Red           color_sys
	Green         color_sys
	Yellow        color_sys
	Blue          color_sys
	Magenta       color_sys
	Cyan          color_sys
	White         color_sys
	Default       color_sys
	BrightBlack   color_sys
	BrightRed     color_sys
	BrightGreen   color_sys
	BrightYellow  color_sys
	BrightBlue    color_sys
	BrightMagenta color_sys
	BrightCyan    color_sys
	BrightWhite   color_sys
}{
	Black:         color_sys{Name: "Black", Fg: 30, Bg: 40, colorbase: colorbase{R: 0, G: 0, B: 0}},
	Red:           color_sys{Name: "Red", Fg: 31, Bg: 41, colorbase: colorbase{R: 128, G: 0, B: 0}},
	Green:         color_sys{Name: "Green", Fg: 32, Bg: 42, colorbase: colorbase{R: 0, G: 128, B: 0}},
	Yellow:        color_sys{Name: "Yellow", Fg: 33, Bg: 43, colorbase: colorbase{R: 128, G: 128, B: 0}},
	Blue:          color_sys{Name: "Blue", Fg: 34, Bg: 44, colorbase: colorbase{R: 0, G: 0, B: 128}},
	Magenta:       color_sys{Name: "Magenta", Fg: 35, Bg: 45, colorbase: colorbase{R: 128, G: 0, B: 128}},
	Cyan:          color_sys{Name: "Cyan", Fg: 36, Bg: 46, colorbase: colorbase{R: 0, G: 128, B: 128}},
	White:         color_sys{Name: "White", Fg: 37, Bg: 47, colorbase: colorbase{R: 192, G: 192, B: 192}},
	Default:       color_sys{Name: "Default", Fg: 39, Bg: 49, colorbase: colorbase{R: 128, G: 128, B: 128}},
	BrightBlack:   color_sys{Name: "BrightBlack", Fg: 90, Bg: 100, colorbase: colorbase{R: 128, G: 128, B: 128}},
	BrightRed:     color_sys{Name: "BrightRed", Fg: 91, Bg: 101, colorbase: colorbase{R: 255, G: 0, B: 0}},
	BrightGreen:   color_sys{Name: "BrightGreen", Fg: 92, Bg: 102, colorbase: colorbase{R: 0, G: 255, B: 0}},
	BrightYellow:  color_sys{Name: "BrightYellow", Fg: 93, Bg: 103, colorbase: colorbase{R: 255, G: 255, B: 0}},
	BrightBlue:    color_sys{Name: "BrightBlue", Fg: 94, Bg: 104, colorbase: colorbase{R: 0, G: 0, B: 255}},
	BrightMagenta: color_sys{Name: "BrightMagenta", Fg: 95, Bg: 105, colorbase: colorbase{R: 255, G: 0, B: 255}},
	BrightCyan:    color_sys{Name: "BrightCyan", Fg: 96, Bg: 106, colorbase: colorbase{R: 0, G: 255, B: 255}},
	BrightWhite:   color_sys{Name: "BrightWhite", Fg: 97, Bg: 107, colorbase: colorbase{R: 255, G: 255, B: 255}},
}
