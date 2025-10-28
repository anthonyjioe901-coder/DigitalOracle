# 🔨 AUCTMAH - Render Deployment Guide

## Quick Setup on Render

### Step 1: Create New Service
1. Go to https://render.com/dashboard
2. Click **"New"** → **"Web Service"**
3. Connect your GitHub repo: `anthonyjioe901-coder/DigitalOracle`
4. Select branch: `main`

### Step 2: Configure Service

Fill in these settings:

```
Name:                auctmah
Environment:         Go
Runtime Version:     1.22
Build Command:       ./build-auctmah.sh
Start Command:       ./Auctmah/app
Port:                8080
```

### Step 3: Environment Variables

Add these environment variables:

```
PORT=8080
```

### Step 4: Deploy

1. Click **"Deploy"**
2. Watch the build logs
3. Once deployed, you'll get a URL like: `https://auctmah.onrender.com`

---

## Troubleshooting

### Build fails: "wasm-pack not found"
**Solution:** The build script installs Rust automatically. Just wait for the build to complete.

### Build timeout (>30 minutes)
**Solution:** 
- Render has a 30-minute timeout
- Try upgrading to a paid plan for longer builds
- Or pre-build WASM locally and commit `pkg/` folder

### Port issues
Make sure `PORT=8080` is set in environment variables

---

## What Gets Deployed

```
Auctmah/
├─ main.go ✅ (compiled to binary)
├─ go.mod ✅
├─ frontend/
│  ├─ src/lib.rs ✅ (compiled to WASM)
│  ├─ Cargo.toml ✅
│  ├─ dist/
│  │  ├─ index.html ✅
│  │  └─ pkg/ ✅ (WASM output)
│  └─ index.html ✅
```

---

## Live URL

Once deployed, your live auction system will be at:
**https://auctmah.onrender.com** 🚀

---

## Manual Local Build (for testing)

```bash
cd Auctmah/frontend
wasm-pack build --target web --release
mkdir -p dist
cp index.html dist/

cd ../..
go build -o auctmah.exe Auctmah/main.go
./auctmah.exe
```

Then open: http://localhost:8080

---

**Happy Auctioning! 🔨💰**
