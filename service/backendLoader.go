package service

import (
	"net/http/httputil"
	"net/url"
	"sync"
)

type Backend struct {
	URL          *url.URL
	Alive        bool
	Mux          sync.Mutex
	ReverseProxy *httputil.ReverseProxy
}

func Backends(urls []string) ([]*Backend, error) {
	var backends []*Backend
	for _, raw := range urls {
		parsed, err := url.Parse(raw)
		if err != nil {
			return nil, err
		}
		proxy := reverseProxy(parsed)
		backends = append(backends, &Backend{
			URL:          parsed,
			Alive:        true,
			ReverseProxy: proxy,
		})
	}
	return backends, nil
}
