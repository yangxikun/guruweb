package output

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func Warn(format string, a ...interface{}) {
	fmt.Fprintln(os.Stderr, color.YellowString("WARN: "+format, a...))
}

func Normal(format string, a ...interface{}) {
	fmt.Fprintln(os.Stdout, color.GreenString("> "+format, a...))
}

func Trace(format string, a ...interface{}) {
	fmt.Fprintln(os.Stdout, color.BlueString("TRACE: "+format, a...))
}
