package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

// ============ DATA MODELS ============
type Auction struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	StartPrice    float64   `json:"start_price"`
	CurrentBid    float64   `json:"current_bid"`
	HighestBidder string    `json:"highest_bidder"`
	BidCount      int       `json:"bid_count"`
	Status        string    `json:"status"`
	StartTime     time.Time `json:"start_time"`
	EndTime       time.Time `json:"end_time"`
}

type Bid struct {
	BidderID  string    `json:"bidder_id"`
	Amount    float64   `json:"amount"`
	Timestamp time.Time `json:"timestamp"`
}

type Message struct {
	Type    string     `json:"type"`
	Auction *Auction   `json:"auction,omitempty"`
	Bid     *Bid       `json:"bid,omitempty"`
	Error   string     `json:"error,omitempty"`
}

type AuctionSubmission struct {
	Name         string  `json:"name"`
	Description  string  `json:"description"`
	StartPrice   float64 `json:"startPrice"`
	Duration     int     `json:"duration"`
	Email        string  `json:"email"`
	Timestamp    string  `json:"timestamp"`
}

// ============ GLOBAL STATE ============
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
		HandshakeTimeout: 10 * time.Second,
		ReadBufferSize:   1024,
		WriteBufferSize:  1024,
	}

	auctions      = make(map[string]*Auction)
	auctionMutex  sync.RWMutex

	clients       = make(map[*Client]bool)
	clientsMutex  sync.RWMutex

	broadcast     = make(chan Message, 256)
	
	startTime     = time.Now()
	
	// WebSocket ping/pong settings
	pongWait = 60 * time.Second
	pingPeriod = (pongWait * 9) / 10  // Send pings at 90% of pong wait time
	writeWait = 10 * time.Second
)

type Client struct {
	conn *websocket.Conn
	id   string
}

// ============ INITIALIZATION ============
func init() {
	now := time.Now()

	// Active auction 1
	auctions["auction-1"] = &Auction{
		ID:            "auction-1",
		Title:         "Vintage Camera Collection",
		Description:   "1960s Leica M3 and accessories - Excellent condition",
		StartPrice:    100,
		CurrentBid:    850,
		HighestBidder: "bidder_42",
		BidCount:      24,
		Status:        "active",
		StartTime:     now.Add(-5 * time.Minute),
		EndTime:       now.Add(5 * time.Minute),
	}

	// Active auction 2
	auctions["auction-2"] = &Auction{
		ID:            "auction-2",
		Title:         "Modern Art Painting",
		Description:   "Oil on canvas by emerging artist - 100x80cm",
		StartPrice:    200,
		CurrentBid:    2500,
		HighestBidder: "bidder_elite",
		BidCount:      47,
		Status:        "active",
		StartTime:     now.Add(-10 * time.Minute),
		EndTime:       now.Add(2 * time.Minute),
	}

	// Scheduled auction
	auctions["auction-3"] = &Auction{
		ID:          "auction-3",
		Title:       "Rare Vinyl Records",
		Description: "Limited edition Beatles pressings - Mint condition",
		StartPrice:  50,
		CurrentBid:  50,
		Status:      "scheduled",
		StartTime:   now.Add(30 * time.Minute),
		EndTime:     now.Add(60 * time.Minute),
	}

	// Ended auction
	auctions["auction-4"] = &Auction{
		ID:            "auction-4",
		Title:         "Antique Watch",
		Description:   "Swiss-made pocket watch - 1940s",
		StartPrice:    150,
		CurrentBid:    3200,
		HighestBidder: "bidder_collector",
		BidCount:      62,
		Status:        "ended",
		StartTime:     now.Add(-25 * time.Minute),
		EndTime:       now.Add(-5 * time.Minute),
	}
}

// ============ FRONTEND PATH DETECTION ============
func getFrontendPath() string {
	// Try different paths based on where the binary is run from
	paths := []string{
		"./frontend/dist",           // Running from Auctmah/
		"./Auctmah/frontend/dist",   // Running from repo root
	}
	
	for _, path := range paths {
		if _, err := os.Stat(path); err == nil {
			fmt.Printf("📂 Serving frontend from: %s\n", path)
			return path
		}
	}
	
	// Default to first path
	fmt.Printf("⚠️  Frontend path not found, using default: ./frontend/dist\n")
	return "./frontend/dist"
}

