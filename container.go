package gocondi

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"syscall"
)

const (
	defaultDatabaseName      = "default"
	envVarPrefix             = "GOCONDI_"
	driverPostgres           = "postgres"
	connectionStringPostgres = "host=%s port=%d user=%s password=%s dbname=%s sslmode=disable"
)

var c *Container

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

func (c *Container) GetDefaultDatabase() (*sql.DB, error) {
	return c.GetDatabase("default")
}

func (c *Container) GetDatabase(name string) (*sql.DB, error) {
	if database, exists := c.databases[name]; exists {
		return database, nil
	} else {
		errorMessage := fmt.Sprintf("Database connection with name %s not exists", name)

		return nil, errors.New(errorMessage)
	}
}

func (c *Container) GetDatabases() []*sql.DB {
	databases := make([]*sql.DB, 0)

	for _, value := range c.databases {
		databases = append(databases, value)
	}

	return databases
}

func (c *Container) GetLogger() *logrus.Logger {
	return c.logger
}

func (c *Container) GetStringParameter(name string) string {
	return getParameter(name)
}

func (c *Container) GetStringArrayParameter(name string) []string {
	var values []string

	if parameter, exists := c.parameters[name]; exists {
		values = (parameter).([]string)
	} else {
		c.logger.Panicf("Parameter \"%s\" not exists", name)
	}

	return values
}

func (c *Container) GetIntParameter(name string) int {
	parameter := getParameter(name)

	if parameter != "" {
		value, err := strconv.ParseInt(parameter, 10, 0)

		if nil != err {
			c.logger.WithError(err).Panicf("Error parsing int parameter \"%s\"", name)
		}

		return int(value)
	}

	return 0
}

func (c *Container) GetIntArrayParameter(name string) []int {
	var values []int

	if parameter, exists := c.parameters[name]; exists {
		values = (parameter).([]int)
	} else {
		c.logger.Panicf("Parameter \"%s\" not exists", name)
	}

	return values
}

func (c *Container) GetInt64Parameter(name string) int64 {
	parameter := getParameter(name)
	value, err := strconv.ParseInt(parameter, 10, 0)

	if nil != err {
		c.logger.WithError(err).Panicf("Error parsing int64 parameter \"%s\"", name)
	}

	return value
}

func (c *Container) GetInt64ArrayParameter(name string) []int64 {
	var values []int64

	if parameter, exists := c.parameters[name]; exists {
		values = (parameter).([]int64)
	} else {
		c.logger.Panicf("Parameter \"%s\" not exists", name)
	}

	return values
}

func (c *Container) GetFloatParameter(name string) float32 {
	parameter := getParameter(name)
	value, err := strconv.ParseFloat(parameter, 0)

	if nil != err {
		c.logger.WithError(err).Panicf("Error parsing float parameter \"%s\"", name)
	}

	return float32(value)
}

func (c *Container) GetFloatArrayParameter(name string) []float32 {
	var values []float32

	if parameter, exists := c.parameters[name]; exists {
		values = (parameter).([]float32)
	} else {
		c.logger.Panicf("Parameter \"%s\" not exists", name)
	}

	return values
}

func (c *Container) GetFloat64Parameter(name string) float64 {
	parameter := getParameter(name)
	value, err := strconv.ParseFloat(parameter, 0)

	if nil != err {
		c.logger.WithError(err).Panicf("Error parsing float64 parameter \"%s\"", name)
	}

	return value
}

func (c *Container) GetFloat64ArrayParameter(name string) []float64 {
	var values []float64

	if parameter, exists := c.parameters[name]; exists {
		values = (parameter).([]float64)
	} else {
		c.logger.Panicf("Parameter \"%s\" not exists", name)
	}

	return values
}

func (c *Container) GetBoolParameter(name string) bool {
	parameter := getParameter(name)
	value, err := strconv.ParseBool(parameter)

	if nil != err {
		c.logger.WithError(err).Panicf("Error parsing bool parameter \"%s\"", name)
	}

	return value
}

func (c *Container) GetBoolArrayParameter(name string) []bool {
	var values []bool

	if parameter, exists := c.parameters[name]; exists {
		values = (parameter).([]bool)
	} else {
		c.logger.Panicf("Parameter \"%s\" not exists", name)
	}

	return values
}

func (c *Container) GetParameters() map[string]interface{} {
	return c.parameters
}

func (c *Container) SetDefaultDatabase(database *sql.DB) *Container {
	return c.SetDatabase("default", database)
}

