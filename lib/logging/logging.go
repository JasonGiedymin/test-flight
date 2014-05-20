package logging

import (
  "os"
  "github.com/kdar/factorlog"
)

var Info = factorlog.New(os.Stdout,factorlog.NewStdFormatter(`%{Color "green"}%{Date} %{Time} - %{Message}%{Color "reset"}`))
var Error = factorlog.New(os.Stdout,factorlog.NewStdFormatter(`%{Color "red"}%{Date} %{Time} %{File}:%{Line} %{Message}%{Color "reset"}`))
