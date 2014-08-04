package testing

import (
    "io"
    "os"
    "reflect"
    "text/template"
    "time"
)

func comma(index int, obj interface{}) bool {
    return index == reflect.ValueOf(obj).Len()-1
}

var FuncMap = template.FuncMap{
    "comma": comma,
    "add": func(a int, b int) int {
        return a + b
    },
}

type fileSystem interface {
    Open(name string) (file, error)
    Stat(name string) (os.FileInfo, error)
}

type file interface {
    io.Closer
    io.Reader
    io.ReaderAt
    io.Seeker
    Stat() (os.FileInfo, error)
}

type osFS struct{}

func (osFS) Open(name string) (file, error) { return nil, nil }

func (osFS) Stat(name string) (os.FileInfo, error) { return nil, nil }

type MockFile struct {
    MName    string
    MSize    int64
    MMode    os.FileMode
    MModTime time.Time
    MIsdir   bool
    MSys     interface{}
}

func (m MockFile) IsDir() bool {
    return m.MIsdir
}
func (m MockFile) ModTime() time.Time {
    return m.MModTime
}
func (m MockFile) Mode() os.FileMode {
    return m.MMode
}
func (m MockFile) Name() string {
    return m.MName
}
func (m MockFile) Size() int64 {
    return m.MSize
}
func (m MockFile) Sys() interface{} {
    return m.MSys
}
