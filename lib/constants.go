package lib

type ConstantValues struct {
    buildFileName  string
    configFileName string
}

func Constants() *ConstantValues {
    return &ConstantValues{
        buildFileName:  "test-flight-build.json",
        configFileName: "test-flight-config.json",
    }
}
