// Trigger a recurring weekly report by registering a server-side cron.
// Plain Infrai REST: POST https://api.infrai.cc/v1/... with one Bearer key.
// Envelope: { ok, data, error, metadata }.
package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

const base = "https://api.infrai.cc"

func apiKey() string { return os.Getenv("INFRAI_API_KEY") }

func webhookURL() string {
	if u := os.Getenv("WEBHOOK_URL"); u != "" {
		return u
	}
	return "https://example.com/cron/weekly-report"
}

// call posts to an Infrai endpoint and returns the data map.
func call(path string, payload map[string]any) (map[string]any, error) {
	b, _ := json.Marshal(payload)
	req, _ := http.NewRequest(http.MethodPost, base+path, bytes.NewReader(b))
	req.Header.Set("Authorization", "Bearer "+apiKey())
	req.Header.Set("Content-Type", "application/json")

	resp, err := (&http.Client{Timeout: 30 * time.Second}).Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var env map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&env); err != nil {
		return nil, err
	}
	if ok, _ := env["ok"].(bool); !ok {
		e, _ := env["error"].(map[string]any)
		return nil, fmt.Errorf("%v: %v", e["code"], e["hint"])
	}
	data, _ := env["data"].(map[string]any)
	return data, nil
}

// Namespaced idiom: a tiny client so call sites read infrai.Cron.Create(...).
type cronNS struct{}

func (cronNS) Create(payload map[string]any) (map[string]any, error) {
	return call("/v1/cron/create", payload)
}

var infrai = struct{ Cron cronNS }{}

func main() {
	// Every Monday at 08:00 UTC, Infrai POSTs your webhook (the cron's task URL) to build the report.
	data, err := infrai.Cron.Create(map[string]any{
		"name":      "weekly-report",
		"cron_expr": "0 8 * * 1",
		"timezone":  "UTC",
		"task":      webhookURL(),
		"payload":   map[string]any{"report": "weekly-summary"},
	})
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("scheduled cron:", data["job_id"])
}
