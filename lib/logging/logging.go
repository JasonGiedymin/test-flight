package logging

import (
    "github.com/kdar/factorlog"
    "os"
)

// The actual loggers
var (
    Log        *factorlog.FactorLog
    LogDebug   *factorlog.FactorLog
    LogConsole *factorlog.FactorLog // channel??
    File       *factorlog.FactorLog
    debugFile  *os.File = getDebugFile()
    verbosity  int
)

const debugLogFormat = `%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "white+b" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}] - %{Message}%{Color "reset"}`
const consoleLogFormat = `%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "white+b" "TRACE"}%{Message}%{Color "reset"}`
const stdLogFormat = `%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "white+b" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}] - %{Message}%{Color "reset"}`
const fileLogFormat = `[%{Date} %{Time}] [%{SEVERITY}] - %{Message}%{Color "reset"}`

var levels = map[int]factorlog.Severity{
    1:  factorlog.INFO,
    2:  factorlog.DEBUG,
    3:  factorlog.TRACE,
    4:  factorlog.TRACE, // 4 is Trace + file logging
}

func getDebugFile() *os.File {
    if verbosity >= 4 {
        newFile, _ := os.Create("debug.log")
        return newFile
    }

    return nil
}

// Sets up Loggers
func Setup() {
    maxLevel := func(verbosity int) int {
        max := len(levels)
        if verbosity > max {
            return max
        }
        return verbosity
    }

    Log = factorlog.New(
        os.Stdout,
        factorlog.NewStdFormatter(stdLogFormat))

    LogDebug = factorlog.New(
        os.Stdout,
        factorlog.NewStdFormatter(debugLogFormat))

    LogConsole = factorlog.New(
        os.Stdout,
        factorlog.NewStdFormatter(consoleLogFormat))

    // TODO: set log level here
    Log.SetMinMaxSeverity(levels[maxLevel(verbosity)], factorlog.ERROR)
    // LogDebug.SetMinMaxSeverity(levels[maxLevel(verbosity)], factorlog.ERROR)
    LogConsole.SetMinMaxSeverity(levels[maxLevel(verbosity)], factorlog.ERROR)

    if verbosity >= 4 {
        File = factorlog.New(
            debugFile,
            factorlog.NewStdFormatter(fileLogFormat))
    }
}

// Checks to see if Loggers should be set
func Load(newVerbosity int) {
    verbosity = newVerbosity
    if Log == nil {
        Setup()
    }
}

// Absolute Errors
func Error(v ...interface{}) {
    if Log == nil {
        Setup()
    }

    Log.Error(v)

    if verbosity >= 4 {
        File.Println(v)
    }
}

// Things that aren't errors but you should know
func Warn(v ...interface{}) {
    if Log == nil {
        Setup()
    }

    Log.Warn(v)
}

// General Application Info
func Info(v ...interface{}) {
    if Log == nil {
        Setup()
    }

    Log.Info(v)
}

// More Info, params, situation, etc...
func Debug(v ...interface{}) {
    if Log == nil {
        Setup()
    }

    Log.Debug(v)
}

// Use Trace when run debugging for debug builds for fine grain info
func Trace(v ...interface{}) {
    if Log == nil {
        Setup()
    }

    Log.Trace(v)
}

func What(v ...interface{}) {
    if Log == nil {
        Setup()
    }

    LogDebug.Trace("*** ", v)
}

func Console(v ...interface{}) {
    if Log == nil {
        Setup()
    }

    LogConsole.Info(v[0])
}
