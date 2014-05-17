package build

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
)

type BuildFile struct {
  Name    string
  Parent  string
  Version string
}

func ReadBuildFile(path string) {
  jsonBlob, _ := ioutil.ReadFile(path + "/build.json")

  var buildFile BuildFile
  err := json.Unmarshal(jsonBlob, &buildFile)
  if err != nil {
    fmt.Println("error:", err)
  }
  fmt.Printf("%+v\n", buildFile)
}
