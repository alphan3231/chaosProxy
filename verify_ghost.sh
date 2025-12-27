#!/bin/bash

echo "👻 Testing GHOST MODE..."

# 1. Stop the target service (We can't stop httpbin.org, so we point proxy to bad port)
echo "🛑 Simulating Backend Failure (Restarting Proxy with BAD target)..."
# Kill existing proxy
pkill -f "cmd/sentinel/main.go" || true
# Ensure port 8081 is free
lsof -ti:8081 | xargs kill -9 2>/dev/null || true
echo "Waiting for port cleanup..."
sleep 2

# Start Proxy pointing to localhost:9999 (which is closed)
export TARGET_URL="http://localhost:9999"
export PORT="8081" # Use different port to avoid conflicts if old one stuck
nohup go run cmd/sentinel/main.go > sentinel_ghost.log 2>&1 &
PROXY_PID=$!

echo "⏳ Waiting for Proxy to start on :8081..."
sleep 3

# 2. Send Request
echo "⚡ Sending request to /get (should be served from Cache)..."
RESPONSE=$(curl -s -i "http://localhost:8081/get?param=test")

# 3. Analyze
echo "📝 Response Headers:"
echo "$RESPONSE" | grep "X-Chaos"

if echo "$RESPONSE" | grep -q "X-Chaos-Ghost: true"; then
    echo "✅ SUCCESS: Ghost Mode Activated! Served from Redis."
else
    echo "❌ FAILURE: Ghost Mode NOT detected."
    cat sentinel_ghost.log
fi

# Cleanup
kill $PROXY_PID
echo "🧹 Cleaned up."
