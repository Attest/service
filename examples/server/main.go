package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"github.com/Attest/service"
)

func main() {
	addr := ":4444"
	if env := os.Getenv("EXAMPLE_HTTP_ADDR"); env != "" {
		addr = env
	}

	g := service.NewSignals(os.Interrupt, os.Kill)
	g.Register(server(addr, g))

	log.Fatal(g.Start())
}

func server(addr string, g *service.Group) func(ctx context.Context) error {
	return func(ctx context.Context) error {
		srv := http.Server{Addr: addr, Handler: handler(g)}

		go func() {
			<-ctx.Done()
			_ = srv.Close()
		}()

		return srv.ListenAndServe()
	}
}

func handler(r *service.Group) http.HandlerFunc {
	f := func(w http.ResponseWriter, req *http.Request) {
		log.Println(req.URL.Path)
		if req.URL.Path == "/live" {
			r.Liveness()(w, req)
			return
		}

		if req.URL.Path == "/ready" {
			r.Readiness()(w, req)
			return
		}

		w.WriteHeader(http.StatusNotFound)
	}

	return http.HandlerFunc(f)
}
