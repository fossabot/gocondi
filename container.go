package gocondi

import (
    "database/sql"
    "github.com/sirupsen/logrus"
    "errors"
    "fmt"
)

const defaultDatabaseName = "default"

var containerObject = new(Container)

type containerInterface interface {
    GetDefaultDatabase() (*sql.DB, error)
    GetDatabase(name string) (*sql.DB, error)
    GetDatabases() []*sql.DB
    GetLogger() *logrus.Logger
    GetParameter(name string) interface{}
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

func (containerObject *Container) GetParameter(name string) interface{} {
    parameter := containerObject.parameters[name]

    if nil == parameter && nil != containerObject.logger {
        containerObject.logger.Warningf("Parameter \"%s\" not exists", name)
    }

    return parameter
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

func GetContainer() *Container {
    return containerObject
}
