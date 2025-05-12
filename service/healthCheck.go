package service

import (
	"log"
	"net/http"
	"time"
)

const healthCheckInterval = 10 * time.Second

func (b *Backend) SetAlive(alive bool) {
	b.Mux.Lock()
	defer b.Mux.Unlock()
	b.Alive = alive
}

func (b *Backend) IsAlive() bool {
	b.Mux.Lock()
	defer b.Mux.Unlock()
	return b.Alive
}

func checkBackendHealth(b *Backend) {
	resp, err := http.Get(b.URL.String())
	if err != nil || resp.StatusCode >= 500 {
		b.SetAlive(false)
		log.Printf("[HealthCheck] %s is DOWN\n", b.URL)
		return
	}
	b.SetAlive(true)
	log.Printf("[HealthCheck] %s is UP\n", b.URL)
}

func StartHealthCheck(pool *ServerPool) {
	ticker := time.NewTicker(healthCheckInterval)
	go func() {
		for {
			<-ticker.C
			log.Println("[HealthCheck] Running...")
			for _, backend := range pool.Servers {
				go checkBackendHealth(backend)
			}
		}
	}()
}
