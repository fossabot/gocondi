package gocondi

import (
    "database/sql"
    "github.com/sirupsen/logrus"
    "errors"
    "fmt"
    "strings"
    "io/ioutil"
    "os"
    "strconv"
)

const (
    defaultDatabaseName = "default"
    envVarPrefix        = "GOCONDI_"
)

var containerObject *Container

type containerInterface interface {
    GetDefaultDatabase() (*sql.DB, error)
    GetDatabase(name string) (*sql.DB, error)
    GetDatabases() []*sql.DB
    GetLogger() *logrus.Logger
    GetParameters() map[string]interface{}
    GetStringParameter(name string) string
    GetStringArrayParameter(name string) []string
    GetIntParameter(name string) int
    GetIntArrayParameter(name string) []int
    GetInt64Parameter(name string) int64
    GetInt64ArrayParameter(name string) []int64
    GetFloatParameter(name string) float32
    GetFloatArrayParameter(name string) []float32
    GetFloat64Parameter(name string) float64
    GetFloat64ArrayParameter(name string) []float64
    GetBoolParameter(name string) bool
    GetBoolArrayParameter(name string) []bool
    SetDefaultDatabase(database *sql.DB) *Container
    SetDatabase(name string, database *sql.DB) *Container
    SetDatabases(databases map[string]*sql.DB) *Container
    SetLogger(logger *logrus.Logger) *Container
    SetParameter(name string, parameter interface{}) *Container
    SetParameters(parameters map[string]interface{}) *Container
}

type Container struct {
    containerInterface
    databases  map[string]*sql.DB
    logger     *logrus.Logger
    parameters map[string]interface{}
}

func (containerObject *Container) GetDefaultDatabase() (*sql.DB, error) {
    return containerObject.GetDatabase("default")
}

func (containerObject *Container) GetDatabase(name string) (*sql.DB, error) {
    if database, exists := containerObject.databases[name]; exists {
        return database, nil
    } else {
        errorMessage := fmt.Sprintf("Database connection with name %s not exists", name)

        return nil, errors.New(errorMessage)
    }
}

func (containerObject *Container) GetDatabases() []*sql.DB {
    databases := make([]*sql.DB, 0)

    for _, value := range containerObject.databases {
        databases = append(databases, value)
    }

    return databases
}

func (containerObject *Container) GetLogger() *logrus.Logger {
    return containerObject.logger
}

func (containerObject *Container) GetStringParameter(name string) string {
    return getParameter(name)
}

func (containerObject *Container) GetStringArrayParameter(name string) []string {
    var values []string

    if parameter, exists := containerObject.parameters[name]; exists {
        values = (parameter).([]string)
    } else {
        containerObject.logger.Panicf("Parameter \"%s\" not exists", name)
    }

    return values
}

func (containerObject *Container) GetIntParameter(name string) int {
    parameter := getParameter(name)
    value, err := strconv.ParseInt(parameter, 10, 0)

    if nil != err {
        containerObject.logger.WithError(err).Panicf("Error parsing int parameter \"%s\"", name)
    }

    return int(value)
}

func (containerObject *Container) GetIntArrayParameter(name string) []int {
    var values []int

    if parameter, exists := containerObject.parameters[name]; exists {
        values = (parameter).([]int)
    } else {
        containerObject.logger.Panicf("Parameter \"%s\" not exists", name)
    }

    return values
}

func (containerObject *Container) GetInt64Parameter(name string) int64 {
    parameter := getParameter(name)
    value, err := strconv.ParseInt(parameter, 10, 0)

    if nil != err {
        containerObject.logger.WithError(err).Panicf("Error parsing int64 parameter \"%s\"", name)
    }

    return value
}

func (containerObject *Container) GetInt64ArrayParameter(name string) []int64 {
    var values []int64

    if parameter, exists := containerObject.parameters[name]; exists {
        values = (parameter).([]int64)
    } else {
        containerObject.logger.Panicf("Parameter \"%s\" not exists", name)
    }

    return values
}

func (containerObject *Container) GetFloatParameter(name string) float32 {
    parameter := getParameter(name)
    value, err := strconv.ParseFloat(parameter, 0)

    if nil != err {
        containerObject.logger.WithError(err).Panicf("Error parsing float parameter \"%s\"", name)
    }

    return float32(value)
}

