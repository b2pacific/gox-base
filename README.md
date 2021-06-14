Gox-Base project provide utilities which is used commonly in all applications.

1. Serialization utils
2. Json file to object
3. Yaml file to object
4. XML file to object
5. ...

# Config

You can use this for common application configuration

```go
type configToTest struct {
App    App    `json:"app" yaml:"app"`
Logger Logger `json:"logger" yaml:"logger"`
}

// Read struct with YAML
conf := &configToTest{}
err := serialization.ReadYamlFromString(yamlConfig, conf)
assert.NoError(t, err)


// Read struct with Json
conf := &configToTest{}
err := serialization.JsonToObject(jsonString, conf)
assert.NoError(t, err)
```

You can also extend an existing type to include more. For example, we will add one more property called
```some``` in ```App``` type

```go

// Extended "App" object with one more type "Some"
// NOTE - in case you have YAM as source you must put `yaml:",inline"`
type extendedAppObject struct {
App         `yaml:",inline"`
Some string `json:"some"`
}

type configToTestWithExtendedAppObject struct {
App    extendedAppObject `json:"app" yaml:"app"`
Logger Logger            `json:"logger" yaml:"logger"`
}
```

# Utility

##### Setup cross function which is usd in almost all apis to be used in gox

```go
import (
"github.com/devlibx/gox-base"
"go.uber.org/zap"
)

zapConfig := zap.NewDevelopmentConfig()
zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
crossFunction := gox.NewCrossFunction(zapConfig.Build())

For test cases:
============== =
import "github.com/devlibx/gox-base/test"
cf, gomockController := test.MockCf(t, zap.InfoLevel)
```

#### Convert anything to string

This utility will convert int, bool, interface{} to string. Object will output a json string.

```golang
out, _ := Stringify(10)
// Output = "10"

boolOut, _ := Stringify(true)
// Output = "true"

type utilTestStruct struct {
IntValue    int    `json:"int"`
BoolValue   bool   `json:"bool"`
StringValue string `json:"string"`
}

objectOut, _ := Stringify(utilTestStruct{
IntValue:    10,
BoolValue:   false,
StringValue: "some value",
})
// Output = {"int":10,"bool":false,"string":"some value"}
```

###### Stringify with error suppressed

If you don't want to handle error and have default value on error then you can use "suppress error" version of this
method

```golang
intOut1 := StringifySuppressError(10, "0")
// Output = "10"

If there is a error when input is bad then you will get the default
value "0"
```

### Logging

You can use logger from Uber

```shell
var cfg zap.Config
cfg.Level = zap.NewAtomicLevelAt(zap.DebugLevel)
// cfg.Encoding = "console"
cfg.Encoding = "json"
cfg.OutputPaths = []string{"stdout", "/tmp/logs1"}
cfg.ErrorOutputPaths = []string{"stdout", "/tmp/logs2"}
cfg.EncoderConfig = zapcore.EncoderConfig{
    MessageKey:     "message",
    LevelKey:       "level",
    EncodeLevel:    zapcore.LowercaseLevelEncoder,
    EncodeDuration: zapcore.SecondsDurationEncoder,
    EncodeCaller:   zapcore.ShortCallerEncoder,
    StacktraceKey:  "stacktrace",
    TimeKey:        "timestamp",
    EncodeTime:     zapcore.ISO8601TimeEncoder,
}

// Production encoder
// encoderCfg := zap.NewProductionEncoderConfig()
// encoderCfg.TimeKey = "timestamp"
// encoderCfg.EncodeTime = zapcore.ISO8601TimeEncoder
// cfg.EncoderConfig = encoderCfg

logger, _ := cfg.Build()
defer logger.Sync()

logger.Debug("Logger from parent level - it has not key-value")

// A module level logger (you can log user Id or any param here)
subModuleLogger := logger.With(zap.String("userId", "1234"))
subModuleLogger.Debug("this is a logger for sub-module")

logger.Debug("Logger from parent level - it has not key-value")
```

Some more example https://github.com/uber-go/zap/blob/master/example_test.go

----

## Read Parameterized Yaml File

Sometime we need to read a Yaml file which contains our configuration. We need to create different files for different
env. You can use to read a Yaml file in a object

1. First all the environment variables are replaced in the configuration file
2. Then base on the env you have passed, it will create the correct yaml file and will build object

```shell
yaml2Go := Yaml2Go{}
err := ReadParameterizedYamlFile("../testdata/app_with_env_var_and_params.yml", &yaml2Go, "prod")
assert.NoError(t, err)
```

##### Your input Yaml file - which provides values for different environments

```yaml
server_config:
  servers:
    testServer:
      host: "env:string: prod=$PRDO_testServer; stage=$STAGE_testServer; default=localhost"
      port: "env:int: prod=80; stage=80; dev=8080; default=8090"
  apis:
    getPost:
      method: POST
      path: /get
      timeout: "env:int: prod=10; stage=20; dev=30; default=1000"

client:
  id: "env:string: prod=prod_client; stage=stage_client; dev=dev_client; default=random_client"
  enabled: "env:bool: prod=true; stage=false; dev=false; default=false"
  price: "env:float: prod=10.001; stage=10.002; dev=10.003; default=10.004"
```

##### Converted Yaml which is loaded into the object

```yaml
client:
  enabled: true
  id: prod_client
  price: 10.001
server_config:
  apis:
    getPost:
      method: POST
      path: /get
      timeout: 10
  servers:
    testServer:
      host: test.prod
      port: 80
```