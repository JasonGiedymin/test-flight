package lib

// == Check Command ==
type CheckCommand struct {
  Controls *FlightControls
  App      *TestFlight
  Dir      *string //`short:"d" long:"dir" description:"directory to run in"`
  SingleFileMode *bool //`short:"s" long:"singlefile" description:"single ansible file to use"`
}
