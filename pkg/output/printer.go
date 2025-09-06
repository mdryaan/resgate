package output

import (
	"fmt"
	"os"
)

func Println(s string)                              { fmt.Fprintln(os.Stdout, s) }
func Printf(format string, a ...interface{})        { fmt.Fprintf(os.Stdout, format, a...) }
func Header(title string)                           { fmt.Fprintln(os.Stdout, Bold(title)) }
func Successf(format string, a ...interface{})      { fmt.Fprintln(os.Stdout, Success(fmt.Sprintf(format, a...))) }
func Warnf(format string, a ...interface{})         { fmt.Fprintln(os.Stdout, Warn(fmt.Sprintf(format, a...))) }
func Infof(format string, a ...interface{})         { fmt.Fprintln(os.Stdout, Info(fmt.Sprintf(format, a...))) }
func Errorf(format string, a ...interface{})        { fmt.Fprintln(os.Stderr, Error(fmt.Sprintf(format, a...))) }

func Fatal(msg string) {
	Errorf("%s", msg)
	os.Exit(1)
}
