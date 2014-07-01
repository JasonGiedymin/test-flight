package lib

// == Destroy Command ==
// Should destroy running containers
type DestroyCommand struct {
  Controls *FlightControls
  App      *TestFlight
  Options  *CommandOptions
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  SingleFileMode bool `short:"s" long:"singlefile" description:"single ansible file to use"`
}