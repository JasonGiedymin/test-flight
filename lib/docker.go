package lib

import (
  Logger "./logging"
  "./types"
  "archive/tar"
  "bufio"
  "bytes"
  "github.com/fsouza/go-dockerclient"
  "os"
  "path/filepath"
  "strings"
  "text/template"
  "time"
)

type TemplateVar struct {
  Meta       *types.ApplicationMeta
  ConfigFile *types.ConfigFile
  BuildFile  *types.BuildFile

  TestDir           string
  Owner             string
  ImageName         string
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
  watch      chan *docker.APIEvents
}

func NewDockerApi(meta *types.ApplicationMeta, configFile *types.ConfigFile, buildFile *types.BuildFile) *DockerApi {
  api := DockerApi{
    meta: meta,
    configFile: configFile,
    buildFile: buildFile,
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

func (api *DockerApi) WatchApiEvents() {
  var listen = func() {
    for {
      event := <- api.watch
      Logger.Info(event)
    }
  }

  go listen()

  if err:= api.client.AddEventListener(api.watch); err != nil {
    Logger.Error(err)
    // return nil
  }

  api.watch <- &docker.APIEvents{}
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

func (api *DockerApi) CreateDocker() error {
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

  // --- test
  // Logger.Trace("Dockerfile:", dockerfile)
  // Logger.Trace("! dockerfile buffered:", dockerfile.Buffered())
  // Logger.Trace("! dockerfile buffer available:", dockerfile.Available())
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

  // tr.WriteHeader(&tar.Header{
  //   Name:       api.meta.Dir,
  //   Size:       int64(dockerfileBuffer.Len()),
  //   ModTime:    currTime,
  //   AccessTime: currTime,
  //   ChangeTime: currTime,
  // }) //new header
  // tr.Write(dockerfileBuffer.Bytes())
  tr.Close()

  Logger.Trace("Dockerfile buffer len", dockerfileBuffer.Len())
  Logger.Trace("Dockerfile:", dockerfileBuffer.String())

  opts := docker.BuildImageOptions{
    Name:         api.buildFile.ImageName,
    InputStream:  tarbuf,
    OutputStream: outputbuf,
  }
  if err := api.client.BuildImage(opts); err != nil {
    Logger.Error(err)
    return err
  }
  // --- test

  return nil
}
