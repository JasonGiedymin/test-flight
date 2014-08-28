// Usage
// x,_ := dc.GetImageDetails(fqImageName)
// Logger.Trace( (*x).Id )

// dc.ListContainers(fqImageName)
// time.Sleep(1 * time.Second)

// dc.DeleteContainer(fqImageName)
// time.Sleep(1 * time.Second)

// time.Sleep(1 * time.Second)
// dc.CreateContainer(fqImageName)
// time.Sleep(1 * time.Second)

// dc.DeleteContainer(fqImageName)
// time.Sleep(1 * time.Second)
// dc.DeleteImage(fqImageName)

// Logger.Trace( dc.ListContainers(fqImageName) )

// dc.ShowImages()
// dc.GetImageDetails(fqImageName)

// dc.StopContainer(id)

package lib

import (
    Logger "github.com/JasonGiedymin/test-flight/lib/logging"

    "archive/tar"
    "bufio"
    "bytes"
    "encoding/json"
    "errors"
    "github.com/fsouza/go-dockerclient"
    "github.com/jmcvetta/napping"
    "io"
    "io/ioutil"
    "net/http"
    "os"
    "path/filepath"
    "strconv"
    "strings"
    "sync"
    "text/template"
    "time"
)

type ApiChannel chan *docker.APIEvents
type ContainerChannel chan string

type TemplateVar struct {
    Meta       *ApplicationMeta
    ConfigFile *ConfigFile
    BuildFile  *BuildFile

    TestDir    string
    Owner      string
    ImageName  string
    Tag        string
    From       string
    Version    string
    Requires   []string
    WorkDir    string
    Env        []DockerEnv
    Expose     []int
    Cmd        string
    AddSimple  []string
    AddComplex []DockerAddComplexEntry
    // AddUser           []DockerAddComplexEntry
    RunTests bool
}

// Proxy Client
type DockerApi struct {
    meta             *ApplicationMeta
    configFile       *ConfigFile
    buildFile        *BuildFile
    buildMatrixEntry *BuildMatrixEntry
    client           *docker.Client
}

func NewDockerApi(
    meta *ApplicationMeta,
    configFile *ConfigFile,
    buildFile *BuildFile,
    // buildMatrixEntry *BuildMatrixEntry,
) *DockerApi {
    api := DockerApi{
        meta:       meta,
        configFile: configFile,
        buildFile:  buildFile,
        // buildMatrixEntry: buildMatrixEntry,
    }

    client, err := docker.NewClient(configFile.DockerEndpoint)
    if err != nil {
        Logger.Error("Docker API Client Error:", err)
        os.Exit(ExitCodes["docker_error"])
    }

    Logger.Debug("Docker client created using endpoint:", configFile.DockerEndpoint)
    api.client = client

    return &api
}

func getTemplateDir(configFile *ConfigFile) RequiredFile {
    return RequiredFile{
        Name: "Test-Flight Template Dir", FileName: configFile.AnsibleTemplatesDir, FileType: "d",
    }
}

func (api *DockerApi) RegisterChannel(eventsChannel chan *docker.APIEvents) {
    if err := api.client.AddEventListener(eventsChannel); err != nil {
        Logger.Error(err)
    }

    eventsChannel <- &docker.APIEvents{Status: "Listening in on Docker Events..."}
}

