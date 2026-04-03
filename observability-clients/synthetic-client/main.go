package main

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"sync"
	"time"
)

var (
	running   bool
	frequency int = 5 // Initial default is 5 seconds
	targetUrl string = "http://api-gateway:8080" // Initial default target
	stats     = map[string]int{"success": 0, "error": 0}
	mu        sync.Mutex
	ctxCancel context.CancelFunc
)

func indexHandler(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("templates/index.html")
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	mu.Lock()
	defer mu.Unlock()
	data := map[string]interface{}{
		"Running":   running,
		"Frequency": frequency,
		"TargetUrl": targetUrl,
		"Success":   stats["success"],
		"Error":     stats["error"],
	}
	tmpl.Execute(w, data)
}

func apiConfigHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		var req struct {
			Action    string `json:"action"`
			Frequency int    `json:"frequency"`
			TargetUrl string `json:"targetUrl"`
		}
		json.NewDecoder(r.Body).Decode(&req)

		mu.Lock()
		if req.Action == "start" && !running {
			running = true
			if req.Frequency > 0 {
				frequency = req.Frequency
			}
			if req.TargetUrl != "" {
				targetUrl = req.TargetUrl
			}
			var ctx context.Context
			ctx, ctxCancel = context.WithCancel(context.Background())
			go runLoadLoop(ctx)
		} else if req.Action == "stop" && running {
			running = false
			if ctxCancel != nil {
				ctxCancel()
			}
		} else if req.Action == "update" {
			if req.Frequency > 0 {
				frequency = req.Frequency
			}
			if req.TargetUrl != "" {
				targetUrl = req.TargetUrl
			}
		}
		mu.Unlock()
	}

	mu.Lock()
	defer mu.Unlock()
	json.NewEncoder(w).Encode(map[string]interface{}{
		"running":   running,
		"frequency": frequency,
		"targetUrl": targetUrl,
		"success":   stats["success"],
		"error":     stats["error"],
	})
}

func runLoadLoop(ctx context.Context) {
	for {
		select {
		case <-ctx.Done():
			return
		default:
			mu.Lock()
			freq := frequency
			target := targetUrl
			mu.Unlock()

			// Launch headless browser journey
			err := executeBrowserJourney(target)
			
			mu.Lock()
			if err != nil {
				log.Println("Journey Error:", err)
				stats["error"]++
			} else {
				stats["success"]++
			}
			mu.Unlock()

			select {
			case <-time.After(time.Duration(freq) * time.Second):
			case <-ctx.Done():
				return
			}
		}
	}
}

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/api/config", apiConfigHandler)
	fmt.Println("Observability Synthetic Client starting on :8090")
	log.Fatal(http.ListenAndServe(":8090", nil))
}
