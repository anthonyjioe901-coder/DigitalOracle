# ⚡ Quick Reference - SocioVault System

## 🚀 Start Here

### Quick Start
```bash
cd c:\Users\aship\Desktop\Digital Orael\signal-bank-landing
./sociovault.exe
# Open: http://localhost:8081
```

### Server Status
- ✅ **Running on**: http://localhost:8081
- ✅ **Port**: 8081
- ✅ **Data Directory**: ./data/

---

## 📋 All Connections Summary

| Component | File | Status |
|-----------|------|--------|
| Frontend HTML | `index.html` | ✅ All forms with correct IDs |
| Frontend CSS | `styles.css` | ✅ Fully styled |
| Frontend JS | `script.js` | ✅ All handlers connected |
| Backend API | `main.go` | ✅ All endpoints ready |
| API Base URL | Automatic detection | ✅ Works locally & production |
| Data Storage | `./data/*.json` | ✅ Auto-created |
| CORS Config | Response headers | ✅ Enabled |

---

## 🔄 Data Flow Diagram

```
┌─────────────┐
│   Browser   │ http://localhost:8081
└──────┬──────┘
       │
       ├─ Form 1: Contribution (email, amount, message)
       │  └─ POST /api/contribute
       │     └─ Backend → Save to contributors.json
       │        └─ Response → Alert + Stats Update
       │
       ├─ Form 2: Help Request (name, email, story, amount)
       │  └─ POST /api/requests
       │     └─ Backend → Save to requests.json
       │        └─ Response → Alert + Stats Update
       │
       ├─ Form 3: Subscribe (email)
       │  └─ POST /api/subscribe
       │     └─ Backend → Save to subscribers.json
       │        └─ Response → Alert
       │
       └─ Auto-refresh (every 30 seconds)
          └─ GET /api/stats
             └─ Backend → Calculate totals
                └─ Response → Update stats display
```

---

## 📝 Form Configuration

### Form 1: Contribution
**Location**: `#contribute-form` section
**Endpoint**: `POST /api/contribute`
**Fields**:
- email (required)
- amount (required, >0)
- message (optional)

**Success**: Alert + Form Reset + Stats Update

---

### Form 2: Help Request
**Location**: `#request-form` section
**Endpoint**: `POST /api/requests`
**Fields**:
- name (required)
- email (required)
- story (required)
- videoUrl (optional)
- amount (required)

**Success**: Alert + Form Reset + Stats Update

---

### Form 3: Subscribe
**Location**: `.signup` form
**Endpoint**: `POST /api/subscribe`
**Fields**:
- email (required)

**Success**: Alert + Form Reset

---

## 🔌 API Endpoints Checklist

```
✅ GET /                    → Serves HTML/CSS/JS
✅ POST /api/contribute     → Accept contribution
✅ GET /api/contribute      → List contributions
✅ POST /api/requests       → Submit help request
✅ GET /api/requests        → List help requests
✅ GET /api/stats           → Real-time statistics
✅ POST /api/vote           → Vote on requests
✅ POST /api/subscribe      → Email subscription
```

---

## 📊 Stats Display Elements

| Element | Data Attribute | Backend Field | Updates From |
|---------|----------------|---------------|--------------|
| Balance | `data-stat="balance"` | `totalBalance` | `/api/stats` |
| Distribution % | `data-stat="distributed"` | `distributedPercent` | `/api/stats` |
| Stories Funded | `data-stat="stories"` | `storiesFunded` | `/api/stats` |

---

## 🧪 Quick Test Checklist

- [ ] Server running: http://localhost:8081 loads
- [ ] Stats display: Shows $48,320 + 92% + 311
- [ ] Contribution form: Submit and see alert
- [ ] Data saved: Check `data/contributors.json`
- [ ] Request form: Submit and see alert
- [ ] Data saved: Check `data/requests.json`
- [ ] Subscribe: Submit and see alert
- [ ] Data saved: Check `data/subscribers.json`
- [ ] Auto-refresh: Wait 30s, stats update without refresh
- [ ] No errors: Check browser console (F12)

---

## 💾 Data Files

| File | Purpose | Auto-Created |
|------|---------|--------------|
| `data/contributors.json` | Contributions | ✅ Yes |
| `data/requests.json` | Help requests | ✅ Yes |
| `data/subscribers.json` | Email list | ✅ Yes |

---

## 🛠️ Troubleshooting

### Issue: Page doesn't load
```
Fix: Check port 8081 is listening
netstat -ano | findstr 8081
Restart: ./sociovault.exe
```

### Issue: Stats not updating
```
Fix: Check browser console (F12)
Look for: "API Base URL: http://localhost:8081/api"
Check Network tab for /api/stats requests
```

### Issue: Form won't submit
```
Fix: Check DevTools Console for errors
Verify form has correct ID (contribute-form, request-form)
Check Network tab for POST request status
```

### Issue: Data not saving
```
Fix: Check data/ directory exists
Check write permissions
Look for errors in Go server output
Try restarting server
```

---

## 📱 Browser Developer Tools Quick Tips

**Check API Base URL:**
```javascript
// Console → Type:
API_BASE
// Should show: http://localhost:8081/api
```

**Check Form Elements:**
```javascript
// Console → Type:
document.getElementById("contribute-form")
document.getElementById("request-form")
document.querySelector(".signup")
// All should return the form element
```

**Check Stats Elements:**
```javascript
// Console → Type:
document.querySelectorAll("[data-stat]")
// Should return 3 h3 elements with data-stat attributes
```

**Trigger Stats Update:**
```javascript
// Console → Type:
loadStats()
// Should fetch fresh stats from API
```

---

## 🔗 Quick Links

| What | Where |
|------|-------|
| Landing Page | http://localhost:8081 |
| API Base | http://localhost:8081/api |
| Stats Endpoint | http://localhost:8081/api/stats |
| Contribution Data | ./data/contributors.json |
| Request Data | ./data/requests.json |
| Subscriber Data | ./data/subscribers.json |

---

## 📚 Documentation Files

| File | Purpose |
|------|---------|
| `README.md` | General overview |
| `IMPLEMENTATION.md` | Setup instructions |
| `LINKAGE_VERIFICATION.md` | All connections verified |
| `SYSTEM_LINKAGE_REPORT.md` | Complete integration report |
| `TESTING_GUIDE.md` | Step-by-step testing |
| `ARCHITECTURE_DIAGRAM.md` | Visual system map |
| `QUICK_REFERENCE.md` | This file |

---

## ✅ System Status

```
Frontend Layer:      ✅ Working
Backend Layer:       ✅ Working
Data Storage:        ✅ Working
API Endpoints:       ✅ All 8 functional
Form Handlers:       ✅ All 3 connected
Stats Display:       ✅ Live updating
CORS Config:         ✅ Enabled
Error Handling:      ✅ Implemented
User Feedback:       ✅ Alerts working
Auto-refresh:        ✅ Every 30 seconds
```

**Status: ✅ FULLY OPERATIONAL**

---

**Last Updated**: October 21, 2025  
**Version**: 1.0  
**All Systems**: GO! 🚀
