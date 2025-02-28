#!/bin/bash

# Frontend Verification Workflow Script
# This script automates the process of building and testing the frontend

# Exit on any error
set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Working directory
REPO_ROOT="/mnt/c/Users/Josh/GolandProjects/pf2eEngine"
FRONTEND_DIR="$REPO_ROOT/frontend"
SERVER_LOG="$REPO_ROOT/server.log"
SERVER_PID_FILE="$REPO_ROOT/server.pid"

# Cleanup function
cleanup() {
  echo -e "${YELLOW}Cleaning up...${NC}"
  
  # Kill server if running
  if [ -f "$SERVER_PID_FILE" ]; then
    kill $(cat "$SERVER_PID_FILE") 2>/dev/null || true
    rm "$SERVER_PID_FILE"
  else
    pkill pf2e-server 2>/dev/null || true
  fi
  
  # Remove log file
  rm -f "$SERVER_LOG"
  
  # Verify server is down
  if pgrep pf2e-server > /dev/null; then
    echo -e "${RED}Warning: Server process still running. Forcing kill...${NC}"
    pkill -9 pf2e-server 2>/dev/null || true
  fi
  
  echo -e "${GREEN}Cleanup complete${NC}"
}

# Register cleanup on script exit
trap cleanup EXIT

# Step 1: Clean up any existing processes
echo -e "${YELLOW}Step 1: Cleaning up existing processes...${NC}"
pkill pf2e-server 2>/dev/null || true
sleep 1

# Step 2: Build the frontend
echo -e "${YELLOW}Step 2: Building frontend...${NC}"
cd "$FRONTEND_DIR"
npm run build
if [ $? -ne 0 ]; then
  echo -e "${RED}Frontend build failed!${NC}"
  exit 1
fi
echo -e "${GREEN}Frontend built successfully${NC}"

# Step 3: Start the server
echo -e "${YELLOW}Step 3: Starting server...${NC}"
cd "$REPO_ROOT"
./pf2e-server > "$SERVER_LOG" 2>&1 &
SERVER_PID=$!
echo $SERVER_PID > "$SERVER_PID_FILE"
echo -e "${GREEN}Server started with PID $SERVER_PID${NC}"

# Step 4: Wait for server to initialize
echo -e "${YELLOW}Step 4: Waiting for server to initialize...${NC}"
sleep 3

# Step 5: Check if server is running
echo -e "${YELLOW}Step 5: Verifying server is running...${NC}"
if ! ps -p $SERVER_PID > /dev/null; then
  echo -e "${RED}Server failed to start! Check logs:${NC}"
  cat "$SERVER_LOG"
  exit 1
fi
echo -e "${GREEN}Server is running${NC}"

# Step 6: Display test URLs
echo -e "${YELLOW}Step 6: Testing URLs${NC}"
echo -e "${GREEN}Main application: http://localhost:8080${NC}"
echo -e "${GREEN}Test page: http://localhost:8080/test.html${NC}"
echo ""
echo -e "${YELLOW}=====================${NC}"
echo -e "${YELLOW}Verification Complete${NC}"
echo -e "${YELLOW}=====================${NC}"
echo ""
echo -e "Keep the server running for manual testing or press Ctrl+C to stop"

# Wait for user signal
wait $SERVER_PID