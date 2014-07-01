package lib

// == Ground Command ==
// Should stop running containers
type GroundCommand struct {
  Controls *FlightControls
  Options  *CommandOptions
  App      *TestFlight
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  SingleFileMode bool `short:"s" long:"singlefile" description:"single ansible file to use"`
}
