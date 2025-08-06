#!/bin/bash

echo "📋 Test API Endpoints"
echo ""

BASE_URL="http://localhost:8080"

# Test Gateway Health
echo "🔍 Testing Gateway Health..."
curl -s "$BASE_URL/health" | python3 -m json.tool
echo ""

# Test Services Health
echo "🔍 Testing Services Health..."
curl -s "$BASE_URL/services/health" | python3 -m json.tool
echo ""

# Test User Service
echo "👤 Testing User Service..."
echo "GET /users:"
curl -s "$BASE_URL/users" | python3 -m json.tool
echo ""

echo "GET /users/1:"
curl -s "$BASE_URL/users/1" | python3 -m json.tool
echo ""

# Test Product Service
echo "📦 Testing Product Service..."
echo "GET /products:"
curl -s "$BASE_URL/products" | python3 -m json.tool
echo ""

echo "GET /products/1:"
curl -s "$BASE_URL/products/1" | python3 -m json.tool
echo ""

# Test Create User
echo "POST /users (Create new user):"
curl -s -X POST "$BASE_URL/users" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test User","email":"test@example.com"}' | python3 -m json.tool
echo ""

# Test Create Product
echo "POST /products (Create new product):"
curl -s -X POST "$BASE_URL/products" \
  -H "Content-Type: application/json" \
  -d '{"name":"Test Product","description":"Test Description","price":100000,"stock":5}' | python3 -m json.tool
echo ""

echo "✅ Test completed!"
