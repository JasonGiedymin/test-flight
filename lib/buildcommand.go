package lib

type BuildCommand struct {
  Controls *FlightControls
  App      *TestFlight
  Options  *CommandOptions
  // Dir      *string //`short:"d" long:"dir" description:"directory to run in"`
  // SingleFileMode *bool //`short:"s" long:"singlefile" description:"single ansible file to use"`
}