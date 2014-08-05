package build

import (
    "github.com/JasonGiedymin/test-flight/lib"
    "reflect"
    "sort"
    "testing"
)

var keyTestData = []struct {
    entry    BuildMatrixEntry
    expected string
}{
    {
        BuildMatrixEntry{
            Language: "c++",
            Version:  "4.6.1",
            OS:       "ubuntu",
            Env:      lib.DockerEnv{"TMPDIR", "/tmp/trash"},
        },
        "(c++,ubuntu,4.6.1,TMPDIR=/tmp/trash)",
    },
    {
        BuildMatrixEntry{
            Language: "c++",
            Version:  "4.6.1",
            OS:       "ubuntu",
            Env:      lib.DockerEnv{"TMPDIR", "/tmp/trash"},
            From:     "override-ansible-ubuntu",
        },
        "(override-ansible-ubuntu,TMPDIR=/tmp/trash)",
    },
}

var testData = []struct {
    Vectors            BuildMatrixVectors
    ExpectedMatrixKeys []string
    SampleBuildFile    lib.BuildFile
}{
    {
        BuildMatrixVectors{
            Language: "golang",
            OS:       []string{"ubuntu"},
            Version: []string{
                "1",
                "2",
                "3",
            },
            Env: []lib.DockerEnv{
                lib.DockerEnv{"trash", "can"},
                lib.DockerEnv{"can", "fruit"},
                lib.DockerEnv{"car", "go"},
            },
        },
        []string{
            "(golang,ubuntu,1,trash=can)",
            "(golang,ubuntu,1,can=fruit)",
            "(golang,ubuntu,1,car=go)",
            "(golang,ubuntu,2,trash=can)",
            "(golang,ubuntu,2,can=fruit)",
            "(golang,ubuntu,2,car=go)",
            "(golang,ubuntu,3,trash=can)",
            "(golang,ubuntu,3,can=fruit)",
            "(golang,ubuntu,3,car=go)",
        },
        *lib.NewBuildFile(),
    },
    {
        BuildMatrixVectors{
            Language: "c++",
            OS:       []string{"ubuntu"},
            Version: []string{
                "4.6.0",
                "4.6.1",
                "4.6.2",
            },
            Env: []lib.DockerEnv{
                lib.DockerEnv{"TMPDIR", "/tmp/trash"},
                lib.DockerEnv{"LIBRARY_PATH", "/usr/lib/my.lib.path/"},
            },
        },
        []string{
            "(c++,ubuntu,4.6.0,TMPDIR=/tmp/trash)",
            "(c++,ubuntu,4.6.0,LIBRARY_PATH=/usr/lib/my.lib.path/)",
            "(c++,ubuntu,4.6.1,TMPDIR=/tmp/trash)",
            "(c++,ubuntu,4.6.1,LIBRARY_PATH=/usr/lib/my.lib.path/)",
            "(c++,ubuntu,4.6.2,TMPDIR=/tmp/trash)",
            "(c++,ubuntu,4.6.2,LIBRARY_PATH=/usr/lib/my.lib.path/)",
        },
        *lib.NewBuildFile(),
    },
    {
        BuildMatrixVectors{
            Language: "c++",
            OS:       []string{"ubuntu"},
            Version: []string{
                "4.6.0",
                "4.6.1",
                "4.6.2",
            },
            Env: []lib.DockerEnv{
                lib.DockerEnv{"TMPDIR", "/tmp/trash"},
                lib.DockerEnv{"LIBRARY_PATH", "/usr/lib/my.lib.path/"},
            },
            From: []string{
                "override-docker-ansible",
                "override-docker-ansible-centos",
            },
        },
        []string{
            "(override-docker-ansible,TMPDIR=/tmp/trash)",
            "(override-docker-ansible,LIBRARY_PATH=/usr/lib/my.lib.path/)",
            "(override-docker-ansible-centos,TMPDIR=/tmp/trash)",
            "(override-docker-ansible-centos,LIBRARY_PATH=/usr/lib/my.lib.path/)",
        },
        *lib.NewBuildFile(),
    },
}

func TestBuildMatrixKey(t *testing.T) {
    for _, data := range keyTestData {
        actual := data.entry.Key()
        if actual != data.expected {
            t.Errorf("entry.Key() method failed to create key properly!\nexpected:\n%s\nactual:\n%s\n", data.expected, actual)
        }
    }
}

func TestMatrixConstruction(t *testing.T) {
    for _, data := range testData {
        expected := data.ExpectedMatrixKeys
        actual := data.Vectors.Product().Keys()

        sort.Strings(expected)

        if !reflect.DeepEqual(expected, actual) {
            t.Errorf("Sets() failed, \nactual: \n%s\n\nexpected: \n%s", actual, expected)
        }
    }
}
