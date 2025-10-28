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

// ============ GLOBAL STATE ============
var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}

	auctions      = make(map[string]*Auction)
	auctionMutex  sync.RWMutex

	clients       = make(map[*Client]bool)
	clientsMutex  sync.RWMutex

	broadcast     = make(chan Message, 256)
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

// ============ MAIN ============
func main() {
	http.HandleFunc("/ws", handleWebSocket)
	http.HandleFunc("/api/auctions", handleAuctions)
	http.HandleFunc("/api/bid", handlePlaceBid)
	http.Handle("/", http.FileServer(http.Dir("./frontend")))

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

	go readMessages(client)
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

	client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
	client.conn.SetPongHandler(func(string) error {
		client.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		return nil
	})

	for {
		var msg Message
		err := client.conn.ReadJSON(&msg)
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				log.Printf("WebSocket error: %v", err)
			}
			break
		}

		if msg.Type == "place_bid" && msg.Bid != nil {
			processBid(msg.Bid, client.id)
		}
	}
}

// ============ BID PROCESSING ============
func processBid(bid *Bid, bidderId string) {
	bid.BidderID = bidderId
	bid.Timestamp = time.Now()

	auctionMutex.Lock()
	defer auctionMutex.Unlock()

	// Find active auction and update
	for _, auction := range auctions {
		if auction.Status == "active" && bid.Amount > auction.CurrentBid {
			auction.CurrentBid = bid.Amount
			auction.HighestBidder = bidderId
			auction.BidCount++

			msg := Message{
				Type:    "bid_accepted",
				Auction: auction,
				Bid:     bid,
			}
			broadcast <- msg
			log.Printf("💰 Bid accepted: $%.2f by %s on %s", bid.Amount, bidderId, auction.Title)
			break
		}
	}
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

// ============ CLIENT HELPERS ============
func (c *Client) send(msg Message) {
	c.conn.WriteJSON(msg)
}
