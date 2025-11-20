package main

import (
	"net/http"
	"sync"
	"time"

	"github.com/tomasen/realip"
	"golang.org/x/time/rate"
)

func (app *application) rateLimiter(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		mu      sync.Mutex
		clients = make(map[string]*client) // IP -> Client struct which contains limiter
	)

	// periodically removing clients if inactive
	go func() {
		for {
			time.Sleep(time.Minute)

			mu.Lock()

			for ip, client := range clients {
				if time.Since(client.lastSeen) > 3*time.Minute {
					delete(clients, ip)
				}
			}

			mu.Unlock()
		}
	}()

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if app.cfg.limiter.enabled {
			ip := realip.FromRequest(r) // fetching the client's IP

			mu.Lock()

			if _, found := clients[ip]; !found {
				clients[ip] = &client{limiter: rate.NewLimiter(rate.Limit(app.cfg.limiter.rps), app.cfg.limiter.burst)}
			}

			clients[ip].lastSeen = time.Now()

			if !clients[ip].limiter.Allow() {
				mu.Unlock() // gotta unlock mutex before returning
				app.rateLimitExceededResponse(w, r)
				return
			}

			mu.Unlock()
		}
		next.ServeHTTP(w, r)
	})
}
