Test-Flight Config File
-----------------------

The config file `test-flight-config.json` should exist either in a user's home
or the local pwd where the flight command is executed.

## Config File Fields

  - `dockerEndPoint`: _Required_, Default=`to what is set in this file`
  - `workDir`: __string__, _Required_, Default=`to what is set in this file`
  - `dockerAdd`: __object{system:array[string], user:array[object{string:string}]},
                 _System Object Required_, Default=`to what is set in this file`

## Example Default Config File

    {
      "dockerEndPoint": "http://localhost:4243",

      "workDir": "/tmp/build",

      "dockerAdd": {
        "system": [
          "docker",
          "artifacts",
          "tasks",
          "tests",
          "vars",
          "templates",
          "handlers",
        ],
        "user": [
          {"name": "temp", "location": "{{.WorkDir}}/temp"},
        ]
      }

    }
