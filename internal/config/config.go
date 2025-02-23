package config

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
)

type config struct {
	DB      database `json:"database"`
	Service service  `json:"service"`
	Routes  []routes `json:"routes"`
}

type service struct {
	Name    string `json:"name"`
	Version string `json:"version"`
	Port    int    `json:"port"`
}

type routes struct {
	Path    string `json:"path"`
	Handler string `json:"handler"`
	Method  string `json:"method"`
}

type Config interface {
	GetPort() int
	GetRoutes() []routes
	GetDB() database
	GetDatabaseUsername() string
	GetDatabasePassword() string
	GetDatabaseHost() string
	GetDatabasePort() int
	GetDatabaseName() string
	GetDatabaseType() string
	GetConnectionString() string
}

type database struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

// NewConfig creates a new Config instance with default values
func NewConfig() Config {
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
	return c
}

// returns the port from the configuration
func (c *config) GetPort() int {
	p := GetEnv("APP_PORT")
	if p == "" {
		return c.Service.Port
	}

	port, _ := strconv.Atoi(p)
	return port
}

// returns the routes from the configuration
func (c *config) GetRoutes() []routes {
	return c.Routes
}

func (c *config) GetDB() database {
	return c.DB
}

func (c *config) GetDatabasePassword() string {
	p := GetEnv("FILE_DATABASE_PASSWORD")
	if p == "" {
		return c.GetDB().Password
	}

	return p
}

var uname string

func (c *config) GetDatabaseUsername() string {
	uname := GetEnv("FILE_DATABASE_USERNAME")
	if uname == "" {
		return c.GetDB().Username
	}

	return uname
}

func (c *config) GetDatabaseHost() string {
	h := GetEnv("DB_HOST")
	if h == "" {
		return c.GetDB().Host
	}

	return h
}

func (c *config) GetDatabasePort() int {
	p := GetEnv("DB_PORT")
	if p == "" {
		return c.GetDB().Port
	}

	port, _ := strconv.Atoi(p)
	return port
}

func (c *config) GetDatabaseName() string {
	name := GetEnv("DB_NAME")
	if name == "" {
		return c.GetDB().Name
	}

	return name
}

// returns the database type eg. postgres, mysql, etc
func (c *config) GetDatabaseType() string {
	t := GetEnv("DB_TYPE")
	if t == "" {
		return c.GetDB().Type
	}

	return t
}

// returns the connection string for the database
func (c *config) GetConnectionString() string {
	// Construct the database connection string here'
	str := fmt.Sprintf("%s:%s@%s:%d/%s?sslmode=disable", c.GetDatabaseUsername(), c.GetDatabasePassword(), c.GetDatabaseHost(), c.GetDatabasePort(), c.GetDatabaseName())

	return c.GetDatabaseType() + "://" + str
}

func GetEnv(key string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Println("Environment variable not set: ", key)

		return ""
	}

	return value
}
