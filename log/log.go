package log

import (
	au "github.com/logrusorgru/aurora"
	"log"
	"path/filepath"
	"runtime"
	"strconv"
)

const (
	debugLvl = iota
	infoLvl
	warnLvl
	errorLvl
)

var verbose = false

// Init initializes the logger
func Init() {
	log.SetFlags(0)
}

func ToggleVerbose() {
	verbose = !verbose
}

func logger(level uint, f string, args ...interface{}) {
	_, fi, line, _ := runtime.Caller(2)
	loc := filepath.Base(fi) + ":" + strconv.Itoa(line) + " "
	var prefix, message au.Value

	switch level {
	case debugLvl:
		prefix = au.Cyan("[DEBUG] " + loc)
		message = au.Cyan(au.Sprintf(f, args...))
	case infoLvl:
		prefix = au.Green("[INFO] " + loc)
		message = au.Green(au.Sprintf(f, args...))
	case warnLvl:
		prefix = au.Yellow("[WARN] " + loc)
		message = au.Yellow(au.Sprintf(f, args...))
	case errorLvl:
		prefix = au.Red("[ERROR] " + loc)
		message = au.Red(au.Sprintf(f, args...))
	default:
		panic("Cannot log!")
	}

	log.SetPrefix(au.Sprintf(au.Bold(prefix)))
	log.Printf(au.Sprintf(message))
	log.SetPrefix("")
}

// Debug writes debug message to log
func Debug(f string, args ...interface{}) {
	if verbose {
		logger(debugLvl, f, args...)
	}
}

// Info writes info message to log
func Info(f string, args ...interface{}) {
	if verbose {
		logger(infoLvl, f, args...)
	}
}

// Warn writes warn message to log
func Warn(f string, args ...interface{}) {
	logger(warnLvl, f, args...)
}

// Error writes error message to log
func Error(f string, args ...interface{}) {
	logger(errorLvl, f, args...)
}
