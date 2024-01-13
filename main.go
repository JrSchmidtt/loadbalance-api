package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"

	"github.com/JrSchmidtt/loadbalance-api/src/core"
)

func main() {
	loadBalancer, err := core.NewLoadBalancer()
	if err != nil {
		panic(err)
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serviceURL := loadBalancer.ChooseService()
		if serviceURL == nil {
			fmt.Println("Error choosing service")
			return
		}
		reverseProxy := httputil.NewSingleHostReverseProxy(serviceURL)
		reverseProxy.ServeHTTP(w, r)
	})
	fmt.Println("Load Balancer started at :8080")
	http.ListenAndServe(":8080", nil)
}