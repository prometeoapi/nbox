package health

import (
	"fmt"
	"math/rand"
	"nbox/internal/entrypoints/api/response"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"
)

var (
	once     sync.Once
	instance *Health
)

type healthOut struct {
	*Health
	Uptime string `json:"uptime"`
}

type Health struct {
	StartedAt time.Time `json:"startedAt"`
	Service   string    `json:"service"`
	Hostname  string    `json:"hostname"`
}

func NewHealthy() *Health {
	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1).Intn(99)

	hostname, err := os.Hostname()
	if hostname == "" || err != nil {
		hostname = "localhost"
	}

	once.Do(func() {
		instance = &Health{
			StartedAt: time.Now(),
			Hostname:  fmt.Sprintf("%s-%d", hostname, r1),
			Service:   "nbox",
		}
	})
	return instance
}

func (u Health) Uptime() string {
	return time.Since(u.StartedAt).String()
}

func (u Health) Healthy(endpoint string) func(http.Handler) http.Handler {
	f := func(h http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			if (r.Method == "GET" || r.Method == "HEAD") && strings.EqualFold(r.URL.Path, endpoint) {
				response.Success(w, r, healthOut{
					&u, u.Uptime(),
				})
				return
			}
			h.ServeHTTP(w, r)
		}
		return http.HandlerFunc(fn)
	}
	return f
}
