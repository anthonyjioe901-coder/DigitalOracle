#!/bin/bash
# Auctmah Build Script for Render

set -e

echo "📦 Installing Rust..."
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y --quiet

export PATH="$HOME/.cargo/bin:$PATH"

echo "⏳ Installing wasm-pack..."
cargo install wasm-pack

echo "🦀 Building Rust WebAssembly..."
cd Auctmah/frontend
wasm-pack build --target web --release

echo "📁 Setting up frontend distribution..."
mkdir -p dist
cp index.html dist/
cp pkg/*.js dist/ 2>/dev/null || true
cp pkg/*.wasm dist/ 2>/dev/null || true
ls -la dist/

echo "📦 Installing Go dependencies..."
cd ..
go mod download
go mod tidy

echo " Building Go backend..."
go build -o app main.go

cd ..
echo "✅ Build complete!"
