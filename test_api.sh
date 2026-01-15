#!/bin/bash

BASE_URL="http://localhost:8080/api"
EMAIL="test_runner_$(date +%s)@example.com"
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

# Extract Token
TOKEN=$(echo $LOGIN_RESPONSE | grep -o '"token":"[^"]*' | grep -o '[^"]*$')

if [ -z "$TOKEN" ]; then
  echo "❌ Login Failed! Response: $LOGIN_RESPONSE"
  exit 1
fi
echo "✅ Token received!"
echo ""

# 4. Admin: Create Product
echo "4. [Admin] Creating New Wine (Test Merlot)..."
CREATE_PROD_RESPONSE=$(curl -s -X POST "$BASE_URL/admin/products" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"name": "Test Merlot", "description": "Smooth red wine", "price": 38.00, "stock": 20, "category": "Red"}')
echo $CREATE_PROD_RESPONSE
echo ""

PID=$(echo $CREATE_PROD_RESPONSE | grep -o '"ID":[0-9]*' | head -1 | grep -o '[0-9]*')
echo "   -> Created Product ID: $PID"
echo ""

# 5. Public: List Products
echo "5. [Public] Listing Products..."
curl -s "$BASE_URL/products"
echo ""
echo ""

# 6. Search Products
echo "6. [Public] Searching for 'merlot'..."
curl -s "$BASE_URL/products?search=merlot"
echo ""
echo ""

# 7. Filter by Category
echo "7. [Public] Filtering by 'Red' category..."
curl -s "$BASE_URL/products?category=Red"
echo ""
echo ""

# 8. Cart: Add Item
echo "8. [User] Adding Product $PID to Cart..."
curl -s -X POST "$BASE_URL/cart" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"product_id\": $PID, \"quantity\": 2}"
echo ""
echo ""

# 9. Cart: View Cart
echo "9. [User] Viewing Cart..."
curl -s "$BASE_URL/cart" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 10. Order: Checkout
echo "10. [User] Checking Out..."
curl -s -X POST "$BASE_URL/orders" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 11. Order: History
echo "11. [User] Viewing Order History..."
curl -s "$BASE_URL/orders" \
  -H "Authorization: Bearer $TOKEN"
echo ""
echo ""

# 12. Create Review
echo "12. [User] Creating Review for Product $PID..."
curl -s -X POST "$BASE_URL/products/$PID/reviews" \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"rating": 5, "comment": "Excellent wine!"}'
echo ""
echo ""

# 13. Get Product Reviews
echo "13. [Public] Getting Reviews for Product $PID..."
curl -s "$BASE_URL/products/$PID/reviews"
echo ""
echo ""

echo "-----------------------------------"
echo "🎉 All tests completed!"
