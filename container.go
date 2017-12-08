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

var containerObject = new(Container)

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
    database := containerObject.databases[name]

    if nil == database {
        errorMessage := fmt.Sprintf("Database connection with name %s not exists", name)

        if nil != containerObject.logger {
            containerObject.logger.Panic(errorMessage)
        }

        return nil, errors.New(errorMessage)
    }

    return database, nil
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
    return getArrayParameter(name)
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
    parametersInString := getArrayParameter(name)
    parameters := make([]int, len(parametersInString))

    for index, stringValue := range parametersInString {
        value, err := strconv.ParseInt(stringValue, 10, 0)

        if nil != err {
            containerObject.logger.WithError(err).Panicf("Error parsing int array parameter \"%s\"", name)
        }

        parameters[index] = int(value)
    }

    return parameters
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
    parametersInString := getArrayParameter(name)
    parameters := make([]int64, len(parametersInString))

    for index, stringValue := range parametersInString {
        value, err := strconv.ParseInt(stringValue, 10, 0)

        if nil != err {
            containerObject.logger.WithError(err).Panicf("Error parsing int64 array parameter \"%s\"", name)
        }

        parameters[index] = value
    }

    return parameters
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
    parametersInString := getArrayParameter(name)
    parameters := make([]float32, len(parametersInString))

    for index, stringValue := range parametersInString {
        value, err := strconv.ParseFloat(stringValue, 0)

        if nil != err {
            containerObject.logger.WithError(err).Panicf("Error parsing float array parameter \"%s\"", name)
        }

        parameters[index] = float32(value)
    }

    return parameters
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
    parametersInString := getArrayParameter(name)
    parameters := make([]float64, len(parametersInString))

    for index, stringValue := range parametersInString {
        value, err := strconv.ParseFloat(stringValue, 0)

        if nil != err {
            containerObject.logger.WithError(err).Panicf("Error parsing float array parameter \"%s\"", name)
        }

        parameters[index] = value
    }

    return parameters
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
    parametersInString := getArrayParameter(name)
    parameters := make([]bool, len(parametersInString))

    for index, stringValue := range parametersInString {
        value, err := strconv.ParseBool(stringValue)

        if nil != err {
            containerObject.logger.WithError(err).Panicf("Error parsing bool array parameter \"%s\"", name)
        }

        parameters[index] = value
    }

    return parameters
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
        containerObject.logger.WithError(err).Warningf("Error reading secrets folder")
        return
    }

    for _, secretFile := range secretFiles {
        // This is for prevent
        if secretFile.IsDir() {
            containerObject.logger.Warningf("Secrets folder has a folder!")
            continue
        }

        secretName := secretFile.Name()

        secretInBytes, err := ioutil.ReadFile(fmt.Sprintf("/run/secrets/%s", secretName))

        if nil != err {
            containerObject.logger.WithError(err).WithField("secret", secretName).Warningf("Error reading secret")
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
    return containerObject
}

func Initialize(logger *logrus.Logger) {
    containerObject.SetLogger(logger)
    containerObject.readSecretsFolder()
    containerObject.readParametersFromEnv()
}

func getParameter(name string) string {
    parameter := containerObject.parameters[name]

    if nil == parameter {
        parameter = getParameterFromSystem(name)
    }

    if nil == parameter && nil != containerObject.logger {
        containerObject.logger.Warningf("Parameter \"%s\" not exists", name)
    }

    return fmt.Sprintf("%v", parameter)
}

func getArrayParameter(name string) []string {
    parameter := getParameter(name)

    return strings.Split(parameter, ",")
}

func getParameterFromSystem(name string) interface{} {
    var parameter interface{}
    parameter = getParameterFromSecrets(name)

    if "" == parameter {
        parameter = getParameterFromEnv(name)
    }

    return parameter
}

func getParameterFromSecrets(name string) string {
    name = strings.ToLower(name)
    path := fmt.Sprintf("/run/secrets/%s", name)
    secret, _ := ioutil.ReadFile(path)

    return string(secret)
}

func getParameterFromEnv(name string) string {
    name = strings.ToUpper(name)

    return os.Getenv(name)
}
