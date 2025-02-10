package config

import (
	"encoding/json"
	"fmt"
	"os"
)

type config struct {
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

// NewConfig creates a new Config instance with default values
func NewConfig() config {
	c := &config{}

	// load the configuratiion file
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
