package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

// Contributor represents a SocioVault contributor
type Contributor struct {
	ID        string    `json:"id"`
	Email     string    `json:"email"`
	Amount    float64   `json:"amount"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"createdAt"`
}

// Request represents a help request
type Request struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Email       string    `json:"email"`
	Story       string    `json:"story"`
	VideoURL    string    `json:"videoUrl"`
	Amount      float64   `json:"amount"`
	Verified    bool      `json:"verified"`
	Votes       int       `json:"votes"`
	CreatedAt   time.Time `json:"createdAt"`
}

// Stats represents overall statistics
type Stats struct {
	TotalBalance        float64 `json:"totalBalance"`
	DistributedPercent  int     `json:"distributedPercent"`
	StoriesFunded       int     `json:"storiesFunded"`
	TotalContributors   int     `json:"totalContributors"`
	ActiveRequests      int     `json:"activeRequests"`
	DailyContributions  float64 `json:"dailyContributions"`
}

var (
	contributors []Contributor
	requests     []Request
	contributorMu sync.Mutex
	requestMu    sync.Mutex
	dataDir      = "./data"
)

func init() {
	// Create data directory if it doesn't exist
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Printf("Failed to create data directory: %v", err)
	}

	// Load existing data
	loadContributors()
	loadRequests()
}

func loadContributors() {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/contributors.json", dataDir))
	if err != nil {
		// File doesn't exist yet, start with empty
		contributors = []Contributor{}
		return
	}
	json.Unmarshal(data, &contributors)
}

func loadRequests() {
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/requests.json", dataDir))
	if err != nil {
		requests = []Request{}
		return
	}
	json.Unmarshal(data, &requests)
}

func saveContributors() {
	data, _ := json.MarshalIndent(contributors, "", "  ")
	ioutil.WriteFile(fmt.Sprintf("%s/contributors.json", dataDir), data, 0644)
}

func saveRequests() {
	data, _ := json.MarshalIndent(requests, "", "  ")
	ioutil.WriteFile(fmt.Sprintf("%s/requests.json", dataDir), data, 0644)
}

func generateID() string {
	return fmt.Sprintf("%d", time.Now().UnixNano())
}

// API Handlers

func handleContribute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		var contrib Contributor
		if err := json.NewDecoder(r.Body).Decode(&contrib); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		contrib.ID = generateID()
		contrib.CreatedAt = time.Now()

		contributorMu.Lock()
		contributors = append(contributors, contrib)
		saveContributors()
		contributorMu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(contrib)
		return
	}

	if r.Method == http.MethodGet {
		contributorMu.Lock()
		defer contributorMu.Unlock()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(contributors)
		return
	}
}

func handleRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method == http.MethodPost {
		var req Request
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			http.Error(w, "Invalid request", http.StatusBadRequest)
			return
		}

		req.ID = generateID()
		req.CreatedAt = time.Now()
		req.Verified = false
		req.Votes = 0

		requestMu.Lock()
		requests = append(requests, req)
		saveRequests()
		requestMu.Unlock()

		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(req)
		return
	}

	if r.Method == http.MethodGet {
		requestMu.Lock()
		defer requestMu.Unlock()
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(requests)
		return
	}
}

func handleStats(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	contributorMu.Lock()
	totalContributors := len(contributors)
	var totalBalance float64
	var dailyContributions float64
	today := time.Now().Format("2006-01-02")

	for _, c := range contributors {
		totalBalance += c.Amount
		if c.CreatedAt.Format("2006-01-02") == today {
			dailyContributions += c.Amount
		}
	}
	contributorMu.Unlock()

	requestMu.Lock()
	activeRequests := len(requests)
	distributedPercent := 92 // Default for now
	storiesFunded := 311      // Default for now

	requestMu.Unlock()

	stats := Stats{
		TotalBalance:       48320 + totalBalance, // Base + new contributions
		DistributedPercent: distributedPercent,
		StoriesFunded:      storiesFunded + len(requests),
		TotalContributors:  2847 + totalContributors,
		ActiveRequests:     activeRequests,
		DailyContributions: dailyContributions,
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(stats)
}

func handleVote(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var voteData struct {
		RequestID string `json:"requestId"`
	}

	if err := json.NewDecoder(r.Body).Decode(&voteData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	requestMu.Lock()
	defer requestMu.Unlock()

	for i, req := range requests {
		if req.ID == voteData.RequestID {
			requests[i].Votes++
			saveRequests()
			w.WriteHeader(http.StatusOK)
			json.NewEncoder(w).Encode(requests[i])
			return
		}
	}

	http.Error(w, "Request not found", http.StatusNotFound)
}

func handleSubscribe(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	if r.Method == http.MethodOptions {
		w.WriteHeader(http.StatusOK)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var subData struct {
		Email string `json:"email"`
	}

	if err := json.NewDecoder(r.Body).Decode(&subData); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Save subscription
	subs, _ := ioutil.ReadFile(fmt.Sprintf("%s/subscribers.json", dataDir))
	var subscribers []map[string]interface{}
	json.Unmarshal(subs, &subscribers)

	subscribers = append(subscribers, map[string]interface{}{
		"email":     subData.Email,
		"timestamp": time.Now(),
	})

	data, _ := json.MarshalIndent(subscribers, "", "  ")
	ioutil.WriteFile(fmt.Sprintf("%s/subscribers.json", dataDir), data, 0644)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "subscribed",
		"email":  subData.Email,
	})
}

func main() {
	// API Routes
	http.HandleFunc("/api/contribute", handleContribute)
	http.HandleFunc("/api/requests", handleRequest)
	http.HandleFunc("/api/stats", handleStats)
	http.HandleFunc("/api/vote", handleVote)
	http.HandleFunc("/api/subscribe", handleSubscribe)

	// Static files
	http.Handle("/", http.FileServer(http.Dir(".")))

	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
	}

	fmt.Printf("SocioVault server listening on :%s\n", port)
	log.Fatal(http.ListenAndServe(":"+port, nil))
}
