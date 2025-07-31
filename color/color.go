package color

import (
	"fmt"
	"strconv"
	"strings"
)

var Ansi = struct {
	Reset string
}{
	Reset: "\033[0m",
}

// There are multiple ways of representing an ANSI color
// All can be expressed as an ANSI sequence for foreground
// color, and for background color
// This interface expresses that commonality
// @note maybe `ansicolor` would be a better name
type color interface {
	AnsiBg() string
	AnsiFg() string
}

// Some colors have names, some don't, but all have RGB values
// colorbase expreses that, and has JSON field mappings are well
// @hote I'm not sure why json mappings to lowercase letters are
// so important
type colorbase struct {
	R byte `json:"r"`
	G byte `json:"g"`
	B byte `json:"b"`
}

// Eg \033[48;2;240;248;255m (AliceBlue - background)
func (c colorbase) AnsiBg() string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", c.R, c.G, c.B)
}

// Eg \033[38;2;240;248;255m (AliceBlue - background)
func (c colorbase) AnsiFg() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

// ANSI sys colors: 16(ish), very basic
// @note why a JSON mapping for a lowercase 'name'
// but not Fg or Bg
type Color_sys struct {
	// embedded RGB
	colorbase
	// ANSI code for foreground color
	Fg int
	// ANSI code for background color
	Bg int
	// Official ANSI name
	Name string `json:"name"`
}

// ANSI code: eg \033[41m (red background)
func (c Color_sys) AnsiBg() string {
	return fmt.Sprintf("\033[%dm", c.Bg)
}

// ANSI code: eg \033[31m (red foreground)
func (c Color_sys) AnsiFg() string {
	return fmt.Sprintf("\033[%dm", c.Fg)
}

// ANSI 256 colors. Similar to system colors, but there are 256 of them.
// They each have an RGB value but the ANSI code is based on the number
// of the color (0-255)
// @note some attempts have been made at naming these colors but nothing feels
// standard or widely adopted
type color_ansi256 struct {
	// embedded RGB
	colorbase
	// Simple 0-255 ordinal for the color
	Val int
}

// Eg \033[48;5;13m
func (c color_ansi256) AnsiBg() string {
	return fmt.Sprintf("\033[48;5;%dm", c.Val)
}

// Eg \033[38;5;13m
func (c color_ansi256) AnsiFg() string {
	return fmt.Sprintf("\033[38;5;%dm", c.Val)
}

// X11 colors. Somewhat similar to Web colors, but they are different and there
// more of them (600+)
// X11 colors have names and RGB values. ANSI does not have
// AnsiBg() and AnsiFg() methods are not necessary. The colorbase methods
// are sufficient.
type color_x11 struct {
	colorbase
	Name string `json:"name"`
}

// Just a string, which can
type Styledstring string

// Convenience method for adding an ANSI sequence
// for a background color
func (s Styledstring) Bg(c color) Styledstring {
	return s.Style(c.AnsiBg())
}

// Convenience method for adding an ANSI sequence
// for a foreground color
func (s Styledstring) Fg(c color) Styledstring {
	return s.Style(c.AnsiFg())
}

// Convenience method for adding an ANSI sequence
// for both a foreground color and a background color
func (s Styledstring) FgBg(fgc, bgc color) Styledstring {
	return s.Fg(fgc).Bg(bgc)
}

// 'Style' a string with ANSI sequences, eg Foreground &
// background color, by prepending an ANSI sequence to the string
// Append an ANSI Reset sequence to the end of the string
//
// strings.Builder is used for a bit of effiency, vs
// fmt.Sprintf() calls, which are simpler but might cause more
// allocations. Styledstring was originally created for logging so
// efficiency is a priority
func (s Styledstring) Style(style string) Styledstring {
	sb := strings.Builder{}
	sb.Grow(len(string(s)) * 2)

	sb.WriteString(style)
	sb.WriteString(string(s))

	if !strings.HasSuffix(string(s), Ansi.Reset) {
		sb.WriteString(Ansi.Reset)
	}
	s = Styledstring(sb.String())

	return s
}

func StyleBool (flag bool) Styledstring {
	s := Styledstring(strconv.FormatBool(flag))
	if flag {
		s = s.Fg(X11.Green)
	} else {
		s = s.Fg(X11.Red)
	}

	return s
}