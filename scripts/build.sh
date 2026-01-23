#!/bin/bash

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}   Apex Blockchain Build Script      ${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""

# Get version from git tag or use default
VERSION=$(git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(date -u '+%Y-%m-%d_%H:%M:%S')
GIT_COMMIT=$(git rev-parse --short HEAD 2>/dev/null || echo "unknown")

LDFLAGS="-X main.Version=${VERSION} -X main.BuildTime=${BUILD_TIME} -X main.GitCommit=${GIT_COMMIT}"

# Create bin directory
mkdir -p bin

echo -e "${YELLOW}Building Apex node...${NC}"
go build -ldflags "${LDFLAGS}" -o bin/apex ./cmd/apex
echo -e "${GREEN}✓ Apex node built successfully${NC}"

echo ""
echo -e "${YELLOW}Building apexctl CLI tool...${NC}"
go build -ldflags "${LDFLAGS}" -o bin/apexctl ./cmd/apexctl
echo -e "${GREEN}✓ apexctl built successfully${NC}"

echo ""
echo -e "${YELLOW}Building genesis tool...${NC}"
go build -ldflags "${LDFLAGS}" -o bin/genesis ./cmd/genesis
echo -e "${GREEN}✓ Genesis tool built successfully${NC}"

echo ""
echo -e "${GREEN}======================================${NC}"
echo -e "${GREEN}Build completed successfully!${NC}"
echo -e "${GREEN}======================================${NC}"
echo ""
echo -e "Version: ${VERSION}"
echo -e "Build Time: ${BUILD_TIME}"
echo -e "Git Commit: ${GIT_COMMIT}"
echo ""
echo -e "Binaries location:"
echo -e "  - ${GREEN}bin/apex${NC} - Main node executable"
echo -e "  - ${GREEN}bin/apexctl${NC} - CLI management tool"
echo -e "  - ${GREEN}bin/genesis${NC} - Genesis block generator"
echo ""