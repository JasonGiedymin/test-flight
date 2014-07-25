package lib

import (
    "encoding/json"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
)

func toBytes(data interface{}) ([]byte, error) {
    bytes, err := json.Marshal(data)
    if err != nil {
        return nil, err
    }
    return bytes, nil
}

// Internal -------------------------------------------------------------------
type DeletedContainer struct {
    Name string
    Id   string
}

// end Internal ---------------------------------------------------------------

//TODO: Merge with config
type ApplicationMeta struct {
    Version     string
    LastCommand string
    CurrentMode string
    ExecPath    string
    Pwd         string
    Dir         string
}

type ApplicationState struct {
    Meta *ApplicationMeta
    // Options    CommandOptions
    ConfigFile *ConfigFile
    BuildFile  *BuildFile
}

func (appState *ApplicationState) SetState(newState string) string {
    appState.Meta.CurrentMode = newState
    Logger.Trace("STATE changed to", appState.Meta.CurrentMode)
    return appState.Meta.CurrentMode
}

type RequiredFile struct {
    Name          string
    FileName      string
    FileType      string // [f]ile, [d]ir
    RequiredFiles []RequiredFile
}

// type MapLike interface {
//   Map(v interface{})
// }

type requiredFileMapFunc func(in RequiredFile) RequiredFile

// requiredFile.Map(MapLike)
func (req *RequiredFile) Map(fx requiredFileMapFunc) []RequiredFile {
    returnSlice := make([]RequiredFile, len(req.RequiredFiles))
    return returnSlice
}
