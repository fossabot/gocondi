package gocondi

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "database/sql"
    "github.com/sirupsen/logrus"
)

const (
    floatArrayParameter     = "float_array_parameter"
    floatParameter          = "float_parameter"
    intArrayParameter       = "int_array_parameter"
    intParameter            = "int_parameter"
    stringArrayParameter    = "string_array_parameter"
    stringParameter         = "string_parameter"
    notExistParameter       = "not_exists"
    sqlConnection           = "host=localhost user=user password=password dbname=iot sslmode=disable"
    testDatabaseName        = "testing"
    nonExistingDatabaseName = "not_exists"
)

var testParameters = map[string]interface{}{
    floatArrayParameter:  []float32{-3, -2, -1, 0, 1, 2, 3},
    floatParameter:       100.01,
    intArrayParameter:    []int{-3, -2, -1, 0, 1, 2, 3},
    intParameter:         100,
    stringArrayParameter: []string{"a", "b", "c"},
    stringParameter:      "string_parameter",
}

func TestGetContainer(t *testing.T) {
    assert.Equal(t, containerObject, GetContainer())
}

func TestContainer_SetParameters(t *testing.T) {
    containerObject.SetParameters(testParameters)
    parameters := containerObject.parameters

    assert.Equal(t, testParameters, parameters)
}

func TestContainer_SetParameter(t *testing.T) {
    for name, parameter := range testParameters {
        containerObject.SetParameter(name, parameter)

        assert.Equal(t, parameter, containerObject.parameters[name])
    }

    assert.Nil(t, containerObject.parameters[notExistParameter])
    assert.NotEqual(t, testParameters[floatArrayParameter], containerObject.parameters[intArrayParameter])
    assert.NotEmpty(t, containerObject.parameters[stringParameter])
}

func TestContainer_SetDefaultDatabase(t *testing.T) {
    db, _ := sql.Open("postgres", sqlConnection)
    containerObject.SetDefaultDatabase(db)

    assert.Equal(t, db, containerObject.databases[defaultDatabaseName])
}

func TestContainer_SetDatabases(t *testing.T) {
    db, _ := sql.Open("postgres", sqlConnection)
    databases := map[string]*sql.DB{testDatabaseName: db}
    containerObject.SetDatabases(databases)

    assert.Equal(t, databases, containerObject.databases)
}

func TestContainer_SetDatabase(t *testing.T) {
    db, _ := sql.Open("postgres", sqlConnection)
    containerObject.SetDatabase(testDatabaseName, db)

    assert.Equal(t, db, containerObject.databases[testDatabaseName])
}

func TestContainer_GetDatabase(t *testing.T) {
    db, _ := sql.Open("postgres", sqlConnection)
    containerObject.SetDatabase(testDatabaseName, db)
    databaseInContainer, err := containerObject.GetDatabase(testDatabaseName)

    assert.Equal(t, db, databaseInContainer)
    assert.NoError(t, err)

    databaseInContainer, err = containerObject.GetDatabase(nonExistingDatabaseName)
    //assert.NotEqual(t, db, databaseInContainer)
    assert.Nil(t, databaseInContainer)
    assert.Error(t, err)
}

func TestContainer_GetDatabases(t *testing.T) {
    db, _ := sql.Open("postgres", sqlConnection)
    databases := map[string]*sql.DB{testDatabaseName: db}
    containerObject.SetDatabases(databases)
    databasesInContainer := containerObject.GetDatabases()

    databasesArray := make([]*sql.DB, 0)

    for _, value := range databasesInContainer {
        databasesArray = append(databasesArray, value)
    }

    assert.Equal(t, databasesArray, databasesInContainer)
}

func TestContainer_GetDefaultDatabase(t *testing.T) {
    db, _ := sql.Open("postgres", sqlConnection)
    containerObject.SetDefaultDatabase(db)
    defaultDatabase, err := containerObject.GetDefaultDatabase()

    assert.Equal(t, db, defaultDatabase)
    assert.NoError(t, err)

    containerObject.databases = *new(map[string]*sql.DB)

    defaultDatabase, err = containerObject.GetDefaultDatabase()

    // assert.NotEqual(t, db, defaultDatabase)
    assert.Nil(t, defaultDatabase)
    assert.Error(t, err)
}

func TestContainer_SetLogger(t *testing.T) {
    logger := logrus.New()
    containerObject.SetLogger(logger)

    assert.Equal(t, logger, containerObject.logger)
}

func TestContainer_GetLogger(t *testing.T) {
    logger := logrus.New()
    containerObject.SetLogger(logger)
    loggerInContainer := containerObject.GetLogger()

    assert.Equal(t, logger, loggerInContainer)
}