#!/usr/bin/env powershell
# Test script to verify SocioVault API is working

$API_BASE = "http://localhost:8081/api"

Write-Host "🔍 Testing SocioVault API Endpoints" -ForegroundColor Cyan
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""

# Test 1: Check if server is running
Write-Host "1️⃣  Testing server connectivity..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$API_BASE/stats" -ErrorAction Stop
    Write-Host "✅ Server is running!" -ForegroundColor Green
    $stats = $response.Content | ConvertFrom-Json
    Write-Host "   Total Balance: `$$($stats.totalBalance)"
    Write-Host "   Stories Funded: $($stats.storiesFunded)"
    Write-Host ""
} catch {
    Write-Host "❌ Server is NOT responding!" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
    exit 1
}

# Test 2: Verify HTML loads
Write-Host "2️⃣  Testing HTML page load..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "http://localhost:8081" -ErrorAction Stop
    if ($response.Content -match "SocioVault") {
        Write-Host "✅ HTML page loads correctly!" -ForegroundColor Green
        Write-Host "   Title contains 'SocioVault': Yes" -ForegroundColor Green
    }
    Write-Host ""
} catch {
    Write-Host "❌ HTML page failed to load!" -ForegroundColor Red
    exit 1
}

# Test 3: Test contribution endpoint
Write-Host "3️⃣  Testing contribution submission..." -ForegroundColor Yellow
try {
    $body = @{
        email = "test@example.com"
        amount = 50.00
        message = "Test contribution from verification script"
    } | ConvertTo-Json
    
    $response = Invoke-WebRequest -Uri "$API_BASE/contribute" -Method POST -ContentType "application/json" -Body $body -ErrorAction Stop
    $contrib = $response.Content | ConvertFrom-Json
    Write-Host "✅ Contribution submitted successfully!" -ForegroundColor Green
    Write-Host "   Email: $($contrib.email)" -ForegroundColor Green
    Write-Host "   Amount: `$$($contrib.amount)" -ForegroundColor Green
    Write-Host "   ID: $($contrib.id)" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Contribution submission failed!" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
}

# Test 4: Get all contributions
Write-Host "4️⃣  Testing get contributions..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$API_BASE/contribute" -Method GET -ErrorAction Stop
    $contribs = $response.Content | ConvertFrom-Json
    Write-Host "✅ Retrieved contributions!" -ForegroundColor Green
    Write-Host "   Total: $($contribs.Count) contributor(s)" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Get contributions failed!" -ForegroundColor Red
}

# Test 5: Test request endpoint
Write-Host "5️⃣  Testing help request submission..." -ForegroundColor Yellow
try {
    $body = @{
        name = "Test User"
        email = "testuser@example.com"
        story = "I need help with emergency expenses"
        videoUrl = "https://example.com/video.mp4"
        amount = 200.00
    } | ConvertTo-Json
    
    $response = Invoke-WebRequest -Uri "$API_BASE/requests" -Method POST -ContentType "application/json" -Body $body -ErrorAction Stop
    $req = $response.Content | ConvertFrom-Json
    Write-Host "✅ Help request submitted successfully!" -ForegroundColor Green
    Write-Host "   Name: $($req.name)" -ForegroundColor Green
    Write-Host "   Amount: `$$($req.amount)" -ForegroundColor Green
    Write-Host "   ID: $($req.id)" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Help request submission failed!" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
}

# Test 6: Get all requests
Write-Host "6️⃣  Testing get help requests..." -ForegroundColor Yellow
try {
    $response = Invoke-WebRequest -Uri "$API_BASE/requests" -Method GET -ErrorAction Stop
    $reqs = $response.Content | ConvertFrom-Json
    Write-Host "✅ Retrieved help requests!" -ForegroundColor Green
    Write-Host "   Total: $($reqs.Count) request(s)" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Get help requests failed!" -ForegroundColor Red
}

# Test 7: Test subscription
Write-Host "7️⃣  Testing email subscription..." -ForegroundColor Yellow
try {
    $body = @{
        email = "subscriber@example.com"
    } | ConvertTo-Json
    
    $response = Invoke-WebRequest -Uri "$API_BASE/subscribe" -Method POST -ContentType "application/json" -Body $body -ErrorAction Stop
    $sub = $response.Content | ConvertFrom-Json
    Write-Host "✅ Subscription successful!" -ForegroundColor Green
    Write-Host "   Status: $($sub.status)" -ForegroundColor Green
    Write-Host "   Email: $($sub.email)" -ForegroundColor Green
    Write-Host ""
} catch {
    Write-Host "❌ Subscription failed!" -ForegroundColor Red
    Write-Host "   Error: $_" -ForegroundColor Red
}

# Test 8: Check data persistence
Write-Host "8️⃣  Checking data persistence..." -ForegroundColor Yellow
$dataDir = "c:\Users\aship\Desktop\Digital Orael\signal-bank-landing\data"
if (Test-Path "$dataDir\contributors.json") {
    $contribCount = (Get-Content "$dataDir\contributors.json" | ConvertFrom-Json).Count
    Write-Host "✅ Contributors data file exists!" -ForegroundColor Green
    Write-Host "   File: contributors.json" -ForegroundColor Green
    Write-Host "   Size: $([Math]::Round((Get-Item "$dataDir\contributors.json").Length / 1KB, 2)) KB" -ForegroundColor Green
}

if (Test-Path "$dataDir\requests.json") {
    $reqCount = (Get-Content "$dataDir\requests.json" | ConvertFrom-Json).Count
    Write-Host "✅ Requests data file exists!" -ForegroundColor Green
    Write-Host "   File: requests.json" -ForegroundColor Green
    Write-Host "   Size: $([Math]::Round((Get-Item "$dataDir\requests.json").Length / 1KB, 2)) KB" -ForegroundColor Green
}

if (Test-Path "$dataDir\subscribers.json") {
    Write-Host "✅ Subscribers data file exists!" -ForegroundColor Green
    Write-Host "   File: subscribers.json" -ForegroundColor Green
}

Write-Host ""
Write-Host "====================================" -ForegroundColor Cyan
Write-Host "✅ ALL TESTS COMPLETED!" -ForegroundColor Green
Write-Host "====================================" -ForegroundColor Cyan
Write-Host ""
Write-Host "🌐 Access the page at: http://localhost:8081" -ForegroundColor Cyan
Write-Host "📊 API Base: http://localhost:8081/api" -ForegroundColor Cyan
