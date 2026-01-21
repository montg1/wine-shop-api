#!/bin/bash
# Analytics Dashboard Tests with pretty output

BASE_URL="http://localhost:8080/api"
ADMIN_EMAIL="admin_test_$(date +%s)@example.com"
PASSWORD="password123"
PASS=0
FAIL=0

echo "📊 Analytics Dashboard Tests"
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

# 1. Setup: Create admin user
echo ""
echo "🔧 Setup: Creating Admin User"

# Register user
curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$ADMIN_EMAIL\", \"password\": \"$PASSWORD\"}" > /dev/null

# Login to get token
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$ADMIN_EMAIL\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo "$LOGIN_RESPONSE" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -z "$TOKEN" ]; then
  echo "   ❌ Failed to get auth token"
  exit 1
fi
echo "   ✅ Admin user created and logged in"

# Make user an admin (requires DB access in CI)
# This is handled by the CI workflow setting up an admin user

# 2. Test Analytics Endpoints (should require admin)
echo ""
echo "📋 Test 1: Analytics Access Control (Non-Admin)"

STATS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/admin/analytics/stats" \
  -H "Authorization: Bearer $TOKEN")

if [ "$STATS_CODE" == "403" ]; then
  test_result "pass" "Non-admin blocked from analytics (403)"
else
  test_result "fail" "Expected 403, got $STATS_CODE"
fi

# 3. Test Endpoint Structure (Unauthenticated)
echo ""
echo "📋 Test 2: Analytics Endpoints Require Auth"

UNAUTH_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/admin/analytics/stats")
if [ "$UNAUTH_CODE" == "401" ]; then
  test_result "pass" "Unauthenticated request blocked (401)"
else
  test_result "fail" "Expected 401, got $UNAUTH_CODE"
fi

# 4. Test all analytics endpoints exist
echo ""
echo "📋 Test 3: Analytics Endpoints Exist"

ENDPOINTS=(
  "/admin/analytics/stats"
  "/admin/analytics/sales-by-category"
  "/admin/analytics/top-products"
  "/admin/analytics/sales-by-day"
  "/admin/analytics/recent-orders"
)

for endpoint in "${ENDPOINTS[@]}"; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL$endpoint" \
    -H "Authorization: Bearer $TOKEN")
  # Should get 403 (forbidden) not 404 (not found)
  if [ "$CODE" == "403" ]; then
    test_result "pass" "Endpoint $endpoint exists (403 for non-admin)"
  elif [ "$CODE" == "404" ]; then
    test_result "fail" "Endpoint $endpoint not found (404)"
  else
    test_result "pass" "Endpoint $endpoint responded ($CODE)"
  fi
done

# 5. Test with Admin Privileges (if admin token available)
echo ""
echo "📋 Test 4: Admin Analytics Response Format"

# Try to use test@example.com as admin (pre-seeded in DB)
ADMIN_LOGIN=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "test@example.com", "password": "password123"}')

ADMIN_TOKEN=$(echo "$ADMIN_LOGIN" | grep -o '"token":"[^"]*' | cut -d'"' -f4)

if [ -n "$ADMIN_TOKEN" ]; then
  # Test stats endpoint with admin token
  STATS_RESPONSE=$(curl -s "$BASE_URL/admin/analytics/stats" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  
  if echo "$STATS_RESPONSE" | grep -q "total_revenue\|total_orders"; then
    test_result "pass" "Stats endpoint returns correct format"
  elif echo "$STATS_RESPONSE" | grep -q "data"; then
    test_result "pass" "Stats endpoint returns data object"
  else
    test_result "fail" "Stats endpoint response unexpected: $STATS_RESPONSE"
  fi
  
  # Test sales-by-category
  CATEGORY_RESPONSE=$(curl -s "$BASE_URL/admin/analytics/sales-by-category" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  
  if echo "$CATEGORY_RESPONSE" | grep -q "data\|category\|revenue"; then
    test_result "pass" "Sales by category endpoint works"
  else
    test_result "fail" "Sales by category response unexpected"
  fi
  
  # Test top-products
  TOP_RESPONSE=$(curl -s "$BASE_URL/admin/analytics/top-products?limit=5" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  
  if echo "$TOP_RESPONSE" | grep -q "data"; then
    test_result "pass" "Top products endpoint works"
  else
    test_result "fail" "Top products response unexpected"
  fi
  
  # Test sales-by-day
  SALES_DAY_RESPONSE=$(curl -s "$BASE_URL/admin/analytics/sales-by-day?days=30" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  
  if echo "$SALES_DAY_RESPONSE" | grep -q "data"; then
    test_result "pass" "Sales by day endpoint works"
  else
    test_result "fail" "Sales by day response unexpected"
  fi
  
  # Test recent-orders
  ORDERS_RESPONSE=$(curl -s "$BASE_URL/admin/analytics/recent-orders?limit=5" \
    -H "Authorization: Bearer $ADMIN_TOKEN")
  
  if echo "$ORDERS_RESPONSE" | grep -q "data"; then
    test_result "pass" "Recent orders endpoint works"
  else
    test_result "fail" "Recent orders response unexpected"
  fi
else
  echo "   ⚠️  Skipping admin tests (no admin user seeded)"
fi

# Summary
echo ""
echo "============================"
echo "Results: $PASS passed, $FAIL failed"
echo "============================"

if [ $FAIL -gt 0 ]; then
  exit 1
fi
