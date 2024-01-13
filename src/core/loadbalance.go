package core

import (
	"fmt"
	"log"
	"net/url"

	"github.com/JrSchmidtt/loadbalance-api/src/config"
)

type LoadBalancer struct {
	Services []config.Service
}

func NewLoadBalancer() (*LoadBalancer, error) {
	config := config.ReadConfig("config.yml")
	return &LoadBalancer{
		Services: config.Services,
	}, nil
}

func (lb *LoadBalancer) ChooseService() *url.URL {
	if len(lb.Services) > 0 {
		for _, service := range lb.Services {
			for _, trafficWeight := range service.TrafficWeights {
				fmt.Printf("Percentage: %d\n", trafficWeight.Percentage)
			}
			fmt.Printf("Service: %s\n", service.Name)
		}
		serviceURL, err := url.Parse(lb.Services[0].URL)
		if err != nil {
			log.Printf("Erro ao fazer o parse da URL do serviço: %v", err)
			return nil
		}
		return serviceURL
	}
	return nil
}

func (lb *LoadBalancer) GetServices() []config.Service {
	return lb.Services
}

func (lb *LoadBalancer) AddService(service config.Service) {
	lb.Services = append(lb.Services, service)
}

func (lb *LoadBalancer) RemoveService(service config.Service) {
	for i, s := range lb.Services {
		if s.Name == service.Name {
			lb.Services = append(lb.Services[:i], lb.Services[i+1:]...)
		}
	}
}

func (lb *LoadBalancer) UpdateService(service config.Service) {
	for i, s := range lb.Services {
		if s.Name == service.Name {
			lb.Services[i] = service
		}
	}
}