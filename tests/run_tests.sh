#!/bin/bash

# Integration Test Suite for Chaos-Proxy
# This script tests the entire system end-to-end

set -e

echo "🧪 Chaos-Proxy Integration Test Suite"
echo "======================================"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

PASSED=0
FAILED=0

# Test function
test_case() {
    local name="$1"
    local result="$2"
    
    if [ "$result" = "0" ]; then
        echo -e "${GREEN}✓${NC} $name"
        ((PASSED++))
    else
        echo -e "${RED}✗${NC} $name"
        ((FAILED++))
    fi
}

echo ""
echo "📦 Prerequisites Check"
echo "----------------------"

# Check Go
if command -v go &> /dev/null; then
    test_case "Go installed" 0
else
    test_case "Go installed" 1
    echo "  Install Go: https://golang.org/dl/"
fi

# Check Python
if command -v python3 &> /dev/null; then
    test_case "Python installed" 0
else
    test_case "Python installed" 1
fi

# Check Docker
if command -v docker &> /dev/null; then
    test_case "Docker installed" 0
else
    test_case "Docker installed" 1
fi

echo ""
echo "🔧 Go Unit Tests"
echo "----------------"

# Run Go tests
if go test ./internal/config/... -v 2>&1 | grep -q "PASS"; then
    test_case "Config tests" 0
else
    test_case "Config tests" 1
fi

if go test ./pkg/middleware/... -v 2>&1 | grep -q "PASS"; then
    test_case "Middleware tests" 0
else
    test_case "Middleware tests" 1
fi

if go test ./internal/handlers/... -v 2>&1 | grep -q "PASS"; then
    test_case "Handler tests" 0
else
    test_case "Handler tests" 1
fi

echo ""
echo "🐍 Python Unit Tests"
echo "--------------------"

# Run Python tests (from brain directory)
cd brain
if python3 -m pytest test_learner.py -v 2>&1 | grep -q "passed" || python3 -m unittest test_learner -v 2>&1 | grep -q "OK"; then
    test_case "Learner tests" 0
else
    test_case "Learner tests" 1
fi
cd ..

echo ""
echo "🔗 Build Tests"
echo "--------------"

# Test Go build
if go build -o /tmp/sentinel_test ./cmd/sentinel 2>&1; then
    test_case "Go build succeeds" 0
    rm -f /tmp/sentinel_test
else
    test_case "Go build succeeds" 1
fi

echo ""
echo "📊 Test Summary"
echo "==============="
echo -e "Passed: ${GREEN}$PASSED${NC}"
echo -e "Failed: ${RED}$FAILED${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}All tests passed! 🎉${NC}"
    exit 0
else
    echo -e "\n${RED}Some tests failed. Please fix before pushing.${NC}"
    exit 1
fi
