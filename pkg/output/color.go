package output

import "github.com/fatih/color"

var (
	greenFn  = color.New(color.FgGreen).SprintFunc()
	yellowFn = color.New(color.FgYellow).SprintFunc()
	redFn    = color.New(color.FgRed).SprintFunc()
	cyanFn   = color.New(color.FgCyan).SprintFunc()
	boldFn   = color.New(color.Bold).SprintFunc()
	dimFn    = color.New(color.Faint).SprintFunc()
)

func Green(s string) string  { return greenFn(s) }
func Yellow(s string) string { return yellowFn(s) }
func Red(s string) string    { return redFn(s) }
func Cyan(s string) string   { return cyanFn(s) }
func Bold(s string) string   { return boldFn(s) }
func Dim(s string) string    { return dimFn(s) }

func Success(s string) string { return greenFn(s) }
func Warn(s string) string    { return yellowFn(s) }
func Error(s string) string   { return redFn(s) }
func Info(s string) string    { return cyanFn(s) }
