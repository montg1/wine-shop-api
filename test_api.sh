#!/bin/bash

BASE_URL="http://localhost:8080/api"
EMAIL="test_runner@example.com"
PASSWORD="password123"

echo "🍷 Starting Wine Shop API Tests..."
echo "-----------------------------------"

# 1. Health Check
echo "1. Checking API Health..."
curl -s "$BASE_URL/health" | grep "ok" && echo " - OK" || echo " - FAILED"
echo ""

# 2. Register
echo "2. Registering User ($EMAIL)..."
curl -s -X POST "$BASE_URL/register" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}"
echo ""
echo ""

# 3. Login
echo "3. Logging in..."
LOGIN_RESPONSE=$(curl -s -X POST "$BASE_URL/login" \
  -H "Content-Type: application/json" \
  -d "{\"email\": \"$EMAIL\", \"password\": \"$PASSWORD\"}")

# Extract Token (Simple parsing)
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
  echo "❌ Login Failed! Response: $LOGIN_RESPONSE"
  exit 1
fi
echo "✅ Token received!"
echo ""

# 4. Admin: Create Product
echo "4. [Admin] Creating New Wine (Pinot Noir)..."
CREATE_PROD_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Pinot Noir", "description": "Elegant red", "price": 45.00, "stock": 10, "category": "Red"}')
echo $CREATE_PROD_RESPONSE
echo ""

# Extract Product ID (Simple parsing approach, assuming output format)
# This is a bit brittle without jq but sufficient for this demo
PID=$(echo $CREATE_PROD_RESPONSE | grep -o '"ID":[0-9]*' | head -1 | grep -o '[0-9]*')
echo "   -> Created Product ID: $PID"
echo ""

# 5. Public: List Products
echo "5. [Public] Listing Products..."
curl -s "$BASE_URL/products"
echo ""
echo ""

# 6. Cart: Add Item
echo "6. [User] Adding Product $PID to Cart..."
curl -s -X POST "$BASE_URL/cart" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"product_id\": $PID, \"quantity\": 2}"
echo ""
echo ""

# 7. Cart: View Cart
echo "7. [User] Viewing Cart..."
curl -s "$BASE_URL/cart" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 8. Order: Checkout
echo "8. [User] Checking Out..."
curl -s -X POST "$BASE_URL/orders" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 9. Order: History
echo "9. [User] Viewing Order History..."
curl -s "$BASE_URL/orders" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

echo "-----------------------------------"
echo "🎉 All tests completed!"
