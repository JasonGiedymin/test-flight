package lib

// == Launch Command ==
type LaunchCommand struct {
  Controls *FlightControls
  Options  *CommandOptions
  App      *TestFlight
}