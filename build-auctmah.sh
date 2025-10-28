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
cp -r pkg/* dist/ 2>/dev/null || true

echo "🚀 Building Go backend..."
cd ../..
go build -o Auctmah/app Auctmah/main.go

echo "✅ Build complete!"