func (api *DockerApi) createTestTemplates(options CommandOptions) error {
    var templateDir = getTemplateDir(api.configFile)

    var inventory = RequiredFile{
        Name: "Test-Flight Test Inventory file", FileName: "inventory", FileType: "f",
    }

    var playbook = RequiredFile{
        Name: "Test-Flight Test Playbook file", FileName: "playbook.yml", FileType: "f",
    }

    // TODO: better to have a type specified that is calculated
    //       much earlier in the process
    var modeDir = func() string {
        if options.SingleFileMode {
            return "filemode"
        } else {
            return "dirmode"
        }
    }()

    // --HERE MUST KNOW ABOUT SUB DIR dirmode/filemode of templates
    // -- needs to work for input and output dir

    // The directory where the templates used to create inventory and playbook
    templates := func() string {
        // could have been simpler but I want to log path at this level
        var templatePath string

        if api.configFile.UseSystemDockerTemplates {
            templatePath = FilePath(api.configFile.TestFlightAssets, "templates", "system")
        } else {
            templatePath = FilePath(api.configFile.TestFlightAssets, "templates", "user")
        }

        // Logger.Debug("Template dir used for ansible tests:", templatePath)
        return templatePath
    }()

    // This is used below in `ExecuteTemplate()`
    templateInputDir := FilePath(templates, modeDir)
    Logger.Debug("Template dir used for generation of ansible files:", templateInputDir)

    // The directory where to put the generated files
    templateOutputDir := FilePath(api.meta.Dir /*api.meta.Dir,*/, templateDir.FileName)
    Logger.Debug("Template dir used to put generated ansible test files:", templateOutputDir)

    createFilesFromTemplate := func(
        templateInputDir string,
        templateOutputDir string,
        requiredFile RequiredFile) error {

        hasFiles, _ := HasRequiredFile(templateOutputDir, requiredFile)

        if hasFiles && api.configFile.OverwriteTemplates || !hasFiles {
            fileToCreate := FilePath(templateOutputDir, requiredFile.FileName)
            var err error
            file, err := os.Create(fileToCreate)
            if err != nil {
                Logger.Error("Error:", err)
                return err
            }

            pattern := filepath.Join(templateInputDir, requiredFile.FileName+"*.tmpl")
            tmpl := template.Must(template.ParseGlob(pattern))

            if err = tmpl.ExecuteTemplate(file, requiredFile.FileName, *api.getTemplateVar()); err != nil {
                Logger.Error("template execution: %s", err)
                return err
            }

            Logger.Debug("Created file from template:", fileToCreate)
        } else if hasFiles && !api.configFile.OverwriteTemplates {
            Logger.Debug(requiredFile.Name, "exists, and system configured to not overwrite.")
        }

        return nil
    }

    // If the test-flight templates dir doesn't exist, create it.
    hasFiles, err := HasRequiredFile(api.meta.Dir, templateDir)

    if !hasFiles { // create if it doesn't exist
        if _, err = CreateFile(&api.meta.Dir, templateDir); err != nil {
            return err
        }
    } else {
        Logger.Debug("Required template files already exist. Overwrite set to:", api.configFile.OverwriteTemplates)
    }

    err = createFilesFromTemplate(templateInputDir, templateOutputDir, inventory)
    if err != nil {
        return err
    }

    err = createFilesFromTemplate(templateInputDir, templateOutputDir, playbook)
    if err != nil {
        return err
    }

    return nil
}

// One big proxy obj to help users. Slowly phase this out.
func (api *DockerApi) getTemplateVar() *TemplateVar {
    return &TemplateVar{
        // Direct:
        Meta:       api.meta,
        ConfigFile: api.configFile,
        BuildFile:  api.buildFile,

        // Helpers for common accessors
        // Keep names simple!
        TestDir:    api.meta.Dir,
        Owner:      api.buildFile.Owner,
        ImageName:  api.buildFile.ImageName,
        Tag:        api.buildFile.Tag,
        From:       api.buildFile.From,
        Requires:   api.buildFile.Requires,
        Version:    api.buildFile.Version,
        WorkDir:    api.buildFile.WorkDir,
        Env:        api.buildFile.Env,
        Expose:     api.buildFile.Expose,
        Cmd:        api.buildFile.Cmd,
        AddSimple:  api.buildFile.Add.Simple,
        AddComplex: api.buildFile.Add.Complex,
        // AddUser:    api.buildFile.Add,
        RunTests: api.buildFile.RunTests,
    }
}

func (api *DockerApi) ShowInfo() {
    Logger.Debug("---------- Test-Flight Docker Info ----------")
    Logger.Debug("Docker Endpoint:", api.configFile.DockerEndpoint)
    Logger.Debug("---------------------------------------------")
}

func (api *DockerApi) ShowImages() error {
    images, _ := api.client.ListImages(true)

    if len(images) <= 0 {
        Logger.Info("No docker images found.")
        return nil
    }

    for _, img := range images {
        Logger.Info("ID: ", img.ID)
        Logger.Info("RepoTags: ", img.RepoTags)
        Logger.Info("Created: ", img.Created)
        Logger.Info("Size: ", img.Size)
        Logger.Info("VirtualSize: ", img.VirtualSize)
        Logger.Info("ParentId: ", img.ParentId)
        Logger.Info("Repository: ", img.Repository)
    }

    return nil
}

