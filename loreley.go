package loreley

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"text/template"
)

const (
	// CodeStart is an escape sequence for starting color coding
	CodeStart = "\x1b"

	// CodeEnd is an escape sequence fo ending color coding
	CodeEnd = `m`

	// AttrForeground is an escape sequence part for setting foreground
	AttrForeground = `3`

	// AttrBackground is an escape sequence part for setting background
	AttrBackground = `4`

	// AttrDefault is an escape sequence part for resetting foreground
	// or background
	AttrDefault = `9`

	// AttrReset is an esscape sequence part for resetting all attributes
	AttrReset = `0`

	// AtrrReverse is an escape sequence part for setting reverse display
	AttrReverse = `7`

	// AttrNoReverse is an escape sequence part for setting reverse display off
	AttrNoReverse = `27`

	// AttrBold is an escape sequence part for setting bold mode on
	AttrBold = `1`

	// AttrNoBold is an escape sequence part for setting bold mode off
	AttrNoBold = `22`

	// AttrForeground256 is an escape sequence part for setting foreground
	// color in 256 color mode
	AttrForeground256 = `38;5`

	// AttrBackground256 is an escape sequence part for setting background
	// color in 256 color mode
	AttrBackground256 = `48;5`

	// StyleReset is an placeholder for resetting all attributes.
	StyleReset = `{reset}`

	// DefaultColor represents default foreground and background color.
	DefaultColor = 0
)

var (
	// CodeRegexp is an regular expression for matching escape codes.
	CodeRegexp = regexp.MustCompile(CodeStart + `[^` + CodeEnd + `]+`)

	// DelimLeft is used for match template syntax (see Go-lang templates).
	DelimLeft = `{`

	// DelimRight is used for match template syntax (see Go-lang templates).
	DelimRight = `}`
)

// Style is a compiled style. Can be used as text/template.
type Style struct {
	*template.Template

	background int
	foreground int

	// NoColors can be set to true to disable escape sequence output,
	// but still render underlying template.
	NoColors bool
}

// ExecuteToString is same, as text/template Execute, but return string
// as a result.
func (style *Style) ExecuteToString(
	data map[string]interface{},
) (string, error) {
	buffer := bytes.Buffer{}

	err := style.Template.Execute(&buffer, data)
	if err != nil {
		return "", err
	}

	return buffer.String(), nil
}

func (style *Style) putReset() string {
	style.background = DefaultColor
	style.foreground = DefaultColor

	return style.getStyleCodes(AttrReset)
}

func (style *Style) putBackground(color int) string {
	style.background = color

	return style.getStyleCodes(AttrBackground256, fmt.Sprint(color))
}

func (style *Style) putDefaultBackground() string {
	style.background = DefaultColor

	return style.getStyleCodes(
		AttrBackground + AttrDefault,
	)
}

func (style *Style) putDefaultForeground() string {
	style.foreground = DefaultColor

	return style.getStyleCodes(
		AttrForeground + AttrDefault,
	)
}

func (style *Style) putForeground(color int) string {
	style.foreground = color

	return style.getStyleCodes(AttrForeground256, fmt.Sprint(color))
}

func (style *Style) putBold() string {
	return style.getStyleCodes(AttrBold)
}

func (style *Style) putNoBold() string {
	return style.getStyleCodes(AttrNoBold)
}

func (style *Style) putReverse() string {
	return style.getStyleCodes(AttrReverse)
}

func (style *Style) putNoReverse() string {
	return style.getStyleCodes(AttrNoReverse)
}

func (style *Style) putTransitionFrom(text string, nextBackground int) string {
	previousBackground := style.background
	previousForeground := style.foreground

	return style.putForeground(previousBackground) +
		style.putBackground(nextBackground) +
		text +
		style.putForeground(previousForeground)
}

func (style *Style) putTransitionTo(nextBackground int, text string) string {
	previousForeground := style.foreground

	return style.putForeground(nextBackground) +
		text +
		style.putBackground(nextBackground) +
		style.putForeground(previousForeground)
}

func (style *Style) getStyleCodes(attr ...string) string {
	if style.NoColors {
		return ``
	}

	return fmt.Sprintf(`%s[%sm`, CodeStart, strings.Join(attr, `;`))
}

// Compile compiles style which is specified by string into Style,
// optionally adding extended formatting functions.
func Compile(
	text string,
	extensions map[string]interface{},
) (*Style, error) {
	style := &Style{}

	functions := map[string]interface{}{
		`bg`:        style.putBackground,
		`fg`:        style.putForeground,
		`nobg`:      style.putDefaultBackground,
		`nofg`:      style.putDefaultForeground,
		`bold`:      style.putBold,
		`nobold`:    style.putNoBold,
		`reverse`:   style.putReverse,
		`noreverse`: style.putNoReverse,
		`reset`:     style.putReset,
		`from`:      style.putTransitionFrom,
		`to`:        style.putTransitionTo,
	}

	for name, function := range extensions {
		functions[name] = function
	}

	template, err := template.New(`style`).Delims(
		DelimLeft,
		DelimRight,
	).Funcs(functions).Parse(
		text,
	)

	if err != nil {
		return nil, err
	}

	return &Style{
		Template: template,
	}, nil
}

// CompileWithReset is same as Compile, but appends {reset} to the end
// of given style string.
func CompileWithReset(
	text string,
	extensions map[string]interface{},
) (*Style, error) {
	return Compile(text+StyleReset, extensions)
}

// TrimStyles removes all escape codes from the given string.
func TrimStyles(input string) string {
	return CodeRegexp.ReplaceAllLiteralString(input, ``)
}
