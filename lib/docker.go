package lib

import (
  Logger "./logging"
  "./types"
  "github.com/fsouza/go-dockerclient"
  "os"
  "path/filepath"
  "text/template"
  "strings"
)

type TemplateVar struct {
  Meta       *types.ApplicationMeta
  ConfigFile *types.ConfigFile
  BuildFile  *types.BuildFile

  Owner             string
  ImageName         string
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
  api := DockerApi{meta: meta, configFile: configFile, buildFile: buildFile}
  client, err := docker.NewClient(configFile.DockerEndpoint)
  if err != nil {
    Logger.Error("Docker API Client Error:", err)
    os.Exit(ExitCodes["docker_error"])
  }

  api.client = client
  return &api
}

func (api *DockerApi) createTestTemplates() error {
  var templateDir = RequiredFile{
    name: "Test-Flight Template Dir", fileName: api.configFile.TemplateDir, FileType: "d",
  }

  var inventory = RequiredFile{
    name: "Test-Flight Test Inventory file", fileName: "inventory", FileType: "f",
  }

  var playbook = RequiredFile{
    name: "Test-Flight Test Playbook file", fileName: "playbook.yml", FileType: "f",
  }

  templateOutputDir := strings.Join([]string{api.meta.Pwd, api.meta.Dir, templateDir.fileName}, "/")
  templateInputDir := api.meta.Pwd + "/templates/"
  // var err error

  createFilesFromTemplate := func(
    templateInputDir string,
    templateOutputDir string,
    requiredFile RequiredFile) error {
    // check that inventory files exist
    if hasFiles, err := HasRequiredFile(&templateOutputDir, requiredFile); err != nil {
      Logger.Error("Error: ", err)
      return err
    } else if hasFiles {
      // Inventory
      fileToCreate := strings.Join([]string{templateOutputDir,requiredFile.fileName}, "/")
      file, _ := os.Create(fileToCreate)

      pattern := filepath.Join(templateInputDir, requiredFile.fileName + "*.tmpl")
      tmpl := template.Must(template.ParseGlob(pattern))

      if err = tmpl.ExecuteTemplate(file, requiredFile.fileName, *api.getTemplateVar()); err != nil {
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
    Owner:             api.buildFile.Owner,
    ImageName:         api.buildFile.ImageName,
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

func (api *DockerApi) CreateTemplate() {
  var (
    pattern         string
    tmpl            *template.Template
    pwd             string
    err             error
    baseTemplateDir string
    testTemplateDir string
  )

  pwd, err = os.Getwd()
  if err != nil {
    // return nil, err
  }

  // baseTemplateDir = api.meta.ExecPath + "./templates/"
  testTemplateDir = pwd + "/" + api.configFile.TemplateDir + "/"
  baseTemplateDir = pwd + "/templates/"

  Logger.Trace("Base Template Dir:", baseTemplateDir)
  Logger.Trace("Test Template Dir:", testTemplateDir)

  // Dockerfile
  pattern = filepath.Join(baseTemplateDir, "Dockerfile*.tmpl")
  tmpl = template.Must(template.ParseGlob(pattern))

  if err = tmpl.ExecuteTemplate(os.Stdout, "Dockerfile", *api.getTemplateVar()); err != nil {
    Logger.Error("template execution: %s", err)
  }



  // Playbook
  pattern = filepath.Join(baseTemplateDir, "playbook.yml*.tmpl")
  tmpl = template.Must(template.ParseGlob(pattern))
  // tmpl := template.Must(template.ParseFiles()

  if err = tmpl.ExecuteTemplate(os.Stdout, "playbook", *api.getTemplateVar()); err != nil {
    Logger.Error("template execution: %s", err)
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

func (api *DockerApi) createDockerFile() string {
  dockerFile := `
  # Dockerfile
  # ----------
  #
  `

  return dockerFile
}

func (api *DockerApi) CreateDocker() {
}
