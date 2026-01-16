#!/bin/bash
# Security tests for CI - no long waits

BASE_URL="http://localhost:8080/api"
EMAIL="security_ci_$(date +%s)@example.com"
PASSWORD="password123"
PASS=0
FAIL=0

echo "🔐 Security Tests (CI Mode)"
echo "============================"

# Test function
test_result() {
  if [ "$1" == "pass" ]; then
    echo "   ✅ $2"
    ((PASS++))
  else
    echo "   ❌ $2"
    ((FAIL++))
  fi
}

# Register & Login
echo ""
echo "📊 Setup: Register & Login"
curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}" > /dev/null

LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -n "$TOKEN" ]; then
  echo "   ✅ Login successful"
else
  echo "   ❌ Login failed!"
  exit 1
fi

# Test 1: RBAC - Admin Access
echo ""
echo "📊 Test 1: RBAC"
ADMIN_CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/admin/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test", "price": 10}')

if [ "$ADMIN_CODE" == "403" ]; then
  test_result "pass" "Non-admin blocked from admin routes (403)"
else
  test_result "fail" "Non-admin accessed admin route (HTTP $ADMIN_CODE)"
fi

# Test 2: Protected Routes Without Token
echo ""
echo "📊 Test 2: JWT Required"
CART_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/cart")
ME_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/me")
ORDERS_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/orders")

[ "$CART_CODE" == "401" ] && test_result "pass" "/cart requires auth (401)" || test_result "fail" "/cart no auth (HTTP $CART_CODE)"
[ "$ME_CODE" == "401" ] && test_result "pass" "/me requires auth (401)" || test_result "fail" "/me no auth (HTTP $ME_CODE)"
[ "$ORDERS_CODE" == "401" ] && test_result "pass" "/orders requires auth (401)" || test_result "fail" "/orders no auth (HTTP $ORDERS_CODE)"

# Test 3: Invalid JWT
echo ""
echo "📊 Test 3: Invalid JWT"
INVALID_CODE=$(curl -s -o /dev/null -w "%{http_code}" "$BASE_URL/cart" \
  -H "Authorization: Bearer INVALID.TOKEN.HERE")

[ "$INVALID_CODE" == "401" ] && test_result "pass" "Invalid JWT rejected (401)" || test_result "fail" "Invalid JWT accepted (HTTP $INVALID_CODE)"

# Test 4: SQL Injection
echo ""
echo "📊 Test 4: SQL Injection"
SQL_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d '{"email": "admin'\''--", "password": "x"}')

if echo "$SQL_RESPONSE" | grep -q "token"; then
  test_result "fail" "SQL Injection may have worked!"
else
  test_result "pass" "SQL Injection blocked"
fi

# Test 5: Rate Limiting (quick check - 11 requests)
echo ""
echo "📊 Test 5: Rate Limiting"
GOT_429=false
for i in {1..11}; do
  CODE=$(curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d '{"email": "fake@test.com", "password": "wrong"}')
  if [ "$CODE" == "429" ]; then
    GOT_429=true
    break
  fi
done

[ "$GOT_429" == "true" ] && test_result "pass" "Rate limiting triggered (429)" || test_result "fail" "Rate limiting not triggered"

# Summary
echo ""
echo "============================"
echo "Results: $PASS passed, $FAIL failed"
echo "============================"

if [ $FAIL -gt 0 ]; then
  exit 1
fi
