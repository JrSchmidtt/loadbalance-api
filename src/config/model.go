package config

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

type Service struct {
	Name           string          `yaml:"name" validate:"required"`
	URL            string          `yaml:"url" validate:"required"`
	Path           string          `yaml:"path" validate:"required"`
	TrafficWeights []TrafficWeight `yaml:"traffic_weight" validate:"required"`
}

type TrafficWeight struct {
	Percentage int `yaml:"percentage" validate:"required"`
}

type Config struct {
	Services []Service `yaml:"services" validate:"required"`
}

var validate *validator.Validate = validator.New()

func (cfg Config) Validate() error {
	fmt.Println(cfg.Services[1].Name)
	fmt.Println("Validating config file...")
	err := validate.Struct(cfg)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}