package log

import (
	"fmt"
	au "github.com/logrusorgru/aurora"
	"log"
	"os"
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

var verbosityLevel int

var (
	appDebug   *log.Logger
	appInfo    *log.Logger
	appWarning *log.Logger
	appError   *log.Logger

	userDebug *log.Logger
	userInfo  *log.Logger
	userError *log.Logger
)

// Init initialize the log with given verbosity level
func Init(verbosity string) error {
	switch verbosity {
	case "DEBUG":
		verbosityLevel = debugLvl
	case "INFO":
		verbosityLevel = infoLvl
	case "WARN":
		verbosityLevel = warnLvl
	case "ERROR":
		verbosityLevel = errorLvl
	default:
		fmt.Printf("Not a valid verbosity level: %s\n", verbosity)
		fmt.Println("Allowed values are DEBUG | INFO | WARN | ERROR")
		return fmt.Errorf("not a valid verbosity level: %s", verbosity)
	}

	if verbosityLevel <= debugLvl {
		file, err := os.OpenFile("vsh_trace.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err != nil {
			log.Fatal(err)
		}

		appDebug = log.New(file,
			"debugLvl: ",
			log.Ldate|log.Ltime)
	}

	if verbosityLevel <= infoLvl {
		appInfo = log.New(os.Stdout,
			"infoLvl: ",
			log.Ldate|log.Ltime)
	}

	if verbosityLevel <= warnLvl {
		appWarning = log.New(os.Stdout,
			"WARNING: ",
			log.Ldate|log.Ltime)
	}

	if verbosityLevel <= errorLvl {
		appError = log.New(os.Stderr,
			"errorLvl: ",
			log.Ldate|log.Ltime)
	}

	if verbosityLevel <= debugLvl {
		userDebug = log.New(os.Stdout, "", 0)
	}

	userInfo = log.New(os.Stdout, "", 0)
	userError = log.New(os.Stderr, "", 0)

	return nil
}

func getCustomPrefix(logger *log.Logger) string {
	_, fi, line, _ := runtime.Caller(2)
	loc := filepath.Base(fi) + ":" + strconv.Itoa(line)
	return logger.Prefix() + loc + " "
}

// AppTrace log application trace
func AppTrace(f string, args ...interface{}) {
	if appDebug != nil {
		appDebug.SetPrefix(getCustomPrefix(appDebug))
		appDebug.Printf(f, args...)
	}
}

// AppInfo log application infoLvl
func AppInfo(f string, args ...interface{}) {
	if appInfo != nil {
		appInfo.SetPrefix(getCustomPrefix(appInfo))
		appInfo.Printf(f, args...)
	}
}

// AppWarning log application warning
func AppWarning(f string, args ...interface{}) {
	if appWarning != nil {
		appWarning.SetPrefix(getCustomPrefix(appWarning))
		appWarning.Printf(f, args...)
	}
}

// AppError log application error
func AppError(f string, args ...interface{}) {
	if appError != nil {
		appError.SetPrefix(getCustomPrefix(appError))
		appError.Printf(f, args...)
	}
}

// UserDebug log user debugLvl
func UserDebug(f string, args ...interface{}) {
	if userDebug != nil {
		userDebug.Printf(f, args...)
	}
}

// UserInfo log user infoLvl
func UserInfo(f string, args ...interface{}) {
	userInfo.Printf(f, args...)
}

// UserError log user error
func UserError(f string, args ...interface{}) {
	message := au.Red(au.Sprintf(f, args...))
	userError.Println(au.Sprintf(message))
}
