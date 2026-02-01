package output

import (
	"fmt"
	"io"
	"os"
)

var out io.Writer = os.Stdout
var errOut io.Writer = os.Stdout

func SetWriter(w io.Writer) {
	out = w
}

func SetErrorWriter(w io.Writer) {
	errOut = w
}

func Info(msg string) {
	_, _ = fmt.Fprintln(out, msg)
}

func Error(msg string) {
	_, _ = fmt.Fprintln(errOut, msg)
}

func Success(msg string) {
	_, _ = fmt.Fprintln(out, msg)
}
