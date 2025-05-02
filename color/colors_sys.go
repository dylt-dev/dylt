package color

var Sys = struct {
	Black Color_sys
	Red Color_sys
	Green Color_sys
	Yellow Color_sys
	Blue Color_sys
	Magenta Color_sys
	Cyan Color_sys
	White Color_sys
	Default Color_sys
	BrightBlack Color_sys
	BrightRed Color_sys
	BrightGreen Color_sys
	BrightYellow Color_sys
	BrightBlue Color_sys
	BrightMagenta Color_sys
	BrightCyan Color_sys
	BrightWhite Color_sys
}{
	Black: Color_sys{Name: "Black", Fg: 30, Bg: 40, colorbase: colorbase{R: 0, G: 0, B: 0} },
	Red: Color_sys{Name: "Red", Fg: 31, Bg: 41, colorbase: colorbase{R: 128, G: 0, B: 0} },
	Green: Color_sys{Name: "Green", Fg: 32, Bg: 42, colorbase: colorbase{R: 0, G: 128, B: 0} },
	Yellow: Color_sys{Name: "Yellow", Fg: 33, Bg: 43, colorbase: colorbase{R: 128, G: 128, B: 0} },
	Blue: Color_sys{Name: "Blue", Fg: 34, Bg: 44, colorbase: colorbase{R: 0, G: 0, B: 128} },
	Magenta: Color_sys{Name: "Magenta", Fg: 35, Bg: 45, colorbase: colorbase{R: 128, G: 0, B: 128} },
	Cyan: Color_sys{Name: "Cyan", Fg: 36, Bg: 46, colorbase: colorbase{R: 0, G: 128, B: 128} },
	White: Color_sys{Name: "White", Fg: 37, Bg: 47, colorbase: colorbase{R: 192, G: 192, B: 192} },
	Default: Color_sys{Name: "Default", Fg: 39, Bg: 49, colorbase: colorbase{R: 128, G: 128, B: 128} },
	BrightBlack: Color_sys{Name: "BrightBlack", Fg: 90, Bg: 100, colorbase: colorbase{R: 128, G: 128, B: 128} },
	BrightRed: Color_sys{Name: "BrightRed", Fg: 91, Bg: 101, colorbase: colorbase{R: 255, G: 0, B: 0} },
	BrightGreen: Color_sys{Name: "BrightGreen", Fg: 92, Bg: 102, colorbase: colorbase{R: 0, G: 255, B: 0} },
	BrightYellow: Color_sys{Name: "BrightYellow", Fg: 93, Bg: 103, colorbase: colorbase{R: 255, G: 255, B: 0} },
	BrightBlue: Color_sys{Name: "BrightBlue", Fg: 94, Bg: 104, colorbase: colorbase{R: 0, G: 0, B: 255} },
	BrightMagenta: Color_sys{Name: "BrightMagenta", Fg: 95, Bg: 105, colorbase: colorbase{R: 255, G: 0, B: 255} },
	BrightCyan: Color_sys{Name: "BrightCyan", Fg: 96, Bg: 106, colorbase: colorbase{R: 0, G: 255, B: 255} },
	BrightWhite: Color_sys{Name: "BrightWhite", Fg: 97, Bg: 107, colorbase: colorbase{R: 255, G: 255, B: 255} },
}
