package lib

import (
    "archive/tar"
    "errors"
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"
    "io/ioutil"
    "os"
    "os/signal"
    "strings"
    "sync"
    "syscall"
)

// var Info = factorlog.New(os.Stdout, factorlog.NewStdFormatter(`%{Color "green"}%{Date} %{Time} %{File}:%{Line} %{Message}%{Color "reset"}`))
// var Error = factorlog.New(os.Stdout, factorlog.NewStdFormatter(`%{Color "red"}%{Date} %{Time} %{File}:%{Line} %{Message}%{Color "reset"}`))

// Converts []FileInfo => []string
func ConvertFiles(files []os.FileInfo) []string {
    convertedFiles := []string{}
    for _, value := range files {
        convertedFiles = append(convertedFiles, value.Name())
    }

    return convertedFiles
}

func findFile(filesFound []string, requiredFile RequiredFile, currDir string) (bool, error) {
    for _, file := range filesFound {
        // TODO: inspect complex types here by splitting and taking [0]
        //       will also have to save name[:len(name)] as next
        //       as the below tries to match "file" == "file/path" which will
        //       fail as it is explicit, but should pass. Need unit tests now

        if file == requiredFile.FileName {
            // Logger.What("->", file)
            if len(requiredFile.RequiredFiles) > 0 && requiredFile.FileType == "d" {
                nextDir := currDir + "/" + requiredFile.FileName
                _, err := HasRequiredFiles(nextDir, requiredFile.RequiredFiles)
                if err != nil {
                    return false, err
                }
            }
            return true, nil
        }
    }

    // msg := "Required file/dir not found: [" + currDir + "/" + requiredFile.FileName + "]"
    // FileCheckFail.New("Required file/dir not found: [%v/%v]", currDir, requiredFile.FileName)

    return false, nil //, errors.New(msg)
}

func CreateFile(dir *string, requiredFile RequiredFile) (*os.File, error) {
    var fileName = *dir + "/" + requiredFile.FileName
    var err error
    var file *os.File

    if requiredFile.FileType == "d" {
        if err = os.Mkdir(fileName, 0755); err != nil {
            return nil, err
        }
    } else if requiredFile.FileType == "f" {
        if _, err = os.Create(fileName); err != nil {
            return nil, nil
        }
    }

    return file, nil
}

// TODO: goroutine + channels to further optimize
/*
 * Reads the current directory and returns
 *   - bool: if all the required files were found
 */
func HasRequiredFiles(dir string, requiredFiles []RequiredFile) (bool, error) {
    if dir == "" {
        dir = defaultDir
    }

    filesFromDisk, err := ioutil.ReadDir(dir)

    if err != nil {
        return false, err
    }

    for _, requiredFile := range requiredFiles {
        if found, err := findFile(ConvertFiles(filesFromDisk), requiredFile, dir); err != nil {
            return false, err
        } else {
            if !found {
                msg := "Can't find " + requiredFile.Name + ": " + FilePath(dir, requiredFile.FileName)
                return false, errors.New(msg)
            } else {
                Logger.Trace("Found:", FilePath(dir, requiredFile.FileName))
            }
        }
    }

    return true, nil
}

func HasRequiredFile(dir string, requiredFile RequiredFile) (bool, error) {
    return HasRequiredFiles(dir, []RequiredFile{requiredFile})
}

func TarDirectory(tw *tar.Writer, dir string, ignoreList []string) error {
    Logger.Trace("Taring: ", dir)

    shouldIgnore := func(file string) bool {
        for _, entry := range ignoreList {
            if entry == file {
                return true
            }
        }
        return false
    }

    var archive = func(files []os.FileInfo) error {
        Logger.Trace("Found files to archive into context: ", len(files))

        for _, file := range files {
            if shouldIgnore(file.Name()) {
                Logger.Debug("Ignoring:", file.Name())
                continue
            }

            fullFilePath := strings.Join([]string{dir, file.Name()}, "/")

            if file.IsDir() {
                TarDirectory(tw, fullFilePath, ignoreList)
                continue
            }

            hdr := &tar.Header{
                Name: fullFilePath,
                Size: file.Size(),
            }
            if err := tw.WriteHeader(hdr); err != nil {
                Logger.Error("Could not write context archive header", err)
                return err
            }

            if bytes, err := ioutil.ReadFile(fullFilePath); err != nil {
                Logger.Error("Could not read context file: ["+fullFilePath+"]", err)
                return err
            } else {
                if _, err := tw.Write(bytes); err != nil {
                    Logger.Error("Could not archive context file: ["+fullFilePath+"]", err)
                    return err
                }
            }

            Logger.Trace("Archived into context the file: [" + fullFilePath + "]")
        }

        Logger.Trace("Successfully archived context", dir)
        return nil
    }

    if filesFromDisk, err := ioutil.ReadDir(dir); err != nil {
        Logger.Error("Error while trying to tar ["+dir+"]", err)
        return err
    } else {
        return archive(filesFromDisk)
    }
}

// listens for ctrl-c from the user
func CaptureUserCancel(containerChannel *ContainerChannel, wg *sync.WaitGroup) {
    syschan := make(chan os.Signal, 1)
    signal.Notify(syschan, os.Interrupt)
    signal.Notify(syschan, syscall.SIGTERM)
    go func() {
        <-syschan
        Logger.Debug("Canceling user operation...")
        if wg != nil {
            *containerChannel <- "canceled"
            wg.Done()
        }
        defer close(*containerChannel)
        return
    }()
}

func FilePath(pathNames ...interface{}) string {
    var paths []string

    for _, value := range pathNames {
        paths = append(paths, value.(string))
    }

    return strings.Join(paths, "/")
}
