#!/bin/bash

# Test script for Go XLSX Upload API
# This script demonstrates how to interact with the API

API_URL="http://localhost:8080"
API_KEY="secret123"

echo "========================================="
echo "Go XLSX Upload API - Test Script"
echo "========================================="
echo ""

# Function to print colored output
print_success() {
    echo -e "\033[0;32m✓ $1\033[0m"
}

print_error() {
    echo -e "\033[0;31m✗ $1\033[0m"
}

print_info() {
    echo -e "\033[0;34mℹ $1\033[0m"
}

# Test 1: Health Check
echo "Test 1: Health Check"
echo "-------------------"
response=$(curl -s "$API_URL/healthz")
if echo "$response" | grep -q "ok"; then
    print_success "Health check passed"
    echo "Response: $response"
else
    print_error "Health check failed"
    echo "Response: $response"
fi
echo ""

# Test 2: Upload XLSX File (requires a sample file)
echo "Test 2: Upload XLSX File"
echo "-------------------"
print_info "Note: This test requires a sample.xlsx file in the current directory"
if [ -f "sample.xlsx" ]; then
    response=$(curl -s -X POST "$API_URL/v1/uploads" \
        -H "X-API-Key: $API_KEY" \
        -F "file=@sample.xlsx")

    if echo "$response" | grep -q "uploadId"; then
        print_success "File upload successful"
        echo "Response: $response"
    else
        print_error "File upload failed"
        echo "Response: $response"
    fi
else
    print_info "Skipping upload test - sample.xlsx not found"
    print_info "To test uploads, create a sample.xlsx file with:"
    echo "   - Header row (e.g., Name, Email, Age)"
    echo "   - At least one data row"
fi
echo ""

# Test 3: List Records
echo "Test 3: List Records"
echo "-------------------"
response=$(curl -s "$API_URL/v1/records?limit=10&offset=0" \
    -H "X-API-Key: $API_KEY")

if echo "$response" | grep -q "records"; then
    print_success "List records successful"
    echo "Response: $response" | jq . 2>/dev/null || echo "$response"
else
    print_error "List records failed"
    echo "Response: $response"
fi
echo ""

# Test 4: Pagination
echo "Test 4: Pagination (limit=5, offset=0)"
echo "-------------------"
response=$(curl -s "$API_URL/v1/records?limit=5&offset=0" \
    -H "X-API-Key: $API_KEY")

if echo "$response" | grep -q "limit"; then
    print_success "Pagination test successful"
    echo "Response: $response" | jq . 2>/dev/null || echo "$response"
else
    print_error "Pagination test failed"
    echo "Response: $response"
fi
echo ""

# Test 5: Missing API Key
echo "Test 5: Missing API Key (should fail)"
echo "-------------------"
response=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$API_URL/v1/records?limit=10")
http_code=$(echo "$response" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$http_code" = "401" ]; then
    print_success "API key validation working correctly (401 returned)"
else
    print_error "API key validation not working (expected 401, got $http_code)"
fi
echo ""

# Test 6: Invalid API Key
echo "Test 6: Invalid API Key (should fail)"
echo "-------------------"
response=$(curl -s -w "\nHTTP_CODE:%{http_code}" "$API_URL/v1/records?limit=10" \
    -H "X-API-Key: wrongkey")
http_code=$(echo "$response" | grep "HTTP_CODE" | cut -d: -f2)

if [ "$http_code" = "401" ]; then
    print_success "Invalid API key rejected correctly (401 returned)"
else
    print_error "Invalid API key not rejected (expected 401, got $http_code)"
fi
echo ""

echo "========================================="
echo "Tests Complete!"
echo "========================================="
