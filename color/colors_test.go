package color

import (
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"os"
	"regexp"
	"slices"
	"strings"
	"testing"
	"text/template"
	"unicode"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

type capColorStruct struct {
	CapColor color_x11
	OriginalName string
}

func TestCapitalizeName (t *testing.T) {
	s0 := "hello"
	s1 := "Hello"
	s2 := "13"
	s3 := ""
	assert.Equal(t, "Hello", createCapitalName(s0))
	assert.Equal(t, "Hello", createCapitalName(s1))
	assert.Equal(t, "13", createCapitalName(s2))
	assert.Equal(t, "", createCapitalName(s3))
}

func TestCleanColorsX11Go (t *testing.T) {
	var tmpl template.Template = *loadAndTestTemplate(t, "./colors_x11.go.clean.tmpl")
	w, err := os.Create("./colors_x11.go")
	require.NoError(t, err)
	tmpl.Execute(w, struct{}{})
}


func TestCreateCapColors (t *testing.T) {
	var colors = readAndTestColorsX11(t)
	var capAliases = createCapitalColorAliases(colors)
	require.NotEmpty(t, capAliases)
	require.LessOrEqual(t, len(capAliases), len(colors))
	for _, capAlias := range capAliases {
		var rx = regexp.MustCompile("^[A-Z].*")
		assert.True(t, rx.Match([]byte(capAlias.CapColor.Name)))
	}
	t.Log(capAliases)
}



// A color is a named RGB tri-variant
// colorStruct gets that
// type colorStruct struct {
// 	R    byte   `json:"r"`
// 	G    byte   `json:"g"`
// 	B    byte   `json:"b"`
// 	Name string `json:"name"`
// }

// Convert a colorStruct's RGB into a single uint32
// func (c colorStruct) rgb() uint32 {
// 	r32 := uint32(c.R)
// 	g32 := uint32(c.G)
// 	b32 := uint32(c.B)
// 	return uint32(r32<<16 + g32<<8 + b32)
// }

// type _color_16 uint32

// func TestGen16(t *testing.T) {
// 	colors := map[byte]_color_16{}
// 	colors[0] = _color_16(rgb(0, 0, 0))
// 	colors[1] = _color_16(rgb(128, 0, 0))
// 	colors[2] = _color_16(rgb(0, 128, 0))
// 	colors[3] = _color_16(rgb(128, 128, 0))
// 	colors[4] = _color_16(rgb(0, 0, 128))
// 	colors[5] = _color_16(rgb(128, 0, 128))
// 	colors[6] = _color_16(rgb(0, 128, 128))
// 	colors[7] = _color_16(rgb(192, 192, 192))
// 	colors[8] = _color_16(rgb(128, 128, 128))
// 	colors[9] = _color_16(rgb(255, 0, 0))
// 	colors[10] = _color_16(rgb(0, 255, 0))
// 	colors[11] = _color_16(rgb(255, 255, 0))
// 	colors[12] = _color_16(rgb(0, 0, 255))
// 	colors[13] = _color_16(rgb(255, 0, 255))
// 	colors[14] = _color_16(rgb(0, 255, 255))
// 	colors[15] = _color_16(rgb(255, 255, 255))

// 	t.Logf("%#v", colors)
// }

// type _color_256 uint32

// X11 color names have a couple of snags.
// - Some are well-behaved alphabetical CamelCase
// - Some end in digits
// - Some have spaces
//
// Eg
// "antique white"
// "AntiqueWhite"
// "AntiqueWhite1"
// "AntiqueWhite2"
// "AntiqueWhite3"
// "AntiqueWhite4"
//
// In all cases, the 'space version' is the same as the non-space version, eg the RGB
// values for "antique white" and "AntiqueWhite" are the same. We don't want to support both
// so the space version needs to be discarded. That means it's important to confirm there
// are no names-with-spaces without a non-spaced alternative.
//
// Names ending in digits are ok. They are just alternatives for a simpler name, eg AntiqueWhite1
// (or 2, 3, etc) for AntiqueWhite. These are simply treated as valid names and get no special treatment.
//
// If a name does not fit into any categories, it is considered a dirty name. The presence of any dirty
// names fails the test.
func TestFilterX11Names(t *testing.T) {
	var colors []color_x11 = readAndTestColorsX11(t)

	var spaceNames = map[string]color_x11{}
	var cleanNames = map[string]color_x11{}
	var dirtyNames = map[string]color_x11{}

	// Step 1 - sort all names into 3 catgeories: names with spaces, names ending in digits, and 'clean' names
	for _, color := range colors {
		if isHasSpace(color.Name) {
			spaceNames[color.Name] = color
		} else if isClean(color.Name) {
			cleanNames[color.Name] = color
		} else {
			dirtyNames[color.Name] = color
		}
	}
	assert.Zero(t, len(dirtyNames))

	t.Log("=== isSpace ===")
	t.Log(strings.Join(slices.Sorted(maps.Keys(spaceNames)), "\n"))
	t.Log("=== isClean ===")
	t.Log(strings.Join(slices.Sorted(maps.Keys(cleanNames)), "\n"))

	// Step 2 - for each space name, normalize the name and see if it exists in the map. If not, add it.
	for name, color := range spaceNames {
		normalName := normalize(name)
		if _, is := cleanNames[normalName]; !is {
			cleanNames[normalName] = color
		}
	}

	encoder := json.NewEncoder(os.Stdout)
	encoder.SetIndent("", "\t")
	encoder.Encode(cleanNames)
}

func TestGenColorsAnsi256Go(t *testing.T) {
	tmpl := template.New("go")
	bufTmpl, err := os.ReadFile("./colors_ansi256.go.tmpl")
	require.NoError(t, err)
	_, err = tmpl.Parse(string(bufTmpl))
	require.NoError(t, err)

	w, err := os.Create("./colors_ansi256.go")
	require.NoError(t, err)
	var colors map[int]color_ansi256 = readAndTestColorsAnsi256(t)
	err = tmpl.Execute(w, colors)
	require.NoError(t, err)
}

func TestGenColorsSysGo(t *testing.T) {
	tmpl := template.New("go")
	bufTmpl, err := os.ReadFile("./colors_sys.go.tmpl")
	require.NoError(t, err)
	_, err = tmpl.Parse(string(bufTmpl))
	require.NoError(t, err)

	w, err := os.Create("./colors_sys.go")
	require.NoError(t, err)
	var colors []Color_sys = readAndTestColorsSys(t)
	err = tmpl.Execute(w, colors)
	require.NoError(t, err)
}

func TestGenColorsX11Go(t *testing.T) {
	w, err := os.Create("./colors_x11.go")
	require.NoError(t, err)
	genColorsX11Go(t, w)
}

func TestGenColorsX11GoStdout(t *testing.T) {
	genColorsX11Go(t, os.Stdout)
}

func TestGenAnsi256Json(t *testing.T) {
	colors := map[int]color_ansi256{}

	colors[0] = color_ansi256{colorbase: colorbase{R: 0, G: 0, B: 0}, Val: 0}
	colors[1] = color_ansi256{colorbase: colorbase{R: 128, G: 0, B: 0}, Val: 1}
	colors[2] = color_ansi256{colorbase: colorbase{R: 0, G: 128, B: 0}, Val: 2}
	colors[3] = color_ansi256{colorbase: colorbase{R: 128, G: 128, B: 0}, Val: 3}
	colors[4] = color_ansi256{colorbase: colorbase{R: 0, G: 0, B: 128}, Val: 4}
	colors[5] = color_ansi256{colorbase: colorbase{R: 128, G: 0, B: 128}, Val: 5}
	colors[6] = color_ansi256{colorbase: colorbase{R: 0, G: 128, B: 128}, Val: 6}
	colors[7] = color_ansi256{colorbase: colorbase{R: 192, G: 192, B: 192}, Val: 7}
	colors[8] = color_ansi256{colorbase: colorbase{R: 128, G: 128, B: 128}, Val: 8}
	colors[9] = color_ansi256{colorbase: colorbase{R: 255, G: 0, B: 0}, Val: 9}
	colors[10] = color_ansi256{colorbase: colorbase{R: 0, G: 255, B: 0}, Val: 10}
	colors[11] = color_ansi256{colorbase: colorbase{R: 255, G: 255, B: 0}, Val: 11}
	colors[12] = color_ansi256{colorbase: colorbase{R: 0, G: 0, B: 255}, Val: 12}
	colors[13] = color_ansi256{colorbase: colorbase{R: 255, G: 0, B: 255}, Val: 13}
	colors[14] = color_ansi256{colorbase: colorbase{R: 0, G: 255, B: 255}, Val: 14}
	colors[15] = color_ansi256{colorbase: colorbase{R: 255, G: 255, B: 255}, Val: 15}

	gradient := []byte{0, 95, 135, 175, 215, 255}
	var r, g, b byte
	var ir, ig, ib int
	var clr color_ansi256
	for i := range 6 * 6 * 6 {
		ir = i / 36
		ig = (i % 36) / 6
		ib = i % 6
		r = gradient[ir]
		g = gradient[ig]
		b = gradient[ib]
		val := i + 16
		clr = color_ansi256{colorbase: colorbase{R: r, G: g, B: b}, Val: val}
		colors[i+16] = clr
	}

	grayGradient := []byte{8, 18, 28, 38, 48, 58, 68, 78, 88, 98, 108, 118, 128, 138, 148, 158, 168, 178, 188, 198, 208, 218, 228, 238}
	for i, n := range grayGradient {
		val := i + 232
		clr = color_ansi256{colorbase: colorbase{R: n, G: n, B: n}, Val: val}
		colors[val] = clr
	}

	w, err := os.Create("colors_256.json")
	require.NoError(t, err)
	encoder := json.NewEncoder(w)
	encoder.SetIndent("", "\t")
	err = encoder.Encode(colors)
	require.NoError(t, err)
	t.Logf("%#v", colors)
}

func TestGetX11Names(t *testing.T) {
	var colors []color_x11 = readAndTestColorsX11(t)
	var names = []string{}
	for _, color := range colors {
		names = append(names, color.Name)
	}
	require.Greater(t, len(names), 0)
}

func TestIsClean(t *testing.T) {
	assert.False(t, isClean("antique white"))
	assert.True(t, isClean("AntiqueWhite"))
	assert.True(t, isClean("AntiqueWhite1"))
}

func TestIsEndsWithDigit(t *testing.T) {
	s0 := "hello"
	s1 := "hello9"
	s2 := "hello999"
	s3 := "hell9o"
	assert.False(t, isEndsWithDigit(s0))
	assert.True(t, isEndsWithDigit(s1))
	assert.True(t, isEndsWithDigit(s2))
	assert.False(t, isEndsWithDigit(s3))
}

func TestIsLower (t *testing.T) {
	s0 := "h"
	// s1 := "H"
	// s2 := ""
	// s3 := "hello"
	// s4 := "hELLO"
	// s5 := "HELLO"
	// s6 := "Hello"
	// s7 := "hElLo"
	// s8 := "HeLlO"
	// s9 := "13"
	assert.True(t, isLowercase(s0))
}

func TestHasIsSpace(t *testing.T) {
	s0 := "hello"
	s1 := " hello"
	s2 := "hello "
	s3 := "he llo"
	assert.False(t, isHasSpace(s0))
	assert.True(t, isHasSpace(s1))
	assert.True(t, isHasSpace(s2))
	assert.True(t, isHasSpace(s3))
}

func TestNormalize(t *testing.T) {
	assert.Equal(t, "AntiqueWhite", normalize("antique white"))
	assert.Equal(t, "LightGoldenrodYellow", normalize("light goldenrod yellow"))
	assert.Equal(t, "DeepSkyBlue", normalize("deep sky blue"))
}

func TestReadColors256Json(t *testing.T) {
	var colors map[int]color_ansi256 = readAndTestColorsAnsi256(t)
	t.Logf("%#v", colors)
}

func TestReadColorsSysJson(t *testing.T) {
	var colors []Color_sys = readAndTestColorsSys(t)
	t.Logf("%#v", colors)
}

func TestReadColorsX11Json(t *testing.T) {
	var colors []color_x11 = readAndTestColorsX11(t)
	t.Logf("%#v", colors)
}

func TestStripTrailingDigits(t *testing.T) {
	assert.Equal(t, "AntiqueWhite", stripTrailingDigits("AntiqueWhite1"))
	assert.Equal(t, "AntiqueWhite", stripTrailingDigits("AntiqueWhite"))
	assert.Equal(t, "firebrick", stripTrailingDigits("firebrick4"))
	assert.Equal(t, "DodgerBlue", stripTrailingDigits("DodgerBlue2"))
}

// func TestRgb(t *testing.T) {
// 	white := colorStruct{R: 255, G: 255, B: 255, Name: "white"}
// 	rgb := white.rgb()
// 	t.Logf("rgb=%X", rgb)
// }

// func TestStyledString(t *testing.T) {
// 	var ss Styledstring = "hello"
// 	ss.FgBg(color_ansi256.Color201, Ansi256.Color194)
// 	fmt.Println(ss)
// }

// func TestWriteAnsi256(t *testing.T) {
// 	sBg := Ansi256.Color194.AnsiBg()
// 	sFg := Ansi256.Color201.AnsiFg()
// 	s := fmt.Sprintf("%s%shello%s", sFg, sBg, Ansi.Reset)
// 	fmt.Println()
// 	fmt.Println()
// 	fmt.Println(s)
// 	fmt.Println()
// 	fmt.Println()
// }

// func TestWriteSys(t *testing.T) {
// 	sBg := Sys.Green.AnsiBg()
// 	sFg := Sys.Red.AnsiFg()
// 	s := fmt.Sprintf("%s%shello%s", sFg, sBg, Ansi.Reset)
// 	fmt.Println()
// 	fmt.Println()
// 	fmt.Println(s)
// 	fmt.Println()
// 	fmt.Println()
// }

func TestWriteX11(t *testing.T) {
	sBg := X11.AntiqueWhite.AnsiBg()
	sFg := X11.DodgerBlue.AnsiFg()
	s := fmt.Sprintf("%s%shello%s", sFg, sBg, Ansi.Reset)
	fmt.Println()
	fmt.Println()
	fmt.Println(s)
	fmt.Println()
	fmt.Println()
}

func TestReadAndTestColorsAnsi256(t *testing.T) {
	_ = readAndTestColorsAnsi256(t)
}

func TestReadAndTestColorsSys(t *testing.T) {
	_ = readAndTestColorsSys(t)
}

// Test method for helper function
func TestReadAndTestColorsX11(t *testing.T) {
	_ = readAndTestColorsX11(t)
}

// Check for colors with the same name that differ only by case
func TestX11CheckForDupes(t *testing.T) {
	var data = map[string][]color_x11{}
	var colors = readAndTestColorsX11(t)
	
	// Create map keyed by lowercased name
	for _, color := range colors {
		name := strings.ToLower(color.Name)
		colorList, is := data[name]
		if !is {
			colorList = []color_x11{}
		}
		colorList = append(colorList, color)
		data[name] = colorList
	}
	t.Logf("len(colors)=%d", len(colors))

	// Create and print histogram of color frequencies aka list lengths
	var dupeFreqs = map[int]int{}
	for _, colorList := range data {
		dupeFreqs[len(colorList)]++
	}
	t.Logf("#%v", dupeFreqs)

	require.Equalf(t, 1, len(dupeFreqs), "Only expecting to see duplication frequencies of 1 (ie no dupes)")
	n, ok := dupeFreqs[1]
	require.True(t, ok)
	require.Equalf(t, len(data), n, "Only expecting to see duplication frequencies of 1 (ie no dupes)")
}

func createCapitalColorAliases (colors []color_x11) []capColorStruct {
	var capColors = make([]capColorStruct, 0, len(colors))
	for _, color := range colors {
		if isLowercase(color.Name) {
			capName := createCapitalName(color.Name)
			var capColor = color_x11 {
				colorbase: color.colorbase,
				Name: capName,
			}
			capColors = append(capColors, capColorStruct{CapColor: capColor, OriginalName: color.Name})
		}
	}

	return capColors
}

func createCapitalName (name string) string {
	if name == "" { return "" }
	return strings.ToUpper(name[0:1]) + name[1:]
}

func createCleanColors (colors []color_x11) []color_x11 {
	var cleanColors []color_x11 = make([]color_x11, 0, len(colors))
	for _, color := range colors {
		if !strings.Contains(color.Name, " ") {
			cleanColors = append(cleanColors, color)
		}
	}

	return cleanColors
}

func genColorsX11Go(t *testing.T, w io.Writer) {
	tmpl := template.New("go")
	bufTmpl, err := os.ReadFile("./colors_x11.go.tmpl")
	require.NoError(t, err)
	_, err = tmpl.Parse(string(bufTmpl))
	require.NoError(t, err)

	var colors []color_x11 = readAndTestColorsX11(t)
	var cleanColors []color_x11 = createCleanColors(colors)
	var capAliases []capColorStruct = createCapitalColorAliases(cleanColors)

	// Anonymous struct to hold template data
	var templateData = struct {
		CleanColors []color_x11
		CapAliases []capColorStruct
	}{
		CleanColors: cleanColors,
		CapAliases: capAliases,
	}
	
	err = tmpl.Execute(w, templateData)
	require.NoError(t, err)
}

// Clean names are nice and simple names. Nothing but alnum.
func isClean(name string) bool {
	rx := regexp.MustCompile("^[[:alnum:]]*$")

	return rx.MatchString(name)
}

func isEndsWithDigit(name string) bool {
	rx := regexp.MustCompile(`.*\d$`)

	return rx.MatchString(name)
}

func isHasSpace(name string) bool {
	return strings.ContainsRune(name, ' ')
}

func isLowercase (name string) bool {
	rx := regexp.MustCompile("^[a-z].*")
	return rx.Match([]byte(name))
}

func loadAndTestTemplate (t *testing.T, path string) *template.Template {
	tmpl := template.New("go")
	bufTmpl, err := os.ReadFile(path)
	require.NoError(t, err)
	_, err = tmpl.Parse(string(bufTmpl))
	require.NoError(t, err)

	return tmpl
}

func normalize(name string) string {
	caser := cases.Title(language.AmericanEnglish)
	titledName := caser.String(name)

	b := strings.Builder{}
	b.Grow(len(titledName))
	for _, c := range titledName {
		if !unicode.IsSpace(c) {
			b.WriteRune(c)
		}
	}
	normalName := b.String()

	return normalName
}

func readAndTestColorsAnsi256(t *testing.T) map[int]color_ansi256 {
	path := "./colors_ansi256.json"
	r, err := os.Open(path)
	require.NoError(t, err)
	require.NotNil(t, r)

	var colors map[int]color_ansi256
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&colors)
	require.NoError(t, err)
	require.NotNil(t, colors)
	require.Equal(t, 256, len(colors))

	return colors
}

