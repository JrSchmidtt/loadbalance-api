package core

import (
	"log"
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

func (lb *LoadBalancer) ChooseService() (*url.URL) {
	lb.mu.Lock()
	defer lb.mu.Unlock()

	service := lb.Services[lb.next]
	lb.next = (lb.next + 1) % len(lb.Services)

	for _, wight := range service.TrafficWeights {
		if wight.Percentage > 0 {
			serviceURL, err := url.Parse(lb.Services[0].URL)
			if err != nil {
				log.Printf("Error parsing URL: %s\n", err)
				return nil
			}
			return serviceURL
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