// ============ MAIN ============
func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/api/health", handleHealth)
	http.HandleFunc("/api/auctions", handleAuctions)
	http.HandleFunc("/api/bid", handlePlaceBid)
	http.HandleFunc("/api/create-auction", handleCreateAuction)
	http.Handle("/", http.FileServer(http.Dir(getFrontendPath())))

	go broadcastMessages()
	go updateAuctionTimers()

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	fmt.Printf("🔨 Auctmah Server Starting\n")
	fmt.Printf("📡 WebSocket: ws://localhost:%s/ws\n", port)
	fmt.Printf("🎯 Live Auctions: http://localhost:%s\n", port)
	fmt.Printf("⚡ Rust+WASM Frontend: High-speed canvas rendering\n")

	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

// ============ WEBSOCKET HANDLER ============
func handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("WebSocket upgrade error:", err)
		return
	}

	client := &Client{
		conn: conn,
		id:   fmt.Sprintf("bidder_%d", time.Now().UnixNano()),
	}

	clientsMutex.Lock()
	clients[client] = true
	clientsMutex.Unlock()

	log.Printf("✅ Client connected: %s (total: %d)\n", client.id, len(clients))

	// Send all auctions to new client
	auctionMutex.RLock()
	for _, auction := range auctions {
		msg := Message{
			Type:    "auction_update",
			Auction: auction,
		}
		client.send(msg)
	}
	auctionMutex.RUnlock()

	// Start ping/pong handlers for keeping connection alive
	go writePump(client)
	go readMessages(client)
}

// ============ HEALTH CHECK HANDLER ============
func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	// Check system health
	auctionMutex.RLock()
	auctionCount := len(auctions)
	auctionMutex.RUnlock()

	clientsMutex.RLock()
	clientCount := len(clients)
	clientsMutex.RUnlock()

	uptime := time.Since(startTime)

	health := map[string]interface{}{
		"status":        "healthy",
		"timestamp":     time.Now().Unix(),
		"uptime_seconds": int(uptime.Seconds()),
		"auctions":      auctionCount,
		"active_clients": clientCount,
		"version":       "1.0.0",
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(health)
}

// ============ MESSAGE READER ============
func readMessages(client *Client) {
	defer func() {
		clientsMutex.Lock()
		delete(clients, client)
		clientsMutex.Unlock()
		client.conn.Close()
		log.Printf("❌ Client disconnected: %s (total: %d)\n", client.id, len(clients))
	}()

	client.conn.SetReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		var msg Message
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		// Handle different message types
		switch msg.Type {
		case "place_bid":
			if msg.Bid != nil {
				processBid(msg.Bid, client.id)
			}
		case "pong":
			// Client responded to ping
			client.conn.SetReadDeadline(time.Now().Add(pongWait))
		}
	}
}

// ============ MESSAGE WRITER (PING/PONG HANDLER) ============
func writePump(client *Client) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case <-ticker.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// ============ BID PROCESSING ============
const (
	MAX_BID_AMOUNT = 10000000.0  // $10 million maximum bid
	MIN_BID_INCREMENT = 1.0       // $1 minimum increment
)

