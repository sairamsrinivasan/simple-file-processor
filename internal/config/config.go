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
	Redis   redis    `json:"redis"`
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

type database struct {
	Type     string `json:"type"`
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Username string `json:"username"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type redis struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	Database int    `json:"database"`
}

type Config interface {
	Port() int
	GetRoutes() []routes
	GetDB() database
	DatabaseUsername() string
	DatabasePassword() string
	DatabaseHost() string
	DatabasePort() int
	DatabaseName() string
	DatabaseType() string
	ConnectionString() string
	RedisAddress() string
	RedisDB() int
	RedisURL() string
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
func (c *config) Port() int {
	p := EnvOrDefault("APP_PORT", strconv.Itoa(c.Service.Port))
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

func (c *config) DatabasePassword() string {
	return EnvOrDefault("FILE_DATABASE_PASSWORD", c.DB.Password)
}

var uname string

func (c *config) DatabaseUsername() string {
	return EnvOrDefault("FILE_DATABASE_USERNAME", c.DB.Username)
}

func (c *config) DatabaseHost() string {
	return EnvOrDefault("DB_HOST", c.DB.Host)
}

func (c *config) DatabasePort() int {
	p := EnvOrDefault("DB_PORT", strconv.Itoa(c.DB.Port))
	port, _ := strconv.Atoi(p)
	return port
}

func (c *config) DatabaseName() string {
	return EnvOrDefault("DB_NAME", c.DB.Name)
}

// returns the database type eg. postgres, mysql, etc
func (c *config) DatabaseType() string {
	return EnvOrDefault("DB_TYPE", c.DB.Type)
}

// returns the connection string for the database
func (c *config) ConnectionString() string {
	// Construct the database connection string here'
	str := fmt.Sprintf("%s:%s@%s:%d/%s?sslmode=disable", c.DatabaseUsername(), c.DatabasePassword(), c.DatabaseHost(), c.DatabasePort(), c.DatabaseName())
	return c.DatabaseType() + "://" + str
}

func (c *config) RedisURL() string {
	// Construct the redis connection string here
	str := fmt.Sprintf("redis://%s:%d/%d", c.Redis.Host, c.Redis.Port, c.Redis.Database)
	return str
}

func (c *config) RedisAddress() string {
	// Construct the address for the redis server
	// includes the host and port
	return fmt.Sprintf("%s:%d", c.Redis.Host, c.Redis.Port)
}

func (c *config) RedisDB() int {
	// returns the redis database number
	return c.Redis.Database
}

func EnvOrDefault(key string, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		fmt.Printf("Environment variable %s not set, using default value.\n", key)
		return defaultValue
	}

	return value
}
