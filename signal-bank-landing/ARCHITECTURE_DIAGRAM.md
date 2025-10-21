# SocioVault System - Visual Integration Map

## 🗺️ Complete System Architecture

```
╔════════════════════════════════════════════════════════════════════════════════╗
║                          SOCIAVAULT LANDING PAGE                              ║
║                        http://localhost:8081                                  ║
╚════════════════════════════════════════════════════════════════════════════════╝

┌──────────────────────────────────────────────────────────────────────────────┐
│                         FRONTEND LAYER (Browser)                            │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  HTML Elements (index.html)                                        │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │                                                                    │     │
│  │  <h3 data-stat="balance">$48,320</h3>                            │     │
│  │  <h3 data-stat="distributed">92%</h3>                           │     │
│  │  <h3 data-stat="stories">311</h3>                               │     │
│  │  │                                                               │     │
│  │  └─► Targets: updateStatsDisplay() in script.js                │     │
│  │                                                                    │     │
│  │  <form id="contribute-form" class="form-grid">                  │     │
│  │    <input name="email" ... required>                             │     │
│  │    <input name="amount" ... required>                            │     │
│  │    <textarea name="message">                                      │     │
│  │    <button>Contribute now</button>                               │     │
│  │  </form>                                                          │     │
│  │  │                                                               │     │
│  │  └─► Targets: handleContribution() in script.js                │     │
│  │                                                                    │     │
│  │  <form id="request-form" class="form-grid">                     │     │
│  │    <input name="name" ... required>                              │     │
│  │    <input name="email" ... required>                             │     │
│  │    <textarea name="story" ... required>                          │     │
│  │    <input name="videoUrl">                                        │     │
│  │    <input name="amount" ... required>                            │     │
│  │    <button>Submit request</button>                               │     │
│  │  </form>                                                          │     │
│  │  │                                                               │     │
│  │  └─► Targets: handleRequestSubmit() in script.js               │     │
│  │                                                                    │     │
│  │  <form class="signup">                                            │     │
│  │    <input type="email" name="email" ... required>                │     │
│  │    <button>Subscribe</button>                                     │     │
│  │  </form>                                                          │     │
│  │  │                                                               │     │
│  │  └─► Targets: handleSubscribe() in script.js                   │     │
│  │                                                                    │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  CSS Styling (styles.css)                                         │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │  ✅ All forms styled                                             │     │
│  │  ✅ Buttons with hover effects                                    │     │
│  │  ✅ Responsive layout                                             │     │
│  │  ✅ Animation effects                                             │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  JavaScript Logic (script.js)                                      │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │                                                                    │     │
│  │  ┌─ API_BASE Detection ────────────────────────────────────┐      │     │
│  │  │  localhost → http://localhost:8081/api              │      │     │
│  │  │  production → ${origin}/api                         │      │     │
│  │  └─────────────────────────────────────────────────────┘      │     │
│  │                                                                    │     │
│  │  ┌─ Form Event Listeners ──────────────────────────────────┐      │     │
│  │  │  DOMContentLoaded                                      │      │     │
│  │  │  ├─ loadStats()                                        │      │     │
│  │  │  └─ setupFormHandlers()                                │      │     │
│  │  │                                                         │      │     │
│  │  │  setupFormHandlers()                                   │      │     │
│  │  │  ├─ contribute-form → handleContribution              │      │     │
│  │  │  ├─ request-form → handleRequestSubmit                │      │     │
│  │  │  └─ .signup → handleSubscribe                         │      │     │
│  │  └─────────────────────────────────────────────────────┘      │     │
│  │                                                                    │     │
│  │  ┌─ Form Submission Flow ──────────────────────────────────┐      │     │
│  │  │                                                         │      │     │
│  │  │  handleContribution()                                 │      │     │
│  │  │  ├─ Get form data                                     │      │     │
│  │  │  ├─ Validate (amount > 0)                            │      │     │
│  │  │  ├─ POST /api/contribute                             │      │     │
│  │  │  ├─ Alert success/error                              │      │     │
│  │  │  ├─ Reset form                                        │      │     │
│  │  │  └─ loadStats()                                       │      │     │
│  │  │                                                         │      │     │
│  │  │  handleRequestSubmit()                                │      │     │
│  │  │  ├─ Get form data                                     │      │     │
│  │  │  ├─ Validate (required fields)                       │      │     │
│  │  │  ├─ POST /api/requests                               │      │     │
│  │  │  ├─ Alert success/error                              │      │     │
│  │  │  ├─ Reset form                                        │      │     │
│  │  │  └─ loadStats()                                       │      │     │
│  │  │                                                         │      │     │
│  │  │  handleSubscribe()                                    │      │     │
│  │  │  ├─ Get email                                         │      │     │
│  │  │  ├─ POST /api/subscribe                              │      │     │
│  │  │  ├─ Alert success/error                              │      │     │
│  │  │  └─ Reset form                                        │      │     │
│  │  └─────────────────────────────────────────────────────┘      │     │
│  │                                                                    │     │
│  │  ┌─ Stats Management ──────────────────────────────────────┐      │     │
│  │  │                                                         │      │     │
│  │  │  loadStats()                                           │      │     │
│  │  │  ├─ GET /api/stats                                    │      │     │
│  │  │  ├─ Parse JSON response                               │      │     │
│  │  │  └─ updateStatsDisplay(stats)                        │      │     │
│  │  │                                                         │      │     │
│  │  │  updateStatsDisplay(stats)                            │      │     │
│  │  │  ├─ Find [data-stat="balance"]                       │      │     │
│  │  │  ├─ Update text: $48,320 + new                        │      │     │
│  │  │  ├─ Find [data-stat="distributed"]                   │      │     │
│  │  │  ├─ Update text: 92%                                  │      │     │
│  │  │  ├─ Find [data-stat="stories"]                       │      │     │
│  │  │  └─ Update text: 311 + new                            │      │     │
│  │  │                                                         │      │     │
│  │  │  setInterval(loadStats, 30000)                        │      │     │
│  │  │  └─ Auto-refresh every 30 seconds                    │      │     │
│  │  └─────────────────────────────────────────────────────┘      │     │
│  │                                                                    │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘


                          HTTP REQUESTS/RESPONSES
                    ═══════════════════════════════════

┌──────────────────────────────────────────────────────────────────────────────┐
│                       BACKEND LAYER (Go Server)                             │
│                      sociovault.exe on port 8081                            │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  API Endpoints (main.go)                                           │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │                                                                    │     │
│  │  GET /  (Static Files)                                           │     │
│  │  ├─ Serves: index.html, styles.css, script.js                   │     │
│  │  └─ Sets: Content-Type: text/html, text/css, text/javascript    │     │
│  │                                                                    │     │
│  │  POST /api/contribute  (handleContribute)                        │     │
│  │  ├─ Receives: JSON {email, amount, message}                    │     │
│  │  ├─ Validation: JSON unmarshal success                          │     │
│  │  ├─ Processing: Create struct, ID, timestamp                    │     │
│  │  ├─ Storage: Append to memory, write to file                   │     │
│  │  └─ Response: 201 Created + JSON response                       │     │
│  │                                                                    │     │
│  │  GET /api/contribute  (handleContribute)                        │     │
│  │  ├─ Processing: Load from file                                 │     │
│  │  ├─ Read: contributors array                                     │     │
│  │  └─ Response: 200 OK + JSON array                               │     │
│  │                                                                    │     │
│  │  POST /api/requests  (handleRequest)                            │     │
│  │  ├─ Receives: JSON {name, email, story, videoUrl, amount}     │     │
│  │  ├─ Validation: JSON unmarshal success                          │     │
│  │  ├─ Processing: Create struct, ID, timestamp                    │     │
│  │  ├─ Storage: Append to memory, write to file                   │     │
│  │  └─ Response: 201 Created + JSON response                       │     │
│  │                                                                    │     │
│  │  GET /api/requests  (handleRequest)                             │     │
│  │  ├─ Processing: Load from file                                 │     │
│  │  ├─ Read: requests array                                         │     │
│  │  └─ Response: 200 OK + JSON array                               │     │
│  │                                                                    │     │
│  │  GET /api/stats  (handleStats)                                  │     │
│  │  ├─ Processing:                                                 │     │
│  │  │  ├─ Read contributors.json → sum amounts                    │     │
│  │  │  ├─ Read requests.json → count entries                      │     │
│  │  │  └─ Calculate: totals + base values                         │     │
│  │  ├─ Returns:                                                     │     │
│  │  │  ├─ totalBalance: 48320 + sum                               │     │
│  │  │  ├─ distributedPercent: 92                                  │     │
│  │  │  ├─ storiesFunded: 311 + requests count                    │     │
│  │  │  ├─ totalContributors: 2847 + contributors count           │     │
│  │  │  ├─ activeRequests: count                                   │     │
│  │  │  └─ dailyContributions: today's sum                        │     │
│  │  └─ Response: 200 OK + JSON stats                              │     │
│  │                                                                    │     │
│  │  POST /api/vote  (handleVote)                                   │     │
│  │  ├─ Receives: JSON {requestId}                                │     │
│  │  ├─ Processing: Find request, increment votes                  │     │
│  │  ├─ Storage: Write to file                                     │     │
│  │  └─ Response: 200 OK + updated request JSON                    │     │
│  │                                                                    │     │
│  │  POST /api/subscribe  (handleSubscribe)                         │     │
│  │  ├─ Receives: JSON {email}                                     │     │
│  │  ├─ Processing: Create subscription record                      │     │
│  │  ├─ Storage: Append to subscribers.json                        │     │
│  │  └─ Response: 201 Created + JSON response                       │     │
│  │                                                                    │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  CORS Configuration (All Endpoints)                              │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │  Access-Control-Allow-Origin: *                                 │     │
│  │  Access-Control-Allow-Methods: GET, POST, OPTIONS               │     │
│  │  Access-Control-Allow-Headers: Content-Type                    │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  Data Structures (Go Structs)                                     │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │                                                                    │     │
│  │  type Contributor struct {                                       │     │
│  │    ID        string    `json:"id"`                              │     │
│  │    Email     string    `json:"email"`                           │     │
│  │    Amount    float64   `json:"amount"`                          │     │
│  │    Message   string    `json:"message"`                         │     │
│  │    CreatedAt time.Time `json:"createdAt"`                       │     │
│  │  }                                                               │     │
│  │                                                                    │     │
│  │  type Request struct {                                           │     │
│  │    ID        string    `json:"id"`                              │     │
│  │    Name      string    `json:"name"`                            │     │
│  │    Email     string    `json:"email"`                           │     │
│  │    Story     string    `json:"story"`                           │     │
│  │    VideoURL  string    `json:"videoUrl"`                        │     │
│  │    Amount    float64   `json:"amount"`                          │     │
│  │    Verified  bool      `json:"verified"`                        │     │
│  │    Votes     int       `json:"votes"`                           │     │
│  │    CreatedAt time.Time `json:"createdAt"`                       │     │
│  │  }                                                               │     │
│  │                                                                    │     │
│  │  type Stats struct {                                             │     │
│  │    TotalBalance        float64 `json:"totalBalance"`            │     │
│  │    DistributedPercent  int     `json:"distributedPercent"`     │     │
│  │    StoriesFunded       int     `json:"storiesFunded"`           │     │
│  │    TotalContributors   int     `json:"totalContributors"`       │     │
│  │    ActiveRequests      int     `json:"activeRequests"`          │     │
│  │    DailyContributions  float64 `json:"dailyContributions"`      │     │
│  │  }                                                               │     │
│  │                                                                    │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘


                            DATA STORAGE LAYER
                       ════════════════════════════════

┌──────────────────────────────────────────────────────────────────────────────┐
│              Data Directory: ./data/                                         │
├──────────────────────────────────────────────────────────────────────────────┤
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  contributors.json                                                │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │  [                                                                 │     │
│  │    {                                                               │     │
│  │      "id": "1729537200000000000",                                │     │
│  │      "email": "user@example.com",                               │     │
│  │      "amount": 50,                                              │     │
│  │      "message": "Support the community",                        │     │
│  │      "createdAt": "2025-10-21T14:00:00Z"                       │     │
│  │    }                                                               │     │
│  │  ]                                                                 │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  requests.json                                                    │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │  [                                                                 │     │
│  │    {                                                               │     │
│  │      "id": "1729537300000000000",                                │     │
│  │      "name": "John Doe",                                         │     │
│  │      "email": "john@example.com",                               │     │
│  │      "story": "Need help with medical bills",                   │     │
│  │      "videoUrl": "https://example.com/video.mp4",              │     │
│  │      "amount": 500,                                             │     │
│  │      "verified": false,                                         │     │
│  │      "votes": 0,                                                │     │
│  │      "createdAt": "2025-10-21T14:05:00Z"                       │     │
│  │    }                                                               │     │
│  │  ]                                                                 │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
│  ┌────────────────────────────────────────────────────────────────────┐     │
│  │  subscribers.json                                                 │     │
│  ├────────────────────────────────────────────────────────────────────┤     │
│  │  [                                                                 │     │
│  │    {                                                               │     │
│  │      "email": "subscriber@example.com",                         │     │
│  │      "timestamp": "2025-10-21T14:10:00Z"                       │     │
│  │    }                                                               │     │
│  │  ]                                                                 │     │
│  └────────────────────────────────────────────────────────────────────┘     │
│                                                                              │
└──────────────────────────────────────────────────────────────────────────────┘


═════════════════════════════════════════════════════════════════════════════════
                           CONNECTION SUMMARY
═════════════════════════════════════════════════════════════════════════════════

Frontend Form                 →  API Endpoint              →  Backend Handler
─────────────────────────────────────────────────────────────────────────────
contribute-form              →  POST /api/contribute      →  handleContribute()
request-form                 →  POST /api/requests        →  handleRequest()
.signup form                 →  POST /api/subscribe       →  handleSubscribe()
[data-stat] elements         →  GET /api/stats            →  handleStats()
(Vote functionality)         →  POST /api/vote            →  handleVote()

Data Flow:
──────────
User Input → Form Validation → JSON Payload → HTTP Request → 
Server Processing → JSON Response → Browser Alert → DOM Update → Form Reset

All connections verified ✅
All endpoints functional ✅
All data models matched ✅
Full integration complete ✅
