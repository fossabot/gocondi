package gocondi

import (
	"database/sql"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"testing"
)

const (
	floatArrayParameter     = "float_array_parameter"
	floatParameter          = "float_parameter"
	float64ArrayParameter   = "float64_array_parameter"
	float64Parameter        = "float64_parameter"
	intArrayParameter       = "int_array_parameter"
	intParameter            = "int_parameter"
	int64ArrayParameter     = "int64_array_parameter"
	int64Parameter          = "int64_parameter"
	stringArrayParameter    = "string_array_parameter"
	stringParameter         = "string_parameter"
	boolParameter           = "bool_parameter"
	boolArrayParameter      = "bool_array_parameter"
	notExistParameter       = "not_exists"
	sqlConnection           = "host=localhost user=user password=password dbname=default sslmode=disable"
	testDatabaseName        = "testing"
	nonExistingDatabaseName = "not_exists"
)

var testParameters = map[string]interface{}{
	boolParameter:         true,
	boolArrayParameter:    []bool{true, false, true, true},
	floatArrayParameter:   []float32{-3, -2, -1, 0, 1, 2, 3},
	floatParameter:        float32(100.01),
	float64ArrayParameter: []float64{-3, -2, -1, 0, 1, 2, 3},
	float64Parameter:      float64(100.01),
	intArrayParameter:     []int{-3, -2, -1, 0, 1, 2, 3},
	intParameter:          100,
	int64ArrayParameter:   []int64{-3, -2, -1, 0, 1, 2, 3},
	int64Parameter:        int64(1),
	stringArrayParameter:  []string{"a", "b", "c"},
	stringParameter:       "string_parameter",
}

func TestGetContainer(t *testing.T) {
	assert.Panics(t, func() {
		GetContainer()
	})

	logger := logrus.New()
	logger.Out = ioutil.Discard
	logger.Level = logrus.PanicLevel

	Initialize(logger)

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

func TestContainer_GetBoolArrayParameter(t *testing.T) {
	containerObject.SetParameter(boolArrayParameter, testParameters[boolArrayParameter])

	parameter := containerObject.GetBoolArrayParameter(boolArrayParameter)

	assert.Equal(t, testParameters[boolArrayParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetBoolArrayParameter(notExistParameter)
	})
}

func TestContainer_GetBoolParameter(t *testing.T) {
	containerObject.SetParameter(boolParameter, testParameters[boolParameter])

	parameter := containerObject.GetBoolParameter(boolParameter)

	assert.Equal(t, testParameters[boolParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetBoolParameter(notExistParameter)
	})
}

func TestContainer_GetFloatArrayParameter(t *testing.T) {
	containerObject.SetParameter(floatArrayParameter, testParameters[floatArrayParameter])

	parameter := containerObject.GetFloatArrayParameter(floatArrayParameter)

	assert.Equal(t, testParameters[floatArrayParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetFloatArrayParameter(notExistParameter)
	})
}

func TestContainer_GetFloatParameter(t *testing.T) {
	containerObject.SetParameter(floatParameter, testParameters[floatParameter])

	parameter := containerObject.GetFloatParameter(floatParameter)

	assert.Equal(t, testParameters[floatParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetFloatParameter(notExistParameter)
	})
}

func TestContainer_GetFloat64ArrayParameter(t *testing.T) {
	containerObject.SetParameter(float64ArrayParameter, testParameters[float64ArrayParameter])

	parameter := containerObject.GetFloat64ArrayParameter(float64ArrayParameter)

	assert.Equal(t, testParameters[float64ArrayParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetFloat64ArrayParameter(notExistParameter)
	})
}

func TestContainer_GetFloat64Parameter(t *testing.T) {
	containerObject.SetParameter(float64Parameter, testParameters[float64Parameter])

	parameter := containerObject.GetFloat64Parameter(float64Parameter)

	assert.Equal(t, testParameters[float64Parameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetFloat64Parameter(notExistParameter)
	})
}

func TestContainer_GetIntArrayParameter(t *testing.T) {
	containerObject.SetParameter(intArrayParameter, testParameters[intArrayParameter])

	parameter := containerObject.GetIntArrayParameter(intArrayParameter)

	assert.Equal(t, testParameters[intArrayParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetIntArrayParameter(notExistParameter)
	})
}

func TestContainer_GetIntParameter(t *testing.T) {
	containerObject.SetParameter(intParameter, testParameters[intParameter])

	parameter := containerObject.GetIntParameter(intParameter)

	assert.Equal(t, testParameters[intParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetIntParameter(notExistParameter)
	})
}

func TestContainer_GetInt64ArrayParameter(t *testing.T) {
	containerObject.SetParameter(int64ArrayParameter, testParameters[int64ArrayParameter])

	parameter := containerObject.GetInt64ArrayParameter(int64ArrayParameter)

	assert.Equal(t, testParameters[int64ArrayParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetInt64ArrayParameter(notExistParameter)
	})
}

func TestContainer_GetInt64Parameter(t *testing.T) {
	containerObject.SetParameter(int64Parameter, testParameters[int64Parameter])

	parameter := containerObject.GetInt64Parameter(int64Parameter)

	assert.Equal(t, testParameters[int64Parameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetInt64Parameter(notExistParameter)
	})
}

func TestContainer_GetStringArrayParameter(t *testing.T) {
	containerObject.SetParameter(stringArrayParameter, testParameters[stringArrayParameter])

	parameter := containerObject.GetStringArrayParameter(stringArrayParameter)

	assert.Equal(t, testParameters[stringArrayParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetStringArrayParameter(notExistParameter)
	})
}

func TestContainer_GetStringParameter(t *testing.T) {
	containerObject.SetParameter(stringParameter, testParameters[stringParameter])

	parameter := containerObject.GetStringParameter(stringParameter)

	assert.Equal(t, testParameters[stringParameter], parameter)
	assert.Panics(t, func() {
		containerObject.GetStringParameter(notExistParameter)
	})
}
