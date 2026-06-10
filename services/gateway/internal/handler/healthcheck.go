package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
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
	resultChan := make(chan result, len(services)+1)

	var wg sync.WaitGroup
	for service, conn := range services {
		wg.Go(func() {
			_, err := conn.Check(ctx, &grpc_health_v1.HealthCheckRequest{Service: ""})
			resultChan <- result{service: service, status: err == nil}
		})
	}

	wg.Go(func() {
		consumers, err := h.getRabbitConsumersAmount()
		resultChan <- result{service: "notification", status: err == nil && consumers == 1}
		if err != nil {
			h.log.Warn("error while requesting rabbitMQ", "err", err.Error())
		}
	})

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

	responses["gateway"] = "SERVING"

	data := jsonutil.Envelope{
		"timestamp": time.Now(),
		"services":  responses,
	}

	jsonutil.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) getRabbitConsumersAmount() (int, error) {
	vhost := "%2f" // Дефолтный vhost "/"
	queue := "email"

	url := fmt.Sprintf("http://%s/api/queues/%s/%s", h.cfg.Rabbit.URL, vhost, queue)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return 0, err
	}

	req.SetBasicAuth(h.cfg.Rabbit.User, h.cfg.Rabbit.Pass)

	client := &http.Client{Timeout: 1 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return 0, fmt.Errorf("rabbit api unavailable: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return 0, fmt.Errorf("queue '%s' not found", queue)
	}
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return 0, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(body))
	}
	var info struct {
		Consumers int `json:"consumers"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&info); err != nil {
		return 0, err
	}

	return info.Consumers, nil
}
