package config

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/yaml.v3"
)

// ReadConfig reads the config file and returns a Config struct.
func ReadConfig(configFile string) Config {
	data, err := os.ReadFile(configFile)
	if err != nil {
		log.Fatal("Error reading config file", err)
	}
	var config Config
	err = yaml.Unmarshal(data, &config)
	if err != nil {
		log.Fatal("Error parsing config file", err)
	}
	err = config.Validate()
	if err != nil {
		log.Fatal("Error validating config file", err)
	}
	// check if traffic weight sum is 100 for each service.
	sum := 0
	for _, service := range config.Services {
		for _, trafficWeight := range service.TrafficWeights {
			sum += trafficWeight.Percentage
		}
		fmt.Printf("Service: %s\n", service.Name)
	}
	if sum != 100 {
		log.Fatal("Traffic weight sum must be 100")
	}
	return config
}