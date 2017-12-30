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

	assert.Equal(t, c, GetContainer())
}

func TestContainer_SetParameters(t *testing.T) {
	c.SetParameters(testParameters)
	parameters := c.parameters

	assert.Equal(t, testParameters, parameters)
}

func TestContainer_SetParameter(t *testing.T) {
	for name, parameter := range testParameters {
		c.SetParameter(name, parameter)

		assert.Equal(t, parameter, c.parameters[name])
	}

	assert.Nil(t, c.parameters[notExistParameter])
	assert.NotEqual(t, testParameters[floatArrayParameter], c.parameters[intArrayParameter])
	assert.NotEmpty(t, c.parameters[stringParameter])
}

func TestContainer_SetDefaultDatabase(t *testing.T) {
	db, _ := sql.Open("postgres", sqlConnection)
	c.SetDefaultDatabase(db)

	assert.Equal(t, db, c.databases[defaultDatabaseName])
}

func TestContainer_SetDatabases(t *testing.T) {
	db, _ := sql.Open("postgres", sqlConnection)
	databases := map[string]*sql.DB{testDatabaseName: db}
	c.SetDatabases(databases)

	assert.Equal(t, databases, c.databases)
}

func TestContainer_SetDatabase(t *testing.T) {
	db, _ := sql.Open("postgres", sqlConnection)
	c.SetDatabase(testDatabaseName, db)

	assert.Equal(t, db, c.databases[testDatabaseName])
}

func TestContainer_GetDatabase(t *testing.T) {
	db, _ := sql.Open("postgres", sqlConnection)
	c.SetDatabase(testDatabaseName, db)
	databaseInContainer, err := c.GetDatabase(testDatabaseName)

	assert.Equal(t, db, databaseInContainer)
	assert.NoError(t, err)

	databaseInContainer, err = c.GetDatabase(nonExistingDatabaseName)

	assert.Nil(t, databaseInContainer)
	assert.Error(t, err)
}

func TestContainer_GetDatabases(t *testing.T) {
	db, _ := sql.Open("postgres", sqlConnection)
	databases := map[string]*sql.DB{testDatabaseName: db}
	c.SetDatabases(databases)
	databasesInContainer := c.GetDatabases()

	databasesArray := make([]*sql.DB, 0)

	for _, value := range databasesInContainer {
		databasesArray = append(databasesArray, value)
	}

	assert.Equal(t, databasesArray, databasesInContainer)
}

func TestContainer_GetDefaultDatabase(t *testing.T) {
	db, _ := sql.Open("postgres", sqlConnection)
	c.SetDefaultDatabase(db)
	defaultDatabase, err := c.GetDefaultDatabase()

	assert.Equal(t, db, defaultDatabase)
	assert.NoError(t, err)

	c.databases = *new(map[string]*sql.DB)

	defaultDatabase, err = c.GetDefaultDatabase()

	assert.Nil(t, defaultDatabase)
	assert.Error(t, err)
}

func TestContainer_SetLogger(t *testing.T) {
	logger := logrus.New()
	c.SetLogger(logger)

	assert.Equal(t, logger, c.logger)
}

func TestContainer_GetLogger(t *testing.T) {
	logger := logrus.New()
	c.SetLogger(logger)
	loggerInContainer := c.GetLogger()

	assert.Equal(t, logger, loggerInContainer)
}

func TestContainer_GetBoolArrayParameter(t *testing.T) {
	c.SetParameter(boolArrayParameter, testParameters[boolArrayParameter])

	parameter := c.GetBoolArrayParameter(boolArrayParameter)

	assert.Equal(t, testParameters[boolArrayParameter], parameter)
	assert.Panics(t, func() {
		c.GetBoolArrayParameter(notExistParameter)
	})
}

func TestContainer_GetBoolParameter(t *testing.T) {
	c.SetParameter(boolParameter, testParameters[boolParameter])

	parameter := c.GetBoolParameter(boolParameter)

	assert.Equal(t, testParameters[boolParameter], parameter)
	assert.Panics(t, func() {
		c.GetBoolParameter(notExistParameter)
	})
}