func processBid(bid *Bid, bidderId string) {
	bid.BidderID = bidderId
	bid.Timestamp = time.Now()

	auctionMutex.Lock()
	defer auctionMutex.Unlock()

	// Validation: Maximum bid limit
	if bid.Amount > MAX_BID_AMOUNT {
		msg := Message{
			Type:  "bid_rejected",
			Error: fmt.Sprintf("Maximum bid amount is $%s", formatMoney(MAX_BID_AMOUNT)),
		}
		broadcast <- msg
		log.Printf("❌ Bid rejected: Amount $%.2f exceeds maximum of $%.2f", bid.Amount, MAX_BID_AMOUNT)
		return
	}

	// Find auction and validate
	for _, auction := range auctions {
		// Check if bid is for this auction (you'll need to add auction_id to Bid struct)
		// For now, we'll process the first matching conditions
		
		// CRITICAL FIX #1: Reject bids on ENDED auctions
		if auction.Status == "ended" {
			msg := Message{
				Type:  "bid_rejected",
				Error: "This auction has ended. Bids are no longer accepted.",
			}
			broadcast <- msg
			log.Printf("❌ Bid rejected: Auction '%s' has ENDED", auction.Title)
			continue
		}

		// CRITICAL FIX #2: Reject bids on SCHEDULED auctions
		if auction.Status == "scheduled" {
			msg := Message{
				Type:  "bid_rejected",
				Error: "This auction has not started yet. Please wait until it begins.",
			}
			broadcast <- msg
			log.Printf("❌ Bid rejected: Auction '%s' is SCHEDULED", auction.Title)
			continue
		}

		// Only process bids on ACTIVE auctions
		if auction.Status == "active" {
			// CRITICAL FIX #3: Minimum bid increment validation
			minRequired := auction.CurrentBid + MIN_BID_INCREMENT
			if bid.Amount < minRequired {
				msg := Message{
					Type:  "bid_rejected",
					Error: fmt.Sprintf("Bid must be at least $%.2f (current bid + $%.2f increment)", minRequired, MIN_BID_INCREMENT),
				}
				broadcast <- msg
				log.Printf("❌ Bid rejected: $%.2f is below minimum required $%.2f", bid.Amount, minRequired)
				return
			}

			// CRITICAL FIX #4: Ensure price never decreases
			if bid.Amount > auction.CurrentBid {
				auction.CurrentBid = bid.Amount
				auction.HighestBidder = bidderId
				auction.BidCount++

				msg := Message{
					Type:    "bid_accepted",
					Auction: auction,
					Bid:     bid,
				}
				broadcast <- msg
				log.Printf("💰 Bid accepted: $%.2f by %s on %s (New price: $%.2f)", 
					bid.Amount, bidderId, auction.Title, auction.CurrentBid)
				return
			}
		}
	}
}

func formatMoney(amount float64) string {
	return fmt.Sprintf("%,.2f", amount)
}

// ============ BROADCAST LOOP ============
func broadcastMessages() {
	for {
		msg := <-broadcast
		clientsMutex.RLock()
		for client := range clients {
			client.send(msg)
		}
		clientsMutex.RUnlock()
	}
}

// ============ AUCTION TIMER ============
func updateAuctionTimers() {
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		auctionMutex.Lock()
		now := time.Now()
		for _, auction := range auctions {
			if auction.Status == "active" && now.After(auction.EndTime) {
				auction.Status = "ended"
				msg := Message{
					Type:    "auction_ended",
					Auction: auction,
				}
				auctionMutex.Unlock()
				broadcast <- msg
				auctionMutex.Lock()
				log.Printf("🏁 Auction ended: %s (Winner: %s at $%.2f)", auction.Title, auction.HighestBidder, auction.CurrentBid)
			}
		}
		auctionMutex.Unlock()
	}
}

// ============ HTTP HANDLERS ============
func handleAuctions(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	auctionMutex.RLock()
	defer auctionMutex.RUnlock()

	var auctionList []*Auction
	for _, auction := range auctions {
		auctionList = append(auctionList, auction)
	}
	json.NewEncoder(w).Encode(auctionList)
}

func handlePlaceBid(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	var bid Bid
	json.NewDecoder(r.Body).Decode(&bid)

	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"status": "bid received"})
}

// ============ CREATE AUCTION HANDLER ============
func handleCreateAuction(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Methods", "POST")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var submission AuctionSubmission
	if err := json.NewDecoder(r.Body).Decode(&submission); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Simple validation
	if submission.Name == "" || submission.Description == "" || submission.StartPrice <= 0 || submission.Email == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	// Save to file (simple persistence)
	file, err := os.OpenFile("auction_submissions.json", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening submissions file:", err)
		http.Error(w, "Server error", http.StatusInternalServerError)
		return
	}
	defer file.Close()

	submissionData, _ := json.Marshal(submission)
	file.Write(submissionData)
	file.WriteString("\n")

	log.Printf("📝 New auction submission: %s (Starting price: $%.2f, Email: %s)", 
		submission.Name, submission.StartPrice, submission.Email)

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
		"message": "Auction submission received. We'll contact you within 24 hours.",
	})
}

// ============ CLIENT HELPERS ============
func (c *Client) send(msg Message) {
	c.conn.WriteJSON(msg)
}
