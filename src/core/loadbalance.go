package core

import (
	"math/rand"
	"fmt"
	"net/url"
	"sync"
	"github.com/JrSchmidtt/loadbalance-api/src/config"
)

type LoadBalancer struct {
	Services []config.Service
	mu 	 sync.Mutex
	next int
}

func NewLoadBalancer() (*LoadBalancer, error) {
	config := config.ReadConfig("config.yml")
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
	currentWeight := 0
	for i := 0; i < len(lb.Services); i++ {
		service := lb.Services[(lb.next+i) % len(lb.Services)] // calculate next index
		for _, weight := range service.TrafficWeights {
			currentWeight += weight.Percentage
			if randNumber < currentWeight {
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