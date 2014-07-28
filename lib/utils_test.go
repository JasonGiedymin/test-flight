package lib

import (
    libtesting "github.com/JasonGiedymin/test-flight/lib/testing"
    "os"
    "testing"
)

func TestConvertFiles(t *testing.T) {
    mockFiles := []os.FileInfo{
        libtesting.MockFile{MName: "TestFile1"},
        libtesting.MockFile{MName: "TestFile2"},
    }

    expected := []string{
        "TestFile1",
        "TestFile2",
    }

    actual := ConvertFiles(mockFiles)

    for i, _ := range actual {
        currActual := actual[i]
        currExpected := expected[i]
        if currActual != currActual {
            t.Error("Expected:", currExpected, ", Got:", currActual)
        }
    }
}
