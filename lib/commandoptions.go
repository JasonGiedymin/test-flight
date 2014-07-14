package lib

type CommandOptions struct {
    Verbose        []bool `short:"v" long:"verbose" description:"verbose output"`
    Configfile     string `short:"c" long:"config" description:"test-flight config file to use"`
    Dir            string `short:"d" long:"dir" description:"directory to run in"`
    SingleFileMode bool   `short:"s" long:"singlefile" description:"single ansible file to use"`
    Force          bool   `short:"f" long:"force" description:"force new image"`
}
