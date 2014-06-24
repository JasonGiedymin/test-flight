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
  Logger "./logging"
  "./types"
  "archive/tar"
  "bufio"
  "bytes"
  "encoding/json"
  "errors"
  "github.com/fsouza/go-dockerclient"
  "github.com/jmcvetta/napping"
  "io/ioutil"
  "net/http"
  "os"
  "path/filepath"
  "strconv"
  "strings"
  "text/template"
  "time"
)

type ApiChannel chan *docker.APIEvents

type TemplateVar struct {
  Meta       *types.ApplicationMeta
  ConfigFile *types.ConfigFile
  BuildFile  *types.BuildFile

  TestDir           string
  Owner             string
  ImageName         string
  Tag               string
  From              string
  Version           string
  RequiresDocker    string
  RequiresDockerUrl string
  WorkDir           string
  Env               map[string]string
  Expose            []int
  Cmd               string
  AddSimple         []string
  AddComplex        []types.DockerAddComplexEntry
  AddUser           []types.DockerAddComplexEntry
  RunTests          bool
}

// Proxy Client
type DockerApi struct {
  meta       *types.ApplicationMeta
  configFile *types.ConfigFile
  buildFile  *types.BuildFile
  client     *docker.Client
}

func NewDockerApi(meta *types.ApplicationMeta, configFile *types.ConfigFile, buildFile *types.BuildFile) *DockerApi {
  api := DockerApi{
    meta:       meta,
    configFile: configFile,
    buildFile:  buildFile,
  }

  client, err := docker.NewClient(configFile.DockerEndpoint)
  if err != nil {
    Logger.Error("Docker API Client Error:", err)
    os.Exit(ExitCodes["docker_error"])
  }

  Logger.Info("Docker endpoint:", configFile.DockerEndpoint)
  api.client = client

  return &api
}

func getTemplateDir(configFile *types.ConfigFile) types.RequiredFile {
  return types.RequiredFile{
    Name: "Test-Flight Template Dir", FileName: configFile.TemplateDir, FileType: "d",
  }
}

func (api *DockerApi) RegisterChannel(eventsChannel chan *docker.APIEvents) {
  if err := api.client.AddEventListener(eventsChannel); err != nil {
    Logger.Error(err)
  }

  eventsChannel <- &docker.APIEvents{Status: "Listening in on Docker Events..."}
}