// ListContainers that are running
func (api *DockerApi) ListContainers(imageName string) ([]string, error) {
    endpoint := api.configFile.DockerEndpoint
    baseUrl := strings.Join(
        []string{
            endpoint,
            "containers",
            "json",
        },
        "/",
    )
    params := strings.Join(
        []string{
            "all=true",
        },
        "&",
    )
    url := baseUrl + "?" + params
    Logger.Trace("ListContainers Api call:", url)

    resp, _ := http.Get(url)
    defer resp.Body.Close()

    if body, err := ioutil.ReadAll(resp.Body); err != nil {
        Logger.Error("Could not contact docker endpoint:", endpoint)
        return nil, err
    } else {
        switch resp.StatusCode {
        case 200:
            var jsonResult []ApiContainer
            if err := json.Unmarshal(body, &jsonResult); err != nil {
                Logger.Error("Error while trying to marshall result, body:", jsonResult, " - Error:", err)
                return nil, err
            } else {
                var ids []string
                for i := range jsonResult {
                    if jsonResult[i].Image == imageName {
                        ids = append(ids, jsonResult[i].Id)
                    }
                }

                Logger.Debug("Containers found for", imageName, ids)
                return ids, nil
            }
        case 404:
            Logger.Warn("Bad param supplied to API")
        case 500:
            Logger.Error("Error while trying to communicate to docker endpoint:", endpoint)
        }
    }

    return nil, nil
}

// Returns container name deleted, empty string if none, and an
// error (if any)
func (api *DockerApi) DeleteContainer(name string) ([]DeletedContainer, error) {
    var deletedContainers []DeletedContainer
    var returnError error

    endpoint := api.configFile.DockerEndpoint

    delete := func(id string) (*string, error) {
        baseUrl := strings.Join(
            []string{
                endpoint,
                "containers",
                id,
            },
            "/",
        )

        params := strings.Join(
            []string{
                "v=true",
                "force=true",
            },
            "&",
        )
        url := baseUrl + "?" + params

        Logger.Trace("DeleteContainer() Api call:", url)

        req, _ := http.NewRequest("DELETE", url, nil)
        resp, _ := http.DefaultClient.Do(req)
        defer resp.Body.Close()

        if _, err := ioutil.ReadAll(resp.Body); err != nil {
            Logger.Error("Could not contact docker endpoint:", endpoint)
            return nil, err
        } else {
            switch resp.StatusCode {
            case 204:
                Logger.Info("Container Deleted:", name, ", with ID:", id)
                return &name, nil
            case 400:
                msg := "Bad Api param supplied while trying to delete a container, notify developers."
                Logger.Error(msg)
                return nil, errors.New(msg)
            case 404:
                Logger.Warn("Container not found, nothing to delete.")
                return nil, nil
            case 406:
                msg := "Container:" + name + " - (" + id + "), is running, cannot delete."
                Logger.Warn(msg)
                return nil, errors.New(msg)
            case 500:
                msg := "Error while trying to communicate to docker endpoint:" + endpoint
                // Logger.Error(msg)
                return nil, errors.New(msg)
            }

            statusCode := strconv.Itoa(resp.StatusCode)
            msg := "API out of sync, contact developers! Status Code:" + statusCode
            return nil, errors.New(msg)
        }
    }

    if images, err := api.ListContainers(name); err != nil {
        Logger.Error("Could not get image details:", err)
        return nil, err
    } else if images != nil {
        Logger.Info("Found", len(images), "to delete...")
        for _, container := range images {
            Logger.Debug("Trying to delete", container)
            if deleted, err := delete(container); err != nil {
                returnError = err
                msg := "Could not delete container: " + err.Error()
                Logger.Error(msg)
            } else if *deleted != "" {
                deletedContainers = append(deletedContainers, DeletedContainer{*deleted, container})
            }
        }
        return deletedContainers, returnError
    } else {
        Logger.Debug("No containers found for", name)
        return deletedContainers, returnError
    }
}

