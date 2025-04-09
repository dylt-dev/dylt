package color

import (
	"fmt"
	"strings"
)

var Ansi = struct {
	Reset string
}{
	Reset: "\033[0m",
}

type color interface {
	AnsiBg() string
	AnsiFg() string
}

type colorbase struct {
	R byte `json:"r"`
	G byte `json:"g"`
	B byte `json:"b"`
}

type color_sys struct {
	colorbase
	Fg   int
	Bg   int
	Name string `json:"name"`
}

func (c color_sys) AnsiBg() string {
	return fmt.Sprintf("\033[%dm", c.Bg)
}

func (c color_sys) AnsiFg() string {
	return fmt.Sprintf("\033[%dm", c.Fg)
}

type color_ansi256 struct {
	colorbase
	Val int
}

func (c color_ansi256) AnsiBg() string {
	return fmt.Sprintf("\033[48;5;%dm", c.Val)
}

func (c color_ansi256) AnsiFg() string {
	return fmt.Sprintf("\033[38;5;%dm", c.Val)
}

type color_x11 struct {
	colorbase
	Name string `json:"name"`
}

func (c color_x11) AnsiBg() string {
	return fmt.Sprintf("\033[48;2;%d;%d;%dm", c.R, c.G, c.B)
}

func (c color_x11) AnsiFg() string {
	return fmt.Sprintf("\033[38;2;%d;%d;%dm", c.R, c.G, c.B)
}

type Styledstring string

// type styledstring struct {
// 	string
// }

func (s Styledstring) Bg(c color) Styledstring {
	return s.Style(c.AnsiBg())
}

func (s Styledstring) Fg(c color) Styledstring {
	return s.Style(c.AnsiFg())
}

func (s Styledstring) FgBg(fgc, bgc color) Styledstring {
	return s.Fg(fgc).Bg(bgc)
}

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
