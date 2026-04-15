package handler

import (
	"context"
	"net/http"
	"sync"
	"time"

	"github.com/ezcnrmn/vaito/services/gateway/internal/lib/jsonutil"
	"google.golang.org/grpc/health/grpc_health_v1"
)

// Healthcheck godoc
//
//	@summary	healthcheck
//	@tags		utility
//	@produce	json
//	@success	200
//	@router		/healthcheck [get]
func (h *Handler) Healthcheck(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	services := map[string]grpc_health_v1.HealthClient{
		"user":    h.health.user,
		"listing": h.health.listing,
	}
	type result struct {
		service string
		status  bool
	}
	resultChan := make(chan result, len(services))

	var wg sync.WaitGroup
	for service, conn := range services {
		wg.Go(func() {
			_, err := conn.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
			resultChan <- result{service: service, status: err == nil}
		})
	}

	go func() {
		wg.Wait()
		close(resultChan)
	}()

	responses := make(map[string]any)
	for r := range resultChan {
		if r.status {
			responses[r.service] = "SERVING"
		} else {
			responses[r.service] = "DOWN"
		}
	}

	data := jsonutil.Envelope{
		"status":    "available",
		"timestamp": time.Now(),
		"services":  responses,
	}

	jsonutil.WriteJSON(w, http.StatusOK, data)
}