func (api *DockerApi) BuildImage(
    buffer *bytes.Buffer,
    imageName string,
    watch *ContainerChannel,
    wg *sync.WaitGroup,
) {
    // out := make(ContainerChannel) // channel to send back
    go CaptureUserCancel(watch, wg)

    go func() error {
        endpoint := api.configFile.DockerEndpoint

        baseUrl := strings.Join(
            []string{
                endpoint,
                "build",
            },
            "/",
        )

        params := strings.Join(
            []string{
                "t=" + imageName,
            },
            "&",
        )
        url := baseUrl + "?" + params
        Logger.Trace("BuildImage() - Api call to:", url)

        type BuildStream struct {
            Stream string
        }
        var jsonResult BuildStream

        bytesReader := bytes.NewReader(buffer.Bytes())
        resp, err := http.Post(url, "application/tar", bytesReader)
        defer resp.Body.Close()
        if err != nil {
            Logger.Error("Could not submit request, err:", err)
            return err
        }

        reader := bufio.NewReader(resp.Body)

        for {
            if line, err := reader.ReadBytes('\n'); err != nil {
                if err == io.EOF {
                    break
                } else {
                    msg := "Error reading from stream on building image, error:" + err.Error()
                    Logger.Error(msg)
                    *watch <- "ok"
                    return errors.New(msg)
                }
            } else {
                if err := json.Unmarshal(line, &jsonResult); err != nil {
                    Logger.Error(err)
                } else {
                    Logger.ConsoleChannel(strings.TrimSpace(jsonResult.Stream))
                }
            }
        }

        *watch <- "ok"
        wg.Done()
        return nil
    }()
}

// We only want to attach to containers associated with a particular
// Image name
func (api *DockerApi) Attach(containerId string) ContainerChannel {
    out := make(ContainerChannel) // channel to send back
    go CaptureUserCancel(&out, nil)

    go func() {
        endpoint := api.configFile.DockerEndpoint

        baseUrl := strings.Join(
            []string{
                endpoint,
                "containers",
                containerId,
                "attach",
            },
            "/",
        )

        params := strings.Join(
            []string{
                "stdout=true",
                "stderr=true",
                "stream=true",
            },
            "&",
        )
        url := baseUrl + "?" + params
        Logger.Trace("Attach() - Api call to:", url)

        // jsonResult := ""
        byteData := []byte{}
        bytesReader := bytes.NewReader(byteData)
        resp, err := http.Post(url, "text/json", bytesReader)
        defer resp.Body.Close()
        if err != nil {
            Logger.Error("Could not submit request, err:", err)
        }

        reader := bufio.NewReader(resp.Body)

        for {
            if line, err := reader.ReadBytes('\n'); err != nil {
                if err == io.EOF {
                    break
                } else {
                    msg := "Error reading from stream on attached container, error:" + err.Error()
                    Logger.Error(msg)
                    // close(out)
                    // return error.New(msg)
                }
            } else {
                Logger.ConsoleChannel(string(bytes.TrimSpace(line)[:]))
            }
        }

        close(out)
    }()

    return out
}

func (api *DockerApi) Destroy(fqImageName string) error {
    if running, err := api.ListContainers(fqImageName); err != nil {
        Logger.Trace("Error while trying to get a list of containers for ", fqImageName)
        return err
    } else {
        for _, container := range running {
            if id, err := api.StopContainer(container); err != nil {
                Logger.Error("Could not stop container", id, " associated with", fqImageName)
                return err
            }
        }

        // Once all stopped, delete all
        if _, err := api.DeleteContainer(fqImageName); err != nil {
            msg := "Error occured while trying to delete a container. " +
                " This attempt failed for the following reason:" + err.Error()

            Logger.Error(msg, err)
            return err
        }

        if _, err := api.DeleteImage(fqImageName); err != nil {
            Logger.Error("Could not delete image,", err)
            return err
        }
    }

    return nil
}

func (api *DockerApi) DeleteImage(name string) (string, error) {
    type ContainerStatus map[string]string

    status := func(statusMap ContainerStatus) string {
        possibleStatus := []string{"Untagged", "Deleted"}
        for _, v := range possibleStatus {
            containerId := statusMap[v]
            if containerId != "" {
                return containerId
            }
        }

        Logger.Error("Api out of sync, could not map status, using UNKNOWN")
        return "UNKNOWN"
    }

    endpoint := api.configFile.DockerEndpoint

    delete := func(imageId string) ([]string, error) {
        var deletedContainers []string

        baseUrl := strings.Join(
            []string{
                endpoint,
                "images",
                imageId,
            },
            "/",
        )

        params := strings.Join(
            []string{
                "force=true",
                "noprune=false",
            },
            "&",
        )
        url := baseUrl + "?" + params
        Logger.Trace("DeleteImage() Api call:", url)

        req, _ := http.NewRequest("DELETE", url, nil)
        resp, _ := http.DefaultClient.Do(req)
        defer resp.Body.Close()

        if body, err := ioutil.ReadAll(resp.Body); err != nil {
            Logger.Error("Could not contact docker endpoint:", endpoint)
        } else {
            switch resp.StatusCode {
            case 200:
                var jsonResult []ContainerStatus

                if err := json.Unmarshal(body, &jsonResult); err != nil {
                    Logger.Error(err)
                } else {
                    for _, v := range jsonResult {
                        deletedContainers = append(deletedContainers, status(v))
                    }

                    return deletedContainers, nil
                }
            case 409:
                msg := "Cannot delete image while in use by a container. Delete the container first."
                Logger.Warn(msg)
                return deletedContainers, nil
            case 404:
                msg := "Image not found, cannot delete."
                Logger.Warn(msg)
                return deletedContainers, nil
            case 500:
                msg := "Error while trying to communicate to docker endpoint:"
                Logger.Error(msg)
                return nil, errors.New(msg)
            }

            msg := "API Out of sync, contact developers. Response code: " + strconv.Itoa(resp.StatusCode)
            Logger.Error(msg)
            return nil, errors.New(msg)
        }

        return deletedContainers, nil
    }

    // Need funcs.map(_.apply) where funcs is type [](func, errorString)
    if image, err := api.GetImageDetails(name); err != nil {
        Logger.Error("Could not get image details:", err)
        return "", nil
    } else if image != nil && err == nil {
        if _, err := delete(image.Id); err != nil {
            Logger.Error("Could not delete image: ", name, ", Id:", image.Id, ", Error was:", err)
            return "", err
        }

        return name, nil
    }

    return "", nil
}

// http get using golang packaged lib
func (api *DockerApi) ShowImageGo() (*ApiDockerImage, error) {
    endpoint := api.configFile.DockerEndpoint

    url := strings.Join(
        []string{
            endpoint,
            "images",
            api.buildFile.ImageName + ":" + api.buildFile.Tag,
            "json",
        },
        "/",
    )

    resp, _ := http.Get(url)
    defer resp.Body.Close()

    if body, err := ioutil.ReadAll(resp.Body); err != nil {
        Logger.Error("Could not contact docker endpoint:", endpoint)
        return nil, err
    } else {
        switch resp.StatusCode {
        case 200:
            var jsonResult ApiDockerImage
            if err := json.Unmarshal(body, &jsonResult); err != nil {
                Logger.Error(err)
                return nil, err
            }
            return &jsonResult, nil
        case 404:
            Logger.Debug("Image not found")
        case 500:
            Logger.Error("Error while trying to communicate to docker endpoint:", endpoint)
        }

        return nil, nil
    }
}

func (api *DockerApi) GetImageDetails(fqImageName string) (*ApiDockerImage, error) {
    url := strings.Join(
        []string{
            api.configFile.DockerEndpoint,
            "images",
            fqImageName,
            "json",
        },
        "/",
    )

    // For testing, also need to `import "encoding/json"`
    // result := map[string]json.RawMessage{}
    result := ApiDockerImage{}

    Logger.Trace("GetImageDetails Api call to:", url)
    if resp, err := napping.Get(url, nil, &result, nil); err != nil {
        Logger.Error("Error while getting Image information from docker,", err)
        return nil, err
    } else {
        switch resp.Status() {
        case 200:
            Logger.Trace(result)
            return &result, nil
        case 404:
            Logger.Debug("Image not found")
            return nil, nil
        }

        return nil, nil
    }
}

