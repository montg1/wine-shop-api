#!/bin/bash

BASE_URL="http://localhost:8080/api"
EMAIL="security_test_$(date +%s)@example.com"
PASSWORD="password123"

echo "🔐 Starting Wine Shop Security Tests..."
echo "========================================"

# -----------------------------
# Test 1: Rate Limiting (Auth)
# -----------------------------
echo ""
echo "📊 TEST 1: Rate Limiting on Login (should block after 10 attempts)"
echo "-------------------------------------------------------------------"

for i in {1..12}; do
  RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"email": "fake@example.com", "password": "wrong"}')
  
  if [ "$RESPONSE" == "429" ]; then
    echo "   Attempt $i: BLOCKED (429 Too Many Requests) ✅"
  else
    echo "   Attempt $i: HTTP $RESPONSE"
  fi
done

# Wait for rate limit to reset
echo ""
echo "⏳ Waiting 60 seconds for rate limit reset..."
sleep 60

# -----------------------------
# Test 2: Register & Login
# -----------------------------
echo ""
echo "📊 TEST 2: User Registration & Login"
echo "-------------------------------------"

echo "   Registering user: $EMAIL"
curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" | head -c 100
echo ""

echo "   Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
  echo "   ❌ Login failed!"
  exit 1
fi
echo "   ✅ Token received!"

# -----------------------------
# Test 3: RBAC - Admin Access
# -----------------------------
echo ""
echo "📊 TEST 3: RBAC - Non-admin trying admin routes"
echo "------------------------------------------------"

ADMIN_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/admin/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Hacked Wine", "price": 0}')

if [ "$ADMIN_RESPONSE" == "403" ]; then
  echo "   ✅ Admin route blocked for non-admin (403 Forbidden)"
else
  echo "   ❌ SECURITY ISSUE: Non-admin accessed admin route! (HTTP $ADMIN_RESPONSE)"
fi

# -----------------------------
# Test 4: JWT Required Routes
# -----------------------------
echo ""
echo "📊 TEST 4: Protected Routes Without Token"
echo "------------------------------------------"

CART_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/cart")
ME_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/me")
ORDERS_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/orders")

if [ "$CART_RESPONSE" == "401" ]; then
  echo "   ✅ /cart blocked without token (401)"
else
  echo "   ❌ /cart accessible without token! (HTTP $CART_RESPONSE)"
fi

if [ "$ME_RESPONSE" == "401" ]; then
  echo "   ✅ /me blocked without token (401)"
else
  echo "   ❌ /me accessible without token! (HTTP $ME_RESPONSE)"
fi

if [ "$ORDERS_RESPONSE" == "401" ]; then
  echo "   ✅ /orders blocked without token (401)"
else
  echo "   ❌ /orders accessible without token! (HTTP $ORDERS_RESPONSE)"
fi

# -----------------------------
# Test 5: Invalid JWT
# -----------------------------
echo ""
echo "📊 TEST 5: Invalid JWT Token"
echo "-----------------------------"

INVALID_JWT="eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.INVALID.SIGNATURE"
INVALID_RESPONSE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/cart" \
  -H "Authorization: Bearer $INVALID_JWT")

if [ "$INVALID_RESPONSE" == "401" ]; then
  echo "   ✅ Invalid JWT rejected (401)"
else
  echo "   ❌ Invalid JWT accepted! (HTTP $INVALID_RESPONSE)"
fi

# -----------------------------
# Test 6: SQL Injection Attempt
# -----------------------------
echo ""
echo "📊 TEST 6: SQL Injection Attempt"
echo "---------------------------------"

SQL_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "admin@test.com'\'' OR 1=1--", "password": "anything"}')

if echo "$SQL_RESPONSE" | grep -q "token"; then
  echo "   ❌ CRITICAL: SQL Injection might have worked!"
else
  echo "   ✅ SQL Injection blocked (no token returned)"
fi

# -----------------------------
# Test 7: XSS in Review Comment
# -----------------------------
echo ""
echo "📊 TEST 7: XSS Prevention (stored in DB)"
echo "-----------------------------------------"

# First get a product ID
PRODUCTS=$(curl -s "$BASE_URL/products")
PID=$(echo $PRODUCTS | grep -o '"ID":[0-9]*' | head -1 | grep -o '[0-9]*')

if [ -n "$PID" ]; then
  XSS_PAYLOAD='{"rating": 5, "comment": "<script>alert(\"xss\")</script>"}'
  XSS_RESPONSE=$(curl -s -X POST "$BASE_URL/products/$PID/reviews" \
    -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d "$XSS_PAYLOAD")
  
  if echo "$XSS_RESPONSE" | grep -q "<script>"; then
    echo "   ⚠️  XSS payload stored (frontend should escape output)"
  else
    echo "   ✅ XSS payload handled"
  fi
else
  echo "   ⚠️  No products found to test"
fi

# -----------------------------
# Summary
# -----------------------------
echo ""
echo "========================================"
echo "🔐 Security Tests Completed!"
echo ""
echo "Coverage:"
echo "   ✓ Rate Limiting"
echo "   ✓ Role-Based Access Control (RBAC)"
echo "   ✓ JWT Authentication"
echo "   ✓ Protected Routes"
echo "   ✓ Invalid Token Handling"
echo "   ✓ SQL Injection Prevention"
echo "   ✓ XSS Prevention"
echo "========================================"