func (containerObject *Container) GetFloatArrayParameter(name string) []float32 {
    var values []float32

    if parameter, exists := containerObject.parameters[name]; exists {
        values = (parameter).([]float32)
    } else {
        containerObject.logger.Panicf("Parameter \"%s\" not exists", name)
    }

    return values
}

func (containerObject *Container) GetFloat64Parameter(name string) float64 {
    parameter := getParameter(name)
    value, err := strconv.ParseFloat(parameter, 0)

    if nil != err {
        containerObject.logger.WithError(err).Panicf("Error parsing float64 parameter \"%s\"", name)
    }

    return value
}

func (containerObject *Container) GetFloat64ArrayParameter(name string) []float64 {
    var values []float64

    if parameter, exists := containerObject.parameters[name]; exists {
        values = (parameter).([]float64)
    } else {
        containerObject.logger.Panicf("Parameter \"%s\" not exists", name)
    }

    return values
}

func (containerObject *Container) GetBoolParameter(name string) bool {
    parameter := getParameter(name)
    value, err := strconv.ParseBool(parameter)

    if nil != err {
        containerObject.logger.WithError(err).Panicf("Error parsing bool parameter \"%s\"", name)
    }

    return value
}

func (containerObject *Container) GetBoolArrayParameter(name string) []bool {
    var values []bool

    if parameter, exists := containerObject.parameters[name]; exists {
        values = (parameter).([]bool)
    } else {
        containerObject.logger.Panicf("Parameter \"%s\" not exists", name)
    }

    return values
}

func (containerObject *Container) GetParameters() map[string]interface{} {
    return containerObject.parameters
}

func (containerObject *Container) SetDefaultDatabase(database *sql.DB) *Container {
    return containerObject.SetDatabase("default", database)
}

func (containerObject *Container) SetDatabase(name string, database *sql.DB) *Container {
    if nil == containerObject.databases {
        containerObject.databases = make(map[string]*sql.DB)
    }

    containerObject.databases[name] = database

    return containerObject
}

func (containerObject *Container) SetDatabases(databases map[string]*sql.DB) *Container {
    containerObject.databases = databases

    return containerObject
}

func (containerObject *Container) SetLogger(logger *logrus.Logger) *Container {
    containerObject.logger = logger

    return containerObject
}

func (containerObject *Container) SetParameter(name string, parameter interface{}) *Container {
    if nil == containerObject.parameters {
        containerObject.parameters = make(map[string]interface{})
    }

    containerObject.parameters[name] = parameter

    return containerObject
}

func (containerObject *Container) SetParameters(parameters map[string]interface{}) *Container {
    containerObject.parameters = parameters

    return containerObject
}

func (containerObject *Container) readSecretsFolder() {
    secretFiles, err := ioutil.ReadDir("/run/secrets")

    if nil != err {
        return
    }

    for _, secretFile := range secretFiles {
        // This is for prevent
        if secretFile.IsDir() {
            containerObject.logger.Warningf("Secrets folder has a subfolder!")
            continue
        }

        secretName := secretFile.Name()

        secretInBytes, err := ioutil.ReadFile(fmt.Sprintf("/run/secrets/%s", secretName))

        if nil != err {
            continue
        }

        secret := string(secretInBytes)

        containerObject.SetParameter(secretName, secret)
    }
}

func (containerObject *Container) readParametersFromEnv() {
    for _, pair := range os.Environ() {
        split := strings.Split(pair, "=")
        name := split[0]
        value := split[1]

        if strings.Contains(name, envVarPrefix) {
            name = strings.Split(name, envVarPrefix)[1]
            name = strings.ToLower(name)

            if "" != name {
                containerObject.SetParameter(name, value)
            }
        }
    }
}

func GetContainer() *Container {
    if nil == containerObject {
        panic("Container isn't initialized. You must use gocondi.Initialize(logger) first.")
    }

    return containerObject
}

func Initialize(logger *logrus.Logger) {
    containerObject = new(Container)
    containerObject.SetLogger(logger)
    containerObject.readSecretsFolder()
    containerObject.readParametersFromEnv()
}

func getParameter(name string) string {
    var value string

    if parameter, exists := containerObject.parameters[name]; exists {
        value = fmt.Sprintf("%v", parameter)
    } else {
        containerObject.logger.Panicf("Parameter \"%s\" not exists", name)
    }

    return value
}
