// ## Deps
//     go get github.com/SpaceMonkeyGo/errors
//

package lib

import (
  // Logger "./logging"
  // "bitbucket.org/kardianos/osext"
  "github.com/SpaceMonkeyGo/errors"
  // "github.com/jessevdk/go-flags"
  // "os"
)

// == Default vars ==
var (
  defaultDir    = "./"
  mainYaml      = RequiredFile{Name: "main yaml", FileName: "main.yml", FileType: "f"}
  BadDir        = errors.NewClass("Can't read the current directory")
  FileCheckFail = errors.NewClass("File Check Failed")
  AnsibleFiles  = []RequiredFile{mainYaml}
)

// Creates a list of required files needed by TestFlight using
// templateDir via config file as a basis.
func TestFlightFiles(templateDir string) []RequiredFile {
  return []RequiredFile {
    {Name: "Test-Flight dir", FileName: templateDir, FileType: "d",//, These will actually be generated
      // RequiredFiles: []RequiredFile{
      //   {Name: "Ansible inventory file used for Test-Flight", FileName: "inventory", FileType: "f"},
      //   {Name: "Ansible playbook file used for Test-Flight", FileName: "playbook.yml", FileType: "f"},
      // },
    },
  }
}

var RequiredFiles = []RequiredFile{
  {Name: "Test-Flight json build file", FileName: "build.json", FileType: "f"},
  {Name: "Ansible handlers dir", FileName: "handlers", FileType: "d", RequiredFiles: AnsibleFiles},
  {Name: "Ansible meta dir", FileName: "meta", FileType: "d", RequiredFiles: AnsibleFiles},
  {Name: "Ansible tasks dir", FileName: "tasks", FileType: "d", RequiredFiles: AnsibleFiles},
  {Name: "Ansible templates dir", FileName: "templates", FileType: "d"},
  {Name: "Ansible test dir", FileName: "tests", FileType: "d", RequiredFiles: AnsibleFiles},
  {Name: "Ansible var dir for variables", FileName: "vars", FileType: "d", RequiredFiles: AnsibleFiles},
  {Name: "Ansible vault dir for encrypted files", FileName: "vault", FileType: "d"},
}
