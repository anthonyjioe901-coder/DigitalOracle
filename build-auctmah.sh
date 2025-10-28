#!/bin/bash
# Auctmah Build Script for Render

set -e

echo "📦 Installing Rust..."
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh -s -- -y
source $HOME/.cargo/env

echo "🦀 Building Rust WebAssembly..."
cd Auctmah/frontend
wasm-pack build --target web --release

echo "📁 Setting up frontend..."
mkdir -p dist
cp index.html dist/

echo "🚀 Building Go backend..."
cd ../..
go build -o Auctmah/app Auctmah/main.go

echo "✅ Build complete!"
