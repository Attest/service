package service

import "net/http"

func (g Group) Liveness() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !g.started() {
			w.WriteHeader(http.StatusServiceUnavailable)
		}
	}
}

func (g Group) Readiness() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		if !g.started() || !g.setup() {
			w.WriteHeader(http.StatusServiceUnavailable)
		}

		g.processLock.RLock()
		defer g.processLock.RUnlock()

		for i := range g.processes {
			r := g.processes[i]

			if !r.Healthy() {
				w.WriteHeader(http.StatusServiceUnavailable)
				return
			}
		}
	}
}
