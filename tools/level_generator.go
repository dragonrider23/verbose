// +build ignore

package main

import (
	"flag"
	"fmt"
	"os"
	"text/template"
)

var outFile string

var levels = [...]string{
	"Debug",
	"Info",
	"Notice",
	"Warning",
	"Error",
	"Critical",
	"Alert",
	"Emergency",
	"Fatal",
}

var fileTmpl = `// This file was generated with level_generator. DO NOT EDIT

package verbose

import "os"

// LogLevel is used to compare levels in a consistant manner
type LogLevel int

// String returns the stringified version of LogLevel.
// I.e., "Error" for LogLevelError, and "Debug" for LogLevelDebug
// It will return an empty string for any undefined level.
func (l LogLevel) String() string {
	if s, ok := levelString[l]; ok {
		return s
	}
	return ""
}

// These are the defined log levels
const ({{range $i, $l := .}}
	LogLevel{{$l}}{{if eq $i 0}} LogLevel = iota{{end}}{{end}}
)

// LogLevel to stringified versions
var levelString = map[LogLevel]string{ {{range .}}
	LogLevel{{.}}:     "{{.}}",{{end}}
}
{{range .}}
// {{.}} - Log {{.}} message
func (l *Logger) {{.}}(v ...interface{}) {
    NewEntry(l).{{.}}(v...){{if eq . "Fatal"}}
	os.Exit(1){{end}}
    return
}
{{end}}
// Panic - Log Panic message
func (l *Logger) Panic(v ...interface{}) {
    NewEntry(l).Panic(v...)
    return
}

// Print - Log Print message
func (l *Logger) Print(v ...interface{}) {
    NewEntry(l).Print(v...)
    return
}

// Printf friendly functions
{{range .}}
// {{.}}f - Log formatted {{.}} message
func (l *Logger) {{.}}f(m string, v ...interface{}) {
	NewEntry(l).{{.}}f(m, v...){{if eq . "Fatal"}}
	os.Exit(1){{end}}
    return
}
{{end}}
// Panicf - Log formatted Panic message
func (l *Logger) Panicf(m string, v ...interface{}) {
	NewEntry(l).Panicf(m, v...)
    return
}

// Printf - Log formatted Print message
func (l *Logger) Printf(m string, v ...interface{}) {
	NewEntry(l).Printf(m, v...)
    return
}

// Println friendly functions
{{range .}}
// {{.}}ln - Log {{.}} message with newline
func (l *Logger) {{.}}ln(v ...interface{}) {
    NewEntry(l).{{.}}ln(v...){{if eq . "Fatal"}}
	os.Exit(1){{end}}
    return
}
{{end}}
// Panicln - Log Panic message with newline
func (l *Logger) Panicln(v ...interface{}) {
    NewEntry(l).Panicln(v...)
    return
}

// Println - Log Print message with newline
func (l *Logger) Println(v ...interface{}) {
    NewEntry(l).Println(v...)
    return
}
`

func init() {
	flag.StringVar(&outFile, "output", "log_levels.go", "Filename of output")
}

func main() {
	flag.Parse()

	tmpl, err := template.New("funcs").Parse(fileTmpl)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}

	file, err := os.OpenFile(outFile, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		fmt.Print(err)
		os.Exit(1)
	}
	defer file.Close()
	tmpl.Execute(file, levels)
}
