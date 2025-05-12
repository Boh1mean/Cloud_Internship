package service

import (
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func reverseProxy(backend *url.URL) *httputil.ReverseProxy {
	proxy := httputil.NewSingleHostReverseProxy(backend)

	proxy.ErrorHandler = func(w http.ResponseWriter, r *http.Request, err error) {
		log.Printf("error proxying to %s: %v", backend.String(), err)
		w.Header().Set("X-Load-Balancer", "GoBalancer")
		w.WriteHeader(http.StatusBadGateway)
		if _, writeErr := w.Write([]byte("Service currently unavailable")); writeErr != nil {
			log.Printf("failed to write error response: %v", writeErr)
		}
	}

	return proxy
}