func TestContainer_GetFloatArrayParameter(t *testing.T) {
	c.SetParameter(floatArrayParameter, testParameters[floatArrayParameter])

	parameter := c.GetFloatArrayParameter(floatArrayParameter)

	assert.Equal(t, testParameters[floatArrayParameter], parameter)
	assert.Panics(t, func() {
		c.GetFloatArrayParameter(notExistParameter)
	})
}

func TestContainer_GetFloatParameter(t *testing.T) {
	c.SetParameter(floatParameter, testParameters[floatParameter])

	parameter := c.GetFloatParameter(floatParameter)

	assert.Equal(t, testParameters[floatParameter], parameter)
	assert.Panics(t, func() {
		c.GetFloatParameter(notExistParameter)
	})
}

func TestContainer_GetFloat64ArrayParameter(t *testing.T) {
	c.SetParameter(float64ArrayParameter, testParameters[float64ArrayParameter])

	parameter := c.GetFloat64ArrayParameter(float64ArrayParameter)

	assert.Equal(t, testParameters[float64ArrayParameter], parameter)
	assert.Panics(t, func() {
		c.GetFloat64ArrayParameter(notExistParameter)
	})
}

func TestContainer_GetFloat64Parameter(t *testing.T) {
	c.SetParameter(float64Parameter, testParameters[float64Parameter])

	parameter := c.GetFloat64Parameter(float64Parameter)

	assert.Equal(t, testParameters[float64Parameter], parameter)
	assert.Panics(t, func() {
		c.GetFloat64Parameter(notExistParameter)
	})
}

func TestContainer_GetIntArrayParameter(t *testing.T) {
	c.SetParameter(intArrayParameter, testParameters[intArrayParameter])

	parameter := c.GetIntArrayParameter(intArrayParameter)

	assert.Equal(t, testParameters[intArrayParameter], parameter)
	assert.Panics(t, func() {
		c.GetIntArrayParameter(notExistParameter)
	})
}

func TestContainer_GetIntParameter(t *testing.T) {
	c.SetParameter(intParameter, testParameters[intParameter])

	parameter := c.GetIntParameter(intParameter)

	assert.Equal(t, testParameters[intParameter], parameter)
	assert.Panics(t, func() {
		c.GetIntParameter(notExistParameter)
	})
}

func TestContainer_GetInt64ArrayParameter(t *testing.T) {
	c.SetParameter(int64ArrayParameter, testParameters[int64ArrayParameter])

	parameter := c.GetInt64ArrayParameter(int64ArrayParameter)

	assert.Equal(t, testParameters[int64ArrayParameter], parameter)
	assert.Panics(t, func() {
		c.GetInt64ArrayParameter(notExistParameter)
	})
}

func TestContainer_GetInt64Parameter(t *testing.T) {
	c.SetParameter(int64Parameter, testParameters[int64Parameter])

	parameter := c.GetInt64Parameter(int64Parameter)

	assert.Equal(t, testParameters[int64Parameter], parameter)
	assert.Panics(t, func() {
		c.GetInt64Parameter(notExistParameter)
	})
}

func TestContainer_GetStringArrayParameter(t *testing.T) {
	c.SetParameter(stringArrayParameter, testParameters[stringArrayParameter])

	parameter := c.GetStringArrayParameter(stringArrayParameter)

	assert.Equal(t, testParameters[stringArrayParameter], parameter)
	assert.Panics(t, func() {
		c.GetStringArrayParameter(notExistParameter)
	})
}

func TestContainer_GetStringParameter(t *testing.T) {
	c.SetParameter(stringParameter, testParameters[stringParameter])

	parameter := c.GetStringParameter(stringParameter)

	assert.Equal(t, testParameters[stringParameter], parameter)
	assert.Panics(t, func() {
		c.GetStringParameter(notExistParameter)
	})
}
