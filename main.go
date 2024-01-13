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
		fmt.Println(err)
		return
	}
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		serviceURL := loadBalancer.ChooseService()
		fmt.Printf("Chosen Service: %s\n", serviceURL)
		reverseProxy := httputil.NewSingleHostReverseProxy(serviceURL)
		reverseProxy.ServeHTTP(w, r)
	})
	fmt.Println("Load Balancer started at :8080")
	http.ListenAndServe(":8080", nil)
}