func (c *Container) SetDatabase(name string, database *sql.DB) *Container {
	if nil == c.databases {
		c.databases = make(map[string]*sql.DB)
	}

	c.databases[name] = database

	return c
}

func (c *Container) SetDatabases(databases map[string]*sql.DB) *Container {
	c.databases = databases

	return c
}

func (c *Container) SetLogger(logger *logrus.Logger) *Container {
	c.logger = logger

	return c
}

func (c *Container) SetParameter(name string, parameter interface{}) *Container {
	if nil == c.parameters {
		c.parameters = make(map[string]interface{})
	}

	c.parameters[name] = parameter

	return c
}

func (c *Container) SetParameters(parameters map[string]interface{}) *Container {
	c.parameters = parameters

	return c
}

func (c *Container) loadSecretsFolder() {
	secretFiles, err := ioutil.ReadDir("/run/secrets")

	if nil != err {
		return
	}

	for _, secretFile := range secretFiles {
		// This is for prevent
		if secretFile.IsDir() {
			c.logger.Warningf("Secrets folder has a subfolder!")
			continue
		}

		secretName := secretFile.Name()
		secretInBytes, err := ioutil.ReadFile(fmt.Sprintf("/run/secrets/%s", secretName))

		if nil != err {
			continue
		}

		secret := string(secretInBytes)

		c.SetParameter(secretName, secret)
	}
}

func (c *Container) loadParametersFromEnv() {
	for _, pair := range os.Environ() {
		split := strings.Split(pair, "=")
		name := split[0]
		value := split[1]

		if strings.Contains(name, envVarPrefix) {
			name = strings.Split(name, envVarPrefix)[1]
			name = strings.ToLower(name)

			if "" != name {
				c.SetParameter(name, value)
			}
		}
	}
}

func (c *Container) loadDefaultDatabase() {
	const logTag = "gocondi.loadDefaultDatabase()"
	c.logger.Debugf("%s -> START", logTag)

	defer func() {
		if err := recover(); err != nil {
			c.logger.WithField("error", err).Error()
		}
	}()

	host := c.GetStringParameter("database_host")
	port := c.GetIntParameter("database_port")
	user := c.GetStringParameter("database_username")
	password := c.GetStringParameter("database_password")
	dbName := c.GetStringParameter("database_name")
	driver := c.GetStringParameter("database_driver")

	if host == "" {
		return
	}

	var connectionString string

	switch driver {
	case driverPostgres:
		connectionString = fmt.Sprintf(connectionStringPostgres, host, port, user, password, dbName)
		break
	default:
		err := errors.New(fmt.Sprintf("Database driver '%s' is not supported", driver))

		c.logger.WithError(err).Panicf("Error loading default database")
		return
	}

	db, err := sql.Open(driver, connectionString)

	if nil != err {
		c.logger.WithError(err).Panicf("Error loading default database")
	}

	c.SetDefaultDatabase(db)

	c.logger.Debugf("%s -> Database connected", logTag)
	c.logger.Debugf("%s -> END", logTag)
}

func (c *Container) CloseDatabases() {
	for name, database := range c.databases {
		if err := database.Close(); err != nil {
			c.logger.WithField("name", name).WithError(err).Errorf("Error closing database")
		} else {
			c.logger.WithField("name", name).Debugf("Closed database")
		}
	}
}

func (c *Container) Close() {
	c.CloseDatabases()
}

func (c *Container) Reload() {
	c.logger.Infof("Reloading config...")
	c.loadSecretsFolder()
	c.loadParametersFromEnv()
	c.loadDefaultDatabase()
}

func GetContainer() *Container {
	if nil == c {
		panic("Container isn't initialized. You must use gocondi.Initialize(logger) first.")
	}

	return c
}

func Initialize(logger *logrus.Logger) {
	c = new(Container)
	c.SetLogger(logger)
	c.loadSecretsFolder()
	c.loadParametersFromEnv()
	c.loadDefaultDatabase()
	listenReloadSignal()
}

func listenReloadSignal() {
	c.logger.Debugf("Listening reload signal...")
	signalChannel := make(chan os.Signal)
	signal.Notify(signalChannel, syscall.SIGHUP)
	go func() {
		<-signalChannel
		GetContainer().Reload()
	}()
}

func getParameter(name string) string {
	var value string

	if parameter, exists := c.parameters[name]; exists {
		value = fmt.Sprintf("%v", parameter)
	} else {
		c.logger.Panicf("Parameter \"%s\" not exists", name)
	}

	return value
}
