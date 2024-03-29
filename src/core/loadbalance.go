package core

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"github.com/JrSchmidtt/loadbalance-api/src/config"
)

type LoadBalancer struct {
	Services []config.Service
	mu 	 sync.Mutex
	next int
	currentWeight int
}

func NewLoadBalancer() (*LoadBalancer, error) {
	config := config.ReadConfig("config.yml")
	for i := range config.Services {
		err := pingService(config.Services[i].URL)
		if err != nil {
			fmt.Printf("service %s is not available\n", config.Services[i].Name)
			config.Services = append(config.Services[:i], config.Services[i+1:]...)
		}
	}
	if len(config.Services) == 0 {
		return nil, fmt.Errorf("no services available")
	}
	return &LoadBalancer{
		Services: config.Services,
	}, nil
}

func (lb *LoadBalancer) ChooseService() *url.URL {
	lb.mu.Lock()
	defer lb.mu.Unlock()
	totalWeight := 0
	for _, service := range lb.Services {
		for _, weight := range service.TrafficWeights {
			totalWeight += weight.Percentage
		}
	}
	randNumber := rand.Intn(totalWeight)
	for i := 0; i < len(lb.Services); i++ {
		service := lb.Services[(lb.next+i) % len(lb.Services)] // calculate next index
		for _, weight := range service.TrafficWeights {
			lb.currentWeight += weight.Percentage
			if randNumber < lb.currentWeight {
				lb.next = (lb.next + i + 1) % len(lb.Services)
				serviceURL, err := url.Parse(service.URL)
				if err != nil {
					fmt.Printf("Error parsing service URL: %s\n", err)
					return nil
				}
				return serviceURL
			}
		}
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

func pingService(serviceURL string) error {
	_, err := http.Get(serviceURL)
	if err != nil {
		return err
	}
	return nil
}