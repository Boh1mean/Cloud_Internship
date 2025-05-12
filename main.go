package main

import (
	"CloudInternship/configs"
	"CloudInternship/service"
	"log"
	"net/http"
)

var serverPool *service.ServerPool

func main() {
	cfg, err := configs.LoadServer("servers.json")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	backends, err := service.Backends(cfg.Servers)
	if err != nil {
		log.Fatalf("Failed to initialize backends: %v", err)
	}

	serverPool = service.NewServerPool(backends)
	service.StartHealthCheck(serverPool)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		server := serverPool.NextServer()
		if server == nil {
			http.Error(w, "No healthy server available", http.StatusServiceUnavailable)
			return
		}

		w.Header().Add("X-Forwarded-Server", server.URL.String())
		server.ReverseProxy.ServeHTTP(w, r) // Assuming ReverseProxy is a field, not a method
	})

	log.Println("Starting load balancer on port", cfg.Port)
	err = http.ListenAndServe(cfg.Port, nil)
	if err != nil {
		log.Fatalf("Error starting load balancer: %s\n", err.Error())
	}
}