// CreateDockerImage creates a docker image by first creating the Dockerfile
// in memory and then tars it, prior to sending it along to the docker
// endpoint.
// Note: once the image creation starts the only way to know if it succeeded
//       is to query for it again.
func (api *DockerApi) CreateDockerImage(fqImageName string, options *CommandOptions) (string, error) {
    Logger.Debug("Creating Docker image by attempting to build Dockerfile: " + fqImageName)

    tmplName := "Dockerfile" // file ext of `.tmpl` is implicit, see below

    var requiredDockerFile = RequiredFile{
        Name: "Test-Flight Dockerfile", FileName: tmplName, FileType: "f",
    }

    modeDir := func() string {
        if options.SingleFileMode {
            return "filemode"
        } else {
            return "dirmode"
        }
    }()

    dockerfileBuffer := bytes.NewBuffer(nil)
    tarbuf := bytes.NewBuffer(nil)

    // Dockerfile mem
    dockerfile := bufio.NewWriter(dockerfileBuffer)
    // Or create file
    templateOutputDir := getTemplateDir(api.configFile)
    Logger.Trace(templateOutputDir, requiredDockerFile.FileName)
    fileToCreate := FilePath(templateOutputDir.FileName, requiredDockerFile.FileName)
    Logger.Debug("Trying to build:", fileToCreate)
    dockerfileOut, err := os.Create(fileToCreate)
    if err != nil {
        Logger.Error("Could not write Dockerfile", err)
        return "", err
    }
    defer dockerfileOut.Close()

    // The directory where the templates used to create inventory and playbook
    templates := func() string {
        // could have been simpler but I want to log path at this level
        var templatePath string

        if api.configFile.UseSystemDockerTemplates {
            templatePath = FilePath(api.configFile.TestFlightAssets, "templates", "system")
        } else {
            templatePath = FilePath(api.configFile.TestFlightAssets, "templates", "user")
        }

        Logger.Debug("Using template dir:", templatePath)
        return templatePath
    }()

    // This is used below in `ExecuteTemplate()`
    templateInputDir := FilePath(templates, modeDir)

    pattern := filepath.Join(templateInputDir, requiredDockerFile.FileName+"*.tmpl")
    tmpl := template.Must(template.ParseGlob(pattern))

    // In Mem
    if err := tmpl.ExecuteTemplate(dockerfile, requiredDockerFile.FileName, *api.getTemplateVar()); err != nil {
        Logger.Error("template execution: %s", err)
        return "", err
    }

    // File
    if err := tmpl.ExecuteTemplate(dockerfileOut, requiredDockerFile.FileName, *api.getTemplateVar()); err != nil {
        Logger.Error("template execution: %s", err)
        return "", err
    }

    // For in-mem flush buffer here
    dockerfile.Flush()

    // Logger.Trace("Dockerfile buffer len", dockerfileBuffer.Len())
    Logger.Trace("Dockerfile:", dockerfileBuffer.String())
    // Logger.Info("Created Dockerfile: " + fqImageName)

    currTime := time.Now()

    // Add Dockerfile to archive, break out
    tr := tar.NewWriter(tarbuf)
    tr.WriteHeader(&tar.Header{
        Name:       "Dockerfile",
        Size:       int64(dockerfileBuffer.Len()),
        ModTime:    currTime,
        AccessTime: currTime,
        ChangeTime: currTime,
    })
    tr.Write(dockerfileBuffer.Bytes())

    // Add Context to archive
    // tar test directory
    TarDirectory(tr, api.meta.Dir, api.buildFile.Ignore)
    tr.Close()

    var wg sync.WaitGroup // used for channels
    wg.Add(1)

    Logger.Trace("Building image...")
    buildChannel := make(ContainerChannel)
    // go watchContainerOn(buildChannel, &wg)
    api.BuildImage(tarbuf, fqImageName, &buildChannel, &wg)
    // go CaptureUserCancel(&buildChannel)

    select {
    case value := <-buildChannel:
        switch value {
        case "ok":
            failMsg := "Docker Image [" + fqImageName + "] failed to build."
            // We got an ok but don't know if it is legit. Possible to get no error
            // but image will be nil (successful api call but image doesn't exist)
            if image, err := api.GetImageDetails(fqImageName); err != nil || image == nil {
                return "", errors.New(failMsg)
            } else {
                Logger.Info("Successfully built Docker image: " + fqImageName)
                return image.Id, nil
            }
        case "canceled":
            msg := "User Canceled docker creation."
            Logger.Warn(msg, "User must manually stop last running container.")
            return fqImageName, errors.New(msg)
        }
    }

    wg.Wait()

    Logger.Info("Successfully built Docker image: " + fqImageName)
    return fqImageName, nil
}

