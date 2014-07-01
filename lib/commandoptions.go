package lib

type CommandOptions struct {
  Configfile      string `short:"c" long:"config" description:"test-flight config file to use"`
  Dir      string `short:"d" long:"dir" description:"directory to run in"`
  SingleFileMode bool `short:"s" long:"singlefile" description:"single ansible file to use"`
  Force    bool   `short:"f" long:"force" description:"force new image"`
  // Check           *lib.CheckCommand
  // Images          *lib.ImagesCommand
  // Build           *lib.BuildCommand
  // Launch          *lib.LaunchCommand
  // Ground          *lib.GroundCommand
  // Destroy         *lib.DestroyCommand
  // Version         *lib.VersionCommand
  // Template        *lib.TemplateCommand
}