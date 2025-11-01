package output

import (
	"fmt"
	"io"
	"os"
)

var out io.Writer = os.Stdout

func SetWriter(w io.Writer) {
	out = w
}

func Info(msg string) {
	fmt.Fprintln(out, msg)
}

func Error(msg string) {
	fmt.Fprintln(out, "❌ "+msg)
}

func Success(msg string) {
	fmt.Fprintln(out, "✅ "+msg)
}
