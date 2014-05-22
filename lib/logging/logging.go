package logging

import (
  "github.com/kdar/factorlog"
  "os"
)

// The actual loggers
var (
  Log       *factorlog.FactorLog
  LogDebug  *factorlog.FactorLog
  File      *factorlog.FactorLog
  debugFile *os.File = getDebugFile()
)

var debugLogFormat = `%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "white+b" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}] - %{Message}%{Color "reset"}`
var stdLogFormat = `%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "white+b" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}] - %{Message}%{Color "reset"}`
var fileLogFormat = `[%{Date} %{Time}] [%{SEVERITY}] - %{Message}%{Color "reset"}`

func getDebugFile() *os.File {
  newFile, _ := os.Create("debug.log")
  return newFile
}

// Sets up Loggers
func setup() {
  Log = factorlog.New(
    os.Stdout,
    factorlog.NewStdFormatter(stdLogFormat))

  LogDebug = factorlog.New(
    os.Stdout,
    factorlog.NewStdFormatter(debugLogFormat))

  File = factorlog.New(
    debugFile,
    factorlog.NewStdFormatter(fileLogFormat))
}

// Checks to see if Loggers should be set
func load() {
  if Log == nil {
    setup()
  }
}

// Absolute Errors
func Error(v ...interface{}) {
  if Log == nil {
    setup()
  }

  Log.Error(v)
  // File.Println(v)
}

// Things that aren't errors but you should know
func Warn(v ...interface{}) {
  if Log == nil {
    setup()
  }

  Log.Warn(v)
}

// General Application Info
func Info(v ...interface{}) {
  if Log == nil {
    setup()
  }

  Log.Info(v)
}

// More Info, params, situation, etc...
func Debug(v ...interface{}) {
  if Log == nil {
    setup()
  }

  LogDebug.Debug(v)
}

// Use Trace when run debugging for debug builds for fine grain info
func Trace(v ...interface{}) {
  if Log == nil {
    setup()
  }

  Log.Trace(v)
}
