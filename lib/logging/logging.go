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

var debugLogFormat = `%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "blue" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}:%{File}:%{Line}] %{Message}%{Color "reset"}`
var stdLogFormat = `%{Color "red" "ERROR"}%{Color "yellow" "WARN"}%{Color "green" "INFO"}%{Color "cyan" "DEBUG"}%{Color "blue" "TRACE"}[%{Date} %{Time}] [%{SEVERITY}] - %{Message}%{Color "reset"}`
var fileLogFormat = `[%{Date} %{Time}] [%{SEVERITY}:%{File}:%{Line}] %{Message}%{Color "reset"}`

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

func Error(v ...interface{}) {
  if Log == nil {
    setup()
  }

  Log.Error(v)
  File.Println(v)
}

func Warn(v ...interface{}) {
  if Log == nil {
    setup()
  }

  Log.Warn(v)
}

func Info(v ...interface{}) {
  if Log == nil {
    setup()
  }

  Log.Info(v)
}

func Debug(v ...interface{}) {
  if Log == nil {
    setup()
  }

  LogDebug.Debug(v)
}

func Trace(v ...interface{}) {
  if Log == nil {
    setup()
  }

  Log.Trace(v)
}
