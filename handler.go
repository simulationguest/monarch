package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Selector struct {
	handlers map[string]http.Handler
}

func NewSelector(c Config) Selector {
	handlers := make(map[string]http.Handler)

	for site, config := range c {
		if _, ok := handlers[site]; ok {
			log.Fatalf("There is already a handler registered for %s\n", site)
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
				log.Fatalf("Failed to parse url for reverse proxy (host: \"%s\", to: \"%s\"): %s\n", site, to, err)
			}
			proxy := httputil.NewSingleHostReverseProxy(target)
			handlers[site] = proxy
		default:
			log.Fatalf("Unrecognized action (host: %s): %s\n", site, config.Action)
		}
	}

	return Selector{
		handlers: handlers,
	}
}

func (s *Selector) ServeString() {
	for site := range s.handlers {
		fmt.Println("Serving site", site)
	}
	s.Serve()
}

func (s *Selector) Serve() {
	for site, handler := range s.handlers {
		site := site
		handler := handler
		server := http.Server{
			Addr:    site,
			Handler: handler,
		}
		go log.Println(server.ListenAndServe())

	}
	select {}
}