func (api *DockerApi) createTestTemplates() error {
  var templateDir = getTemplateDir(api.configFile)

  var inventory = types.RequiredFile{
    Name: "Test-Flight Test Inventory file", FileName: "inventory", FileType: "f",
  }

  var playbook = types.RequiredFile{
    Name: "Test-Flight Test Playbook file", FileName: "playbook.yml", FileType: "f",
  }

  templateOutputDir := strings.Join([]string{api.meta.Pwd, api.meta.Dir, templateDir.FileName}, "/")
  templateInputDir := api.meta.Pwd + "/templates/"

  createFilesFromTemplate := func(
    templateInputDir string,
    templateOutputDir string,
    requiredFile types.RequiredFile) error {
    // check that inventory files exist
    if hasFiles, err := HasRequiredFile(templateOutputDir, requiredFile); err != nil {
      Logger.Error("Error: ", err)
      return err
    } else if hasFiles {
      // Inventory
      fileToCreate := strings.Join([]string{templateOutputDir, requiredFile.FileName}, "/")
      file, _ := os.Create(fileToCreate)

      pattern := filepath.Join(templateInputDir, requiredFile.FileName+"*.tmpl")
      tmpl := template.Must(template.ParseGlob(pattern))

      if err = tmpl.ExecuteTemplate(file, requiredFile.FileName, *api.getTemplateVar()); err != nil {
        Logger.Error("template execution: %s", err)
        return err
      }

      Logger.Debug("Created file from template:", fileToCreate)
    }

    // check that dir exists
    // TODO: major cleanup here, need another pass
    if hasFiles, err := HasRequiredFile(api.meta.Dir, templateDir); err != nil {
      return err
    } else if !hasFiles { // create it doesn't
      if _, err = CreateFile(&api.meta.Dir, templateDir); err != nil {
        return err
      }
    }
    return nil
  }

  // if api.build
  _ = createFilesFromTemplate(templateInputDir, templateOutputDir, inventory)
  _ = createFilesFromTemplate(templateInputDir, templateOutputDir, playbook)

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
    TestDir:           api.meta.Dir,
    Owner:             api.buildFile.Owner,
    ImageName:         api.buildFile.ImageName,
    Tag:               api.buildFile.Tag,
    From:              api.buildFile.From,
    Version:           api.buildFile.Version,
    RequiresDocker:    api.buildFile.RequiresDocker,
    RequiresDockerUrl: api.buildFile.RequiresDockerUrl,
    WorkDir:           api.configFile.WorkDir,
    Env:               api.buildFile.Env,
    Expose:            api.buildFile.Expose,
    Cmd:               api.buildFile.Cmd,
    AddSimple:         api.configFile.DockerAdd.Simple,
    AddComplex:        api.configFile.DockerAdd.Complex,
    AddUser:           api.buildFile.Add,
    RunTests:          api.buildFile.RunTests,
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
      "all=0",
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
      var jsonResult []types.ApiContainer
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
func (api *DockerApi) DeleteContainer(name string) ([]types.DeletedContainer, error) {
  var deletedContainers []types.DeletedContainer
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
        deletedContainers = append(deletedContainers, types.DeletedContainer{*deleted, container})
      }
    }
    return deletedContainers, returnError
  } else {
    Logger.Debug("Nothing to delete, no containers found for", name)
    return deletedContainers, returnError
  }
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

  logStatus := func(statusMap ContainerStatus) string {
    possibleStatus := []string{"Untagged", "Deleted"}
    for _, v := range possibleStatus {
      containerId := statusMap[v]
      if containerId != "" {
        Logger.Debug(v, "container Id:", containerId)
        return containerId
      } else {
        Logger.Warn("Unknown status found:", statusMap)
      }

      return containerId
    }

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
          // return nil, err
        } else {
          for _, v := range jsonResult {
            deletedContainers = append(deletedContainers, logStatus(v))
          }

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

      msg := "API Out of sync, contact developers"
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
    Logger.Warn("-->", image.Id)
    if _, err := delete(image.Id); err != nil {
      Logger.Error("Could not delete image: ", name, ", Id:", image.Id, ", Error was:", err)
      return "", err
    }

    return name, nil
  }

  return "", nil
}