func readAndTestColorsSys(t *testing.T) []Color_sys {
	path := "./colors_sys.json"
	r, err := os.Open(path)
	require.NoError(t, err)
	require.NotNil(t, r)

	var colors []Color_sys
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&colors)
	require.NoError(t, err)
	require.NotNil(t, colors)
	require.Equal(t, 17, len(colors))

	return colors
}

// Helper function
// - Read in the file of x11 colors.
// - Decode the file s valid JSON
// = Confirm the file contains 1 or more colors
func readAndTestColorsX11(t *testing.T) []color_x11 {
	path := "./colors_x11.json"
	r, err := os.Open(path)
	require.NoError(t, err)
	require.NotNil(t, r)

	var colors []color_x11
	decoder := json.NewDecoder(r)
	err = decoder.Decode(&colors)
	require.NoError(t, err)
	require.Greater(t, len(colors), 0)

	return colors
}

// Return a string with the trailing digit or digits
// removed, if any.
func stripTrailingDigits(s string) string {
	rx := regexp.MustCompile(`^(.*)(\d+)$`)
	matches := rx.FindStringSubmatch(s)
	if len(matches) == 0 {
		fmt.Printf("No match - %s (len(matches)=%d)\n", s, len(matches))
		return s
	}

	return matches[1]
}
