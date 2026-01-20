#!/bin/bash
# Integration Tests with pretty output

BASE_URL="http://localhost:8080/api"
EMAIL="test_runner_$(date +%s)@example.com"
PASSWORD="password123"
PASS=0
FAIL=0

echo "🔗 Integration Tests"
echo "============================"

# Test helper function
test_result() {
  if [ "$1" == "pass" ]; then
    echo "   ✅ $2"
    ((PASS++))
  else
    echo "   ❌ $2"
    ((FAIL++))
  fi
}

# 1. Health Check
echo ""
echo "📋 Test 1: API Health Check"
HEALTH=$(curl -s "$BASE_URL/health")
if echo "$HEALTH" | grep -q "ok"; then
  test_result "pass" "API is healthy"
else
  test_result "fail" "API health check failed"
fi

# 2. User Registration
echo ""
echo "📋 Test 2: User Registration"
REG_RESPONSE=$(curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

if echo "$REG_RESPONSE" | grep -q "registration success\|already exists"; then
  test_result "pass" "User registration works"
else
  test_result "fail" "User registration failed"
fi

# 3. User Login
echo ""
echo "📋 Test 3: User Login"
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -n "$TOKEN" ]; then
  test_result "pass" "Login successful, token received"
else
  test_result "fail" "Login failed - no token"
  echo "============================"
  echo "Results: $PASS passed, $FAIL failed"
  echo "============================"
  exit 1
fi

# 4. Get Products (Public)
echo ""
echo "📋 Test 4: Public Endpoints"
PRODUCTS=$(curl -s "$BASE_URL/products")
if echo "$PRODUCTS" | grep -q "data"; then
  test_result "pass" "GET /products returns data"
else
  test_result "fail" "GET /products failed"
fi

# Search
SEARCH=$(curl -s "$BASE_URL/products?search=wine")
if echo "$SEARCH" | grep -q "data"; then
  test_result "pass" "GET /products?search works"
else
  test_result "fail" "Search failed"
fi

# Filter
FILTER=$(curl -s "$BASE_URL/products?category=Red")
if echo "$FILTER" | grep -q "data"; then
  test_result "pass" "GET /products?category works"
else
  test_result "fail" "Filter failed"
fi

# 5. Protected Endpoints
echo ""
echo "📋 Test 5: Protected Endpoints (User)"

# View Cart
CART=$(curl -s "$BASE_URL/cart" -H "Authorization: Bearer $TOKEN")
if echo "$CART" | grep -q "data\|ID"; then
  test_result "pass" "GET /cart works"
else
  test_result "fail" "GET /cart failed"
fi

# View Orders
ORDERS=$(curl -s "$BASE_URL/orders" -H "Authorization: Bearer $TOKEN")
if echo "$ORDERS" | grep -q "data"; then
  test_result "pass" "GET /orders works"
else
  test_result "fail" "GET /orders failed"
fi

# Get User Info
ME=$(curl -s "$BASE_URL/me" -H "Authorization: Bearer $TOKEN")
if echo "$ME" | grep -q "data\|email"; then
  test_result "pass" "GET /me works"
else
  test_result "fail" "GET /me failed"
fi

# 6. Admin Access Control
echo ""
echo "📋 Test 6: Admin Access Control"
ADMIN_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/admin/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "price": 10}')

if [ "$ADMIN_RESPONSE" == "403" ]; then
  test_result "pass" "Non-admin blocked from admin routes (403)"
else
  test_result "fail" "Admin access control issue (HTTP $ADMIN_RESPONSE)"
fi

# Summary
echo ""
echo "============================"
echo "Results: $PASS passed, $FAIL failed"
echo "============================"

if [ $FAIL -gt 0 ]; then
  exit 1
fi
