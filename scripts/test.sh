#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}   Apex Blockchain Test Suite        ${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""

# Run unit tests
echo -e "${YELLOW}Running unit tests...${NC}"
go test -v -race -cover ./pkg/...

echo ""
echo -e "${YELLOW}Running integration tests...${NC}"
go test -v -race -tags=integration ./tests/...

# Generate coverage report
echo ""
echo -e "${YELLOW}Generating coverage report...${NC}"
mkdir -p coverage
go test -coverprofile=coverage/coverage.out ./pkg/...
go tool cover -html=coverage/coverage.out -o coverage/coverage.html

echo ""
echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}All tests passed!${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""
echo -e "Coverage report: ${GREEN}coverage/coverage.html${NC}"
echo ""