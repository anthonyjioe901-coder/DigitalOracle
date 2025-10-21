# 🧪 Complete Testing Guide

## Quick Start Testing

### 1. Verify Server is Running
```
Port: 8081
URL: http://localhost:8081
```

Open browser → should see the SocioVault landing page with all stats visible.

### 2. Test Stats Display

**What to check:**
- Community balance: $48,320+ (updates when contributions added)
- Distributed %: 92%
- Stories funded: 311+

**How to verify it's linked:**
```javascript
// Open browser dev tools (F12)
// Go to Console tab
// You should see:
// "API Base URL: http://localhost:8081/api"
```

### 3. Test Contribution Form

**Form Fields:**
- Email: your@email.com
- Amount: 50.00
- Message: "Test contribution"

**Expected Result:**
1. ✅ Alert: "Thank you for your contribution! 🎉"
2. ✅ Form clears
3. ✅ Stats update in real-time
4. ✅ Data saved to file

**Verify Data Saved:**
```
c:\Users\aship\Desktop\Digital Orael\signal-bank-landing\data\contributors.json
```
Should contain your entry with: email, amount, message, id, createdAt

### 4. Test Help Request Form

**Form Fields:**
- Name: John Doe
- Email: john@example.com
- Story: "I need emergency help with medical bills"
- Video URL: (optional) https://example.com/video.mp4
- Amount: 500.00

**Expected Result:**
1. ✅ Alert: "Your request has been submitted!"
2. ✅ Form clears
3. ✅ Stats update (stories count increases)
4. ✅ Data saved to file

**Verify Data Saved:**
```
c:\Users\aship\Desktop\Digital Orael\signal-bank-landing\data\requests.json
```

### 5. Test Email Subscription

**Form Field:**
- Email: subscriber@example.com

**Expected Result:**
1. ✅ Alert: "Successfully subscribed! ✨"
2. ✅ Form clears
3. ✅ Data saved

**Verify Data Saved:**
```
c:\Users\aship\Desktop\Digital Orael\signal-bank-landing\data\subscribers.json
```

### 6. Test Stats Auto-Refresh

1. Submit a contribution
2. Wait 30 seconds (watch the stats)
3. Stats should auto-update WITHOUT page refresh
4. No manual F5 needed

**Why it works:**
- JavaScript calls `loadStats()` every 30 seconds
- Fetches latest data from `/api/stats`
- Updates all `data-stat` elements

## Browser Developer Tools Debugging

### Check Network Requests

1. Open DevTools (F12)
2. Click "Network" tab
3. Reload page
4. Look for requests to:
   - `stats` ✅ Should be 200 OK
   - `contribute` ✅ Should be 201 Created (after form submit)
   - `requests` ✅ Should be 201 Created (after form submit)
   - `subscribe` ✅ Should be 201 Created (after form submit)

### Check Console for Errors

1. Open DevTools (F12)
2. Click "Console" tab
3. Should see:
   - ✅ `API Base URL: http://localhost:8081/api`
   - ❌ NO red errors
   - ⚠️ Any warnings are OK

### Check Form Elements

```javascript
// Open Console (F12)
// Type and run:

document.getElementById("contribute-form")
// Should return: <form id="contribute-form" ...>

document.getElementById("request-form")
// Should return: <form id="request-form" ...>

document.querySelector(".signup")
// Should return: <form class="signup" ...>
```

### Check Data Attributes

```javascript
// In Console:
document.querySelectorAll("[data-stat]")
// Should return 3 elements:
// [0] h3[data-stat="balance"]
// [1] h3[data-stat="distributed"]
// [2] h3[data-stat="stories"]
```

## Command Line Testing

### Test with curl

```powershell
# Get stats
curl http://localhost:8081/api/stats

# Test contribution (Windows PowerShell)
$body = '{"email":"test@example.com","amount":50,"message":"test"}'
curl -X POST -H "Content-Type: application/json" -d $body http://localhost:8081/api/contribute

# Test request
$body = '{"name":"Test","email":"test@example.com","story":"help needed","amount":100}'
curl -X POST -H "Content-Type: application/json" -d $body http://localhost:8081/api/requests

# Test subscribe
$body = '{"email":"sub@example.com"}'
curl -X POST -H "Content-Type: application/json" -d $body http://localhost:8081/api/subscribe
```

### Check Data Files

```powershell
# View contributors
cat c:\Users\aship\Desktop\Digital Orael\signal-bank-landing\data\contributors.json

# View requests
cat c:\Users\aship\Desktop\Digital Orael\signal-bank-landing\data\requests.json

# View subscribers
cat c:\Users\aship\Desktop\Digital Orael\signal-bank-landing\data\subscribers.json
```

## Complete Integration Test Scenario

### Scenario: User Journey

1. **User visits page**
   - Browser opens: http://localhost:8081
   - ✅ Page loads with SocioVault branding
   - ✅ Navigation visible
   - ✅ Stats display (💰 $48,320, ✓ 92%, ❤️ 311)

2. **User makes a contribution**
   - Scrolls to "Make a contribution" section
   - Fills form: email@example.com, $100, "I believe in community"
   - Clicks "Contribute now"
   - ✅ Alert appears
   - ✅ Form clears
   - ✅ Stats update (+$100 to balance)

3. **User requests help**
   - Scrolls to "Ready to request support?"
   - Fills form: John Doe, john@email.com, "Need help with rent", $1500
   - Clicks "Submit request"
   - ✅ Alert appears
   - ✅ Stats update (stories +1)

4. **User subscribes for updates**
   - Enters: updates@email.com
   - Clicks "Subscribe"
   - ✅ Alert appears
   - ✅ Form clears

5. **Stats auto-refresh**
   - Wait 30 seconds
   - Stats update automatically
   - No page refresh needed

6. **Verify data persistence**
   - Check `data/contributors.json` → has 1 entry
   - Check `data/requests.json` → has 1 entry
   - Check `data/subscribers.json` → has 1 entry

## Troubleshooting

### Stats not updating?
1. Check browser console for errors
2. Check DevTools Network tab for failed requests
3. Verify API endpoint: http://localhost:8081/api/stats
4. Restart server if needed

### Forms not submitting?
1. Check console for JavaScript errors
2. Check Network tab for POST requests
3. Verify form has correct ID: `id="contribute-form"` or `id="request-form"`
4. Check that data/directory exists

### Page not loading?
1. Verify port 8081 is listening: `netstat -ano | findstr 8081`
2. Restart server: `cd signal-bank-landing && ./sociovault.exe`
3. Check firewall isn't blocking port 8081

### Data not saving?
1. Check `data/` directory exists
2. Check write permissions on directory
3. Check for error messages in Go server output
4. Try creating file manually to test permissions

## Success Indicators

When everything is working correctly, you should see:

✅ Page loads instantly  
✅ Stats display on load  
✅ Contribution form submits and alerts success  
✅ Help request form submits and alerts success  
✅ Email subscription submits and alerts success  
✅ JSON files created in data/ directory  
✅ Stats auto-refresh every 30 seconds  
✅ No console errors  
✅ All Network requests are 200/201  
✅ Forms clear after submission  

🎉 **Everything is working!**
