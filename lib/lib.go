// ## Deps
//     go get github.com/SpaceMonkeyGo/errors
//

package lib

import (
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "github.com/SpaceMonkeyGo/errors"
)

// == Default vars ==
var (
    defaultDir    = "./"
    BadDir        = errors.NewClass("Can't read the current directory")
    FileCheckFail = errors.NewClass("File Check Failed")
    mainYaml      = RequiredFile{Name: "main yaml", FileName: "main.yml", FileType: "f"}
    buildFile     = RequiredFile{Name: Constants().buildFileName, FileName: Constants().buildFileName, FileType: "f"}
    AnsibleFiles  = []RequiredFile{mainYaml}
)

// Creates a list of required files needed by TestFlight using
// templateDir via config file as a basis.
func TestFlightFiles(templateDir string) []RequiredFile {
    return []RequiredFile{
        {Name: "Test-Flight dir", FileName: templateDir, FileType: "d"}, //, These will actually be generated
        // RequiredFiles: []RequiredFile{
        //   {Name: "Ansible inventory file used for Test-Flight", FileName: "inventory", FileType: "f"},
        //   {Name: "Ansible playbook file used for Test-Flight", FileName: "playbook.yml", FileType: "f"},
        // },

    }
}

// Generates required files
// Note: complex types are free form and can be aything
//       where simple types are top level only.
func generateRequiredFilesFrom(buildfile *BuildFile) ([]RequiredFile, error) {
    requiredFiles := []RequiredFile{}

    for _, file := range buildfile.Add.Simple {
        name := "Ansible " + file + " dir"
        requiredFiles = append(requiredFiles, RequiredFile{Name: name, FileName: file, FileType: "d" /*,RequiredFiles: AnsibleFiles*/})
    }

    // Treat Complex as fully qualified paths and as OS Files
    for _, file := range buildfile.Add.Complex {
        if file.Name == "" || file.Location == "" {
            msg := "Found empty values for " + Constants().buildFileName + " Complex entry. Check syntax, try again."
            Logger.Error(msg)
            return nil, errors.New(msg)
        }

        name := "file " + file.Name
        requiredFiles = append(requiredFiles, RequiredFile{Name: name, FileName: file.Location, FileType: "f"})
    }

    Logger.Trace(requiredFiles)

    return requiredFiles, nil
}

var RequiredFiles = []RequiredFile{
    {Name: "Test-Flight json build file", FileName: Constants().buildFileName, FileType: "f"},

    {Name: "Ansible handlers dir", FileName: "handlers", FileType: "d", RequiredFiles: AnsibleFiles},
    {Name: "Ansible meta dir", FileName: "meta", FileType: "d", RequiredFiles: AnsibleFiles},
    {Name: "Ansible tasks dir", FileName: "tasks", FileType: "d", RequiredFiles: AnsibleFiles},
    {Name: "Ansible templates dir", FileName: "templates", FileType: "d"},
    {Name: "Ansible test dir", FileName: "tests", FileType: "d", RequiredFiles: AnsibleFiles},
    {Name: "Ansible defaults dir", FileName: "defaults", FileType: "d", RequiredFiles: AnsibleFiles},
    {Name: "Ansible var dir for variables", FileName: "vars", FileType: "d", RequiredFiles: AnsibleFiles},
    {Name: "Ansible vault dir for encrypted files", FileName: "vault", FileType: "d"},
}