// http get using golang packaged lib
func (api *DockerApi) ShowImageGo() (*types.ApiDockerImage, error) {
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
      var jsonResult types.ApiDockerImage
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

func (api *DockerApi) GetImageDetails(fqImageName string) (*types.ApiDockerImage, error) {
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
  result := types.ApiDockerImage{}

  Logger.Trace("GetImageDetails Api call to:", url)
  if resp, err := napping.Get(url, nil, &result, nil); err != nil {
    Logger.Error("Error while getting Image information from docker,", err)
    return nil, err
  } else {
    switch resp.Status() {
    case 200:
      return &result, nil
    case 404:
      Logger.Debug("Image not found")
    }

    return &types.ApiDockerImage{}, nil
  }
}

// CreateDockerImage creates a docker image by first creating the Dockerfile
// in memory and then tars it, prior to sending it along to the docker
// endpoint.
func (api *DockerApi) CreateDockerImage(fqImageName string) (string, error) {
  Logger.Debug("Creating Docker image by attempting to build Dockerfile: " + fqImageName)

  dockerfileBuffer := bytes.NewBuffer(nil)
  tarbuf := bytes.NewBuffer(nil)
  outputbuf := bytes.NewBuffer(nil)
  dockerfile := bufio.NewWriter(dockerfileBuffer)

  var requiredDockerFile = types.RequiredFile{
    Name: "Test-Flight Dockerfile", FileName: "Dockerfile", FileType: "f",
  }

  // var templateDir = getTemplateDir(api.configFile) // might need later
  // templateOutputDir := strings.Join([]string{api.meta.Pwd, api.meta.Dir, templateDir.FileName}, "/")
  templateInputDir := api.meta.Pwd + "/templates/"

  pattern := filepath.Join(templateInputDir, requiredDockerFile.FileName+"*.tmpl")
  tmpl := template.Must(template.ParseGlob(pattern))

  if err := tmpl.ExecuteTemplate(dockerfile, requiredDockerFile.FileName, *api.getTemplateVar()); err != nil {
    Logger.Error("template execution: %s", err)
    return "", err
  }

  dockerfile.Flush()

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
  TarDirectory(tr, api.meta.Dir)
  tr.Close()

  Logger.Trace("Dockerfile buffer len", dockerfileBuffer.Len())
  Logger.Trace("Dockerfile:", dockerfileBuffer.String())
  Logger.Info("Created Dockerfile: " + fqImageName)

  opts := docker.BuildImageOptions{
    Name:         fqImageName,
    InputStream:  tarbuf,
    OutputStream: outputbuf,
  }

  if err := api.client.BuildImage(opts); err != nil {
    Logger.Error("Error while building Docker image: "+fqImageName, err)
    return "", err
  }

  Logger.Info("Successfully built Docker image: " + fqImageName)
  return fqImageName, nil
}

func (api *DockerApi) CreateContainer(fqImageName string) (*types.ApiPostResponse, error) {
  if _, err := api.DeleteContainer(fqImageName); err != nil {
    msg := "Error occured while trying to delete a container. " +
      " This attempt failed for the following reason:" + err.Error()

    Logger.Error(msg, err)
    return nil, err
  }

  endpoint := api.configFile.DockerEndpoint

  postBody := types.ApiPostRequest{
    Image:       fqImageName,
    OpenStdin:   true,
    AttachStdin: true,
  }

  url := strings.Join(
    []string{
      endpoint,
      "containers",
      "create",
    },
    "/",
  )
  Logger.Trace("CreateContainer() - Api call to:", url)

  jsonResult := types.ApiPostResponse{}
  bytesReader, _ := postBody.Bytes()
  resp, _ := http.Post(url, "text/json", bytes.NewReader(bytesReader))
  defer resp.Body.Close()

  if body, err := ioutil.ReadAll(resp.Body); err != nil {
    Logger.Error("Could not contact docker endpoint:", endpoint)
    return nil, err
  } else {
    switch resp.StatusCode {
    case 201:
      if err := json.Unmarshal(body, &jsonResult); err != nil {
        Logger.Error(err)
        return nil, err
      }
      Logger.Info("Created container:", fqImageName)
      return &jsonResult, nil
    case 404:
      Logger.Warn("No such container")
    case 406:
      Logger.Warn("Impossible to attach (container not running)")
    case 500:
      Logger.Error("Error while trying to communicate to docker endpoint:", endpoint)
    }

    msg := "Unexpected response code: " + string(resp.StatusCode)
    Logger.Error(msg)
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

  postBody := types.DockerHostConfig{}

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

  // jsonResult := types.ApiPostResponse{}
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
      // if err := json.Unmarshal(body, &jsonResult); err != nil {
      //   Logger.Error(err)
      //   return nil, err
      // }
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

func (api *DockerApi) CreateContainer2() error {
  opts := docker.CreateContainerOptions{
    Name: api.buildFile.ImageName,
    Config: &docker.Config{
      Image:       "test-docker-name1",
      OpenStdin:   true,
      AttachStdin: true,
    },
  }

  if container, err := api.client.CreateContainer(opts); err != nil {
    Logger.Error(err)
    return err
  } else {
    Logger.Info("Container created: ", container.ID[:12])
  }
  return nil
}
