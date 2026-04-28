#!/bin/bash

BASE_URL="http://localhost:8081"

echo "=== Testing Profile Service API ==="
echo ""

# 1. Health check
echo "1. GET /healthz"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/healthz"
echo ""

# 2. Create profile
echo "2. POST /profiles"
CREATE_RESPONSE=$(curl -s -w "\n%{http_code}" -X POST "$BASE_URL/profiles" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Doe",
    "login": "johndoe",
    "email": "john@example.com",
    "bio": "Software developer",
    "image": "https://example.com/avatar.jpg"
  }')
echo "$CREATE_RESPONSE"
HTTP_CODE=$(echo "$CREATE_RESPONSE" | tail -n1)
PROFILE_ID=$(echo "$CREATE_RESPONSE" | grep -o '"id":"[^"]*"' | cut -d'"' -f4)
echo "Profile ID: $PROFILE_ID"
echo ""

# 3. Get profile by ID
echo "3. GET /profiles/:id"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/profiles/$PROFILE_ID"
echo ""

# 4. Get profile by login
echo "4. GET /profiles/login/:login"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/profiles/login/johndoe"
echo ""

# 5. Update profile
echo "5. PUT /profiles/:id"
curl -s -w "\nHTTP Status: %{http_code}\n" -X PUT "$BASE_URL/profiles/$PROFILE_ID" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Smith",
    "email": "john.smith@example.com",
    "bio": "Senior software developer"
  }'
echo ""

# 6. Verify update
echo "6. GET /profiles/:id (after update)"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/profiles/$PROFILE_ID"
echo ""

# 7. Create second profile for uniqueness test
echo "7. POST /profiles (second profile)"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/profiles" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Jane Doe",
    "login": "janedoe",
    "email": "jane@example.com",
    "bio": "Designer"
  }'
echo ""

# 8. Test duplicate login (should fail)
echo "8. POST /profiles (duplicate login - should fail)"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/profiles" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "John Duplicate",
    "login": "johndoe",
    "email": "john.duplicate@example.com"
  }'
echo ""

# 9. Test invalid data (should fail)
echo "9. POST /profiles (invalid data - should fail)"
curl -s -w "\nHTTP Status: %{http_code}\n" -X POST "$BASE_URL/profiles" \
  -H "Content-Type: application/json" \
  -d '{
    "name": "A",
    "login": "ab",
    "email": "invalid-email"
  }'
echo ""

# 10. Delete profile
echo "10. DELETE /profiles/:id"
curl -s -w "\nHTTP Status: %{http_code}\n" -X DELETE "$BASE_URL/profiles/$PROFILE_ID"
echo ""

# 11. Verify deletion
echo "11. GET /profiles/:id (after deletion - should fail)"
curl -s -w "\nHTTP Status: %{http_code}\n" "$BASE_URL/profiles/$PROFILE_ID"
echo ""

echo "=== Test Complete ==="
