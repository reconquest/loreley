package loreley

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Compile_CompilesEmptyStringToEmptyStyle(t *testing.T) {
	assertExecutedTemplate(t, ``, ``, nil)
}

func Test_Compile_CompilesStringToString(t *testing.T) {
	assertExecutedTemplate(t, `root beer`, `root beer`, nil)
}

func Test_Compile_CompilesAsGoTemplate(t *testing.T) {
	assertExecutedTemplate(
		t,
		`{if .sweet}bubblegum{end}`,
		`bubblegum`,
		map[string]interface{}{"sweet": true},
	)
}

func Test_Style_Execute_PutsEscapeSequenceToChangeBgColor(t *testing.T) {
	assertExecutedTemplate(t, `{bg 2}finn`, "\x1b[48;5;2mfinn", nil)
}

func Test_Style_Execute_PutsEscapeSequenceToChangeFgColor(t *testing.T) {
	assertExecutedTemplate(t, `{fg 2}finn`, "\x1b[38;5;2mfinn", nil)
}

func Test_Style_Execute_ResetsBackgroundColorToDefault(t *testing.T) {
	assertExecutedTemplate(t, `{nobg}finn`, "\x1b[49mfinn", nil)
}

func Test_Style_Execute_ResetsForegroundColorToDefault(t *testing.T) {
	assertExecutedTemplate(t, `{nofg}finn`, "\x1b[39mfinn", nil)
}

func Test_Style_Execute_PutsReverseCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{reverse}jake`, "\x1b[7mjake", nil)
}

func Test_Style_Execute_PutsNoReverseCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{noreverse}jake`, "\x1b[27mjake", nil)
}

func Test_Style_Execute_PutsBoldCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{bold}jake`, "\x1b[1mjake", nil)
}

func Test_Style_Execute_PutsNoBoldCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{nobold}finn`, "\x1b[22mfinn", nil)
}

func Test_Style_Execute_PutsResetCodeSequence(t *testing.T) {
	assertExecutedTemplate(t, `{reset}finn`, "\x1b[0mfinn", nil)
}

func Test_Style_Execute_PutsBgWithFg(t *testing.T) {
	assertExecutedTemplate(
		t,
		`{fg 1}{bg 2}finn`,
		"\x1b[38;5;1m\x1b[48;5;2mfinn", nil,
	)
}

func Test_Style_Execute_PutsTransitionStringFromOneBgToAnother(t *testing.T) {
	assertExecutedTemplate(
		t,
		`{fg 6}{bg 2}finn{from "" 4}jake`,
		"\x1b[38;5;6m\x1b[48;5;2mfinn"+
			"\x1b[38;5;2m\x1b[48;5;4m\x1b[38;5;6mjake", nil,
	)
}

func Test_Style_Execute_PutsTransitionStringFromOneBgToAnotherInverted(
	t *testing.T,
) {
	assertExecutedTemplate(
		t,
		`{fg 6}{bg 2}finn{to 4 ""}jake`,
		"\x1b[38;5;6m\x1b[48;5;2mfinn"+
			"\x1b[38;5;4m\x1b[48;5;4m\x1b[38;5;6mjake", nil,
	)
}

func Test_CompileWithReset_AddsResetToEndOfStyle(t *testing.T) {
	test := assert.New(t)

	style, err := CompileWithReset(`123`, nil)
	test.Nil(err)

	actual, err := style.ExecuteToString(nil)
	test.Nil(err)
	test.Equal("123\x1b[0m", actual)
}

func Text_TrimStyles_RemovesAllEscapeCodesFromString(t *testing.T) {
	assert.New(t).Equal(
		`finn`,
		TrimStyles(
			"\x1b[38;5;1m\x1b[48;5;2mfinn",
		),
	)
}

func assertExecutedTemplate(
	t *testing.T,
	template string,
	expected string,
	data map[string]interface{},
) {
	test := assert.New(t)

	style, err := Compile(template, nil)
	test.Nil(err)

	actual, err := style.ExecuteToString(data)
	test.Nil(err)
	test.Equal(expected, actual)
}
