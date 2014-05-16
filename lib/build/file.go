package build

import (
  "io/ioutil"
  "encoding/json"
  "fmt"
)

type BuildFile struct {
  Name string
  Parent string
  version string
}

func TryJson(path string) {
  jsonBlob, _ := ioutil.ReadFile(path+"/build.json")

  var buildFile BuildFile
	err := json.Unmarshal(jsonBlob, &buildFile)
	if err != nil {
		fmt.Println("error:", err)
	}
	fmt.Printf("%+v\n", buildFile)
}
