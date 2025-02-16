package config

import (
	"encoding/json"
	"fmt"
	"os"

	"strconv"
)

type config struct {
	db      database `json:"database"`
	Service service  `json:"service"`
	Routes  []Routes `json:"routes"`
}

type service struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Port    int    `json:"port"`
}

type Routes struct {
	Path    string `json:"path"`
	Handler string `json:"handler"`
	Method  string `json:"method"`
}

type Config interface {
	GetPort() int
	GetRoutes() []Routes
}

type database struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Database string `json:"database"`
}

// NewConfig creates a new Config instance with default values
func NewConfig() config {
	c := &config{}

	// load the configuratiion file using runtime
	conf, err := os.ReadFile("config/configuration.json")
	if err != nil {
		fmt.Printf("Error loading configuration file: %v\n", err)
		panic(err)
	}

	// unmarshal the configuration JSON into a Config struct
	// and change the configuration string to a byte array before unmarshalling
	err = json.Unmarshal([]byte(conf), &c)
	if err != nil {
		fmt.Printf("Error unmarshalling configuration: %v\n", err)
		panic(err)
	}

	// return the Config struct
	return *c
}

// returns the port from the configuration
func (c *config) GetPort() int {
	return c.Service.Port
}

// returns the routes from the configuration
func (c *config) GetRoutes() []Routes {
	return c.Routes
}

func (c *config) GetDB() database {
	return c.db
}

var pwd string

func (c *config) GetDatabasePassword() string {
	if pwd != "" {
		return pwd
	}

	pwd = GetEnv("PSQL_FILE_DATABASE_PASSWORD")
	if pwd == "" {
		return c.GetDB().Password
	}

	return pwd
}

var uname string

func (c *config) GetDatabaseUsername() string {
	if uname != "" {
		return uname
	}

	uname = GetEnv("PSQL_FILE_DATABASE_USERNAME")
	if uname == "" {
		return c.GetDB().Username
	}

	return uname
}

func (c *config) GetDatabaseHost() string {
	host := GetEnv("DB_HOST")
	if host == "" {
		return c.GetDB().Host
	}

	return host
}

func (c *config) GetDatabasePort() int {
	p := GetEnv("DB_PORT")
	if p == "" {
		return c.GetDB().Port
	}

	p_int, err := strconv.Atoi(p)

	if err != nil {
		fmt.Println("Error converting DB_PORT to int: ", err)
		return c.GetDB().Port
	}

	return p_int
}

func (c *config) GetDatabaseName() string {
	name := GetEnv("DB_NAME")
	if name == "" {
		return c.GetDB().Database
	}

	return name
}

func (c *config) GetDatabaseType() string {
	t := GetEnv("DB_TYPE")
	if t == "" {
		return c.GetDB().Type
	}

	return t
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Println("Environment variable not set: ", key)

		return ""
	}

	return value
}
