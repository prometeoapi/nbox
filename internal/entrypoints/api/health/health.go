package health

import (
	"fmt"
	"math/rand"
	"nbox/internal/entrypoints/api/response"
	"net/http"
	"os"
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

// Healthy check
// @Summary health
// @Description status format json
// @Tags status
// @Produce json
// @Success 200 {object} health.Health{}
// @Router /health [get]
// @Accept       json
// @Produce      json
func (u Health) Healthy(w http.ResponseWriter, r *http.Request) {
	response.Success(w, r, healthOut{
		&u, u.Uptime(),
	})
}
