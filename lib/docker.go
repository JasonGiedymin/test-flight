package lib

import (
  Logger "./logging"
  "./types"
  "archive/tar"
  "bufio"
  "bytes"
  "github.com/fsouza/go-dockerclient"
  "github.com/jmcvetta/napping"
  "net/http"
  "os"
  "path/filepath"
  "strings"
  "text/template"
  "time"
  "encoding/json"
  "io/ioutil"
  "errors"
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
    if hasFiles, err := HasRequiredFile(&templateOutputDir, requiredFile); err != nil {
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
    if hasFiles, err := HasRequiredFile(&api.meta.Dir, templateDir); err != nil {
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

func (api *DockerApi) ShowImages() {
  images, _ := api.client.ListImages(true)

  if len(images) <= 0 {
    Logger.Info("No docker images found.")
    return
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
}

func (api *DockerApi) ListContainers(imageName string) ([]string, error){
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
      "all=1",
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
        Logger.Error(err)
        return nil, err
      } else {
        var ids []string
        for i := range jsonResult {
          if jsonResult[i].Image == imageName {
            ids = append(ids, jsonResult[i].Id)
          }
        }

        Logger.Trace("Containers found for", imageName, ids)
        return ids, nil
      }
    case 404:
      Logger.Trace("Bad param supplied to API")
    case 500:
      Logger.Error("Error while trying to communicate to docker endpoint:", endpoint)
    }
  }

  return nil, nil
}

// Returns container name deleted, empty string if none, and an
// error (if any)
func (api *DockerApi) DeleteContainer(name string) ([]string, error) {
  var deletedContainers []string
  endpoint := api.configFile.DockerEndpoint

  delete := func(id string) (*string, error) {
    url := strings.Join(
      []string{
        endpoint,
        "containers",
        id,
      },
      "/",
    )
    Logger.Trace("DeleteContainer Api call:", url)

    req, _ := http.NewRequest("DELETE",url, nil)
    resp, _ := http.DefaultClient.Do(req)
    defer resp.Body.Close()

    if _, err := ioutil.ReadAll(resp.Body); err != nil {
      Logger.Error("Could not contact docker endpoint:", endpoint)
      return nil, err
    } else {
      switch resp.StatusCode {
      case 204:
        Logger.Info("Container Deleted:", name)
        return &name, nil
      case 400:
        msg := "Bad Api param supplied while trying to delete a container, notify developers."
        Logger.Error(msg)
        return nil, errors.New(msg)
      case 404:
        Logger.Error("Container not found, nothing to delete.")
        var emptyString string
        return &emptyString, nil
      case 500:
        msg := "Error while trying to communicate to docker endpoint:" + endpoint
        Logger.Error(msg)
        return nil, errors.New(msg)
      }

      return nil, errors.New("API out of sync, contact developers!")
    }
  }

  if images, err := api.ListContainers(name); err != nil {
    Logger.Error("Could not get image details:", err)
    return nil, err
  } else if images != nil {
    for _, container := range images {
      if deleted, _ := delete(container); *deleted != "" {
        deletedContainers = append(deletedContainers, *deleted)
      }
    }
    return deletedContainers, nil
  } else {
    Logger.Trace("Nothing to delete, no containers found for", name)
    return deletedContainers, nil
  }
}

func (api *DockerApi) DeleteImage(name string) {
  endpoint := api.configFile.DockerEndpoint

  delete := func(imageId string) {
    url := strings.Join(
      []string{
        endpoint,
        "images",
        name,
      },
      "/",
    )

    req, _ := http.NewRequest("DELETE",url, nil)
    resp, _ := http.DefaultClient.Do(req)
    defer resp.Body.Close()

    if body, err := ioutil.ReadAll(resp.Body); err != nil {
      Logger.Error("Could not contact docker endpoint:", endpoint)
      // return nil, err
    } else {
      switch resp.StatusCode {
      case 200:
        // var deleted []string
        var jsonResult []map[string]string
        if err := json.Unmarshal(body, &jsonResult); err != nil {
          Logger.Error(err)
          // return nil, err
        } else {
          for k, v := range jsonResult {
            Logger.Trace(v)
          }
          Logger.Info(jsonResult)
        }
      case 409:
        Logger.Error("Cannot delete image while in use by a container. Delete the container first.")
      case 404:
        Logger.Error("Image not found, cannot delete.")
      case 500:
        Logger.Error("Error while trying to communicate to docker endpoint:", endpoint)
      }
    }
  }

  if image, err := api.GetImageDetails(name); err != nil {
    Logger.Error("Could not get image details:", err)
  } else if image != nil {
    delete(name)
  }
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
      Logger.Trace("Image not found")
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
      Logger.Trace("Image not found")
    }

    return &types.ApiDockerImage{}, nil
  }
}

func (api *DockerApi) CreateDockerImage() error {
  Logger.Info("Attempting to build Dockerfile: " + api.buildFile.ImageName)

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
    return err
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
  Logger.Info("Created Dockerfile: " + api.buildFile.ImageName)

  opts := docker.BuildImageOptions{
    Name:         api.buildFile.ImageName,
    InputStream:  tarbuf,
    OutputStream: outputbuf,
  }

  if err := api.client.BuildImage(opts); err != nil {
    Logger.Error("Error while building Docker image: "+api.buildFile.ImageName, err)
    return err
  }

  Logger.Info("Successfully built Docker image: " + api.buildFile.ImageName)

  return nil
}

func (api *DockerApi) CreateContainer() error {
  opts := docker.CreateContainerOptions{
    Name: api.buildFile.ImageName,
    Config: &docker.Config{
      Image:       "test-docker-name1",
      OpenStdin:   true,
      AttachStdin: true,
      Cmd:         []string{"bash"},
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
