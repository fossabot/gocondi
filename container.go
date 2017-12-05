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

const defaultDatabaseName = "default"

var containerObject = new(Container)

// In the container initialization first list all files in /run/secret
// second read all env var becoming with "GOCONDI
// and last in the array
// TODO Read from yml
type containerInterface interface {
    GetDefaultDatabase() (*sql.DB, error)
    GetDatabase(name string) (*sql.DB, error)
    GetDatabases() []*sql.DB
    GetLogger() *logrus.Logger
    GetStringParameter(name string) string
    GetStringArrayParameter(name string) string
    GetIntParameter(name string) int
    GetFloatParameter(name string) float32
    GetBoolParameter(name string) bool
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

func (containerObject *Container) GetIntParameter(name string) int {
    parameter := getParameter(name)
    value, _ := strconv.ParseInt(parameter, 10, 0)

    return int(value)
}

func (containerObject *Container) GetFloatParameter(name string) float32 {
    parameter := getParameter(name)
    value, _ := strconv.ParseFloat(parameter, 0)

    return float32(value)
}

func (containerObject *Container) GetBoolParameter(name string) bool {
    parameter := getParameter(name)
    value, _ := strconv.ParseBool(parameter)

    return value
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
        name := pair[0]

        containerObject.logger.WithField("var", name).Debugf("Reading env var")
    }
}

func GetContainer() *Container {
    return containerObject
}

func New() *Container {
    containerObject.readSecretsFolder()
    containerObject.readParametersFromEnv()

    return containerObject
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
