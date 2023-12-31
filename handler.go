package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Selector struct {
	handlers map[Site]http.Handler
}

func NewSelector(c Config) Selector {
	handlers := make(map[Site]http.Handler)

	for host, config := range c {
		u, err := url.Parse(host)
		if err != nil {
			log.Fatalf("Failed to parse host \"%s\": %s\n", host, err)
		}

		site := NewSite(u)

		if _, ok := handlers[site]; ok {
			log.Fatalf("There is already a handler registered for %s\n", u)
		}

		switch config.Action {
		case ActionRespond:
			response := config.Data["with"].(string)
			handlers[site] = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprintf(w, "%s", response)
			})
		case ActionReverseProxy:
			to := config.Data["to"].(string)
			target, err := url.Parse(to)
			if err != nil {
				log.Fatalf("Failed to parse url for reverse proxy (host: \"%s\", to: \"%s\"): %s\n", host, to, err)
			}
			proxy := httputil.NewSingleHostReverseProxy(target)
			handlers[site] = proxy
		default:
			log.Fatalf("Unrecognized action (host: %s): %s\n", host, config.Action)
		}
	}

	return Selector{
		handlers: handlers,
	}
}

func (s *Selector) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	site := NewSite(r.URL)
	handler, ok := s.handlers[site]
	if !ok {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	handler.ServeHTTP(w, r)
}
