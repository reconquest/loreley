# loreley

Easy and extensible colorizer for the programs' output.

Basically, loreley turns this:

```
{bold}{fg 15}{bg 27} hello {from "" 29} there {to 16 ""},
```

Into this:

![2016-06-27-13t53t45](https://raw.githubusercontent.com/reconquest/loreley/master/demo.png)

# Usage

```go
package main

import "fmt"
import "github.com/reconquest/loreley"

func main() {
	text, err := loreley.CompileAndExecuteToString(
		`{bold}{fg 15}{bg 27} hello {from "" 29} {.where} {to 16 ""}{reset}`,
		nil,
		map[string]interface{}{"where": "there"},
	)
	if err != nil {
		fmt.Errorf(`can't compile loreley template: %s`, err)
	}

	fmt.Println(text)
}
```

# Reference

loreley extends Go-lang template system. So, fully syntax is supported with
exception, that `{` and `}` will be used as delimiters.

All `<color>`, accepted by loreley, should be the 256-color code.

Available template functions:

* `{bg <color>}` sets background color for the next text;
* `{fg <color>}` sets foreground color for the next text;
* `{nobg}` resets background color to the default;
* `{nofg}` resets foreground color to the default;
* `{bold}` set bold mode on for the next text;
* `{nobold}` set bold mode to off;
* `{reverse}` set reverse mode on for the next text;
* `{noreverse}` set reverse mode off;
* `{reset}` resets all styles to default;
* `{from <text> <bg>}` reuse current fg as specified `<text>`'s bg color,
  specified `<bg>` will be used as fg color and as bg color for the following
  text;
* `{to <bg> <text>}` reuse current bg as specified `<text>`'s bg color,
  specified `<bg>` will be used as fg color for the following text;
