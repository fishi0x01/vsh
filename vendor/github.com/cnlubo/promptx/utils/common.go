package utils

import (
	"bytes"
	"fmt"
	"github.com/pkg/errors"
	"os"
	"text/template"
)

type winsize struct {
	Row    uint16
	Col    uint16
	Xpixel uint16
	Ypixel uint16
}

// ErrEOF is the error returned from prompts when EOF is encountered.
var ErrEOF = errors.New("^D")

// ErrInterrupt is the error returned from prompts when an interrupt (ctrl-c) is
// encountered.
var ErrInterrupt = errors.New("^C")

// ErrAbort is the error returned when confirm prompts are supplied "n"
var ErrAbort = errors.New("")

func CheckErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return false
	}
	return true
}

func CheckAndExit(err error) {
	if !CheckErr(err) {
		os.Exit(1)
	}
}

func Render(tpl *template.Template, data interface{}) []byte {
	var buf bytes.Buffer
	err := tpl.Execute(&buf, data)
	if err != nil {
		return []byte(fmt.Sprintf("%v", data))
	}
	return buf.Bytes()
}