func (api *DockerApi) CreateContainer(fqImageName string) (*ApiPostResponse, error) {
    if _, err := api.DeleteContainer(fqImageName); err != nil {
        msg := "Error occured while trying to delete a container. " +
            " This attempt failed for the following reason:" + err.Error()

        Logger.Error(msg, err)
        return nil, err
    }

    endpoint := api.configFile.DockerEndpoint

    postBody := ApiPostRequest{
        Image:        fqImageName,
        OpenStdin:    true,
        AttachStdin:  false,
        AttachStdout: true,
        AttachStderr: true,
        Cmd:          api.buildFile.LaunchCmd,
    }

    url := FilePath(endpoint, "containers", "create")
    Logger.Trace("CreateContainer() - Api call to:", url)

    jsonResult := ApiPostResponse{}
    bytesReader, _ := postBody.Bytes()
    resp, _ := http.Post(url, "application/json", bytes.NewReader(bytesReader))
    defer resp.Body.Close()

    if body, err := ioutil.ReadAll(resp.Body); err != nil {
        Logger.Error("Could not contact docker endpoint:", endpoint)
        return nil, err
    } else {
        msg := "Unexpected response code: " + string(resp.StatusCode)

        switch resp.StatusCode {
        case 201:
            if err := json.Unmarshal(body, &jsonResult); err != nil {
                Logger.Error(err)
                return nil, err
            }
            Logger.Info("Created container:", fqImageName)
            return &jsonResult, nil
        case 404:
            msg = "No such container"
            Logger.Warn(msg)
        case 406:
            msg = "Impossible to attach (container not running)"
            Logger.Warn(msg)
        case 500:
            // noticed that when docker endpoint fails, it fails with just
            // status 500, no message back. Logs do show error though somewhat
            // slim details. i.e No command specified where command is really CMD :-(
            msg = "Server and/or API error while trying to communicate with docker " +
                "endpoint: " + endpoint + ". Could be malformed Dockerfile (template). " +
                "Check your remote docker endpoint logs."
            Logger.Error(msg)
        }

        return nil, errors.New(msg)
    }
}

func (api *DockerApi) StopContainer(containerId string) (string, error) {
    endpoint := api.configFile.DockerEndpoint

    url := strings.Join(
        []string{
            endpoint,
            "containers",
            containerId,
            "stop",
        },
        "/",
    )
    Logger.Trace("StopContainer() - Api call to:", url)

    resp, _ := http.Post(url, "text/json", nil)
    defer resp.Body.Close()

    if _, err := ioutil.ReadAll(resp.Body); err != nil {
        Logger.Error("Could not contact docker endpoint:", endpoint)
        return "", err
    } else {
        switch resp.StatusCode {
        case 204:
            Logger.Info("Stopped container:", containerId)
            return containerId, nil
        case 404:
            msg := "No such container"
            Logger.Warn(msg)
            return "", nil
        case 500:
            msg := "Error while trying to communicate to docker endpoint: " + endpoint
            Logger.Error(msg)
            return "", nil
        }

        msg := "Unexpected response code: " + string(resp.StatusCode)
        Logger.Error(msg)
        return "", errors.New(msg)
    }
}

func (api *DockerApi) StartContainer(id string) (*string, error) {
    endpoint := api.configFile.DockerEndpoint

    postBody := DockerHostConfig{}

    url := strings.Join(
        []string{
            endpoint,
            "containers",
            id,
            "start",
        },
        "/",
    )
    Logger.Trace("StartContainer() - Api call to:", url)

    var jsonResult string
    bytesReader, _ := postBody.Bytes()

    resp, _ := http.Post(url, "text/json", bytes.NewReader(bytesReader))
    defer resp.Body.Close()

    if _, err := ioutil.ReadAll(resp.Body); err != nil {
        Logger.Error("Could not contact docker endpoint:", endpoint)
        return nil, err
    } else {
        switch resp.StatusCode {
        case 204:
            Logger.Info("Started container:", id)
            return &jsonResult, nil
        case 404:
            Logger.Warn("No such container")
        case 500:
            Logger.Error("Error while trying to communicate to docker endpoint:", endpoint)
        }

        msg := "Unexpected response code: " + string(resp.StatusCode)
        Logger.Error(msg)
        return nil, errors.New(msg)
    }
}
