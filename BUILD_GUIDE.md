# Apex Blockchain - Complete Build and Setup Guide

## 🚀 Quick Start

This guide will help you build and run the Apex blockchain from source, and fix common issues including the Service Worker storage error.

## 📋 Prerequisites

### System Requirements
- **OS**: Windows 10/11, Ubuntu 20.04+, or macOS 10.15+
- **RAM**: Minimum 8GB (16GB recommended)
- **Disk**: 20GB free space
- **Go**: Version 1.21 or higher

### Installing Go (Windows)

1. Download Go from https://golang.org/dl/
2. Install and add to PATH:
```powershell
# Add to environment variables
setx GOPATH "%USERPROFILE%\go"
setx PATH "%PATH%;%GOPATH%\bin;C:\Program Files\Go\bin"
```

3. Verify installation:
```powershell
go version
```

## 🔧 Building from Source

### 1. Clone or Set Up Project

If you have the source code:
```powershell
cd C:\Users\Administrator\APEX
```

If cloning from git:
```powershell
git clone https://github.com/apex-blockchain/apex.git
cd apex
```

### 2. Install Dependencies

```powershell
go mod download
go mod tidy
```

### 3. Build All Binaries

#### Option A: Using Build Script (Recommended)

**On Windows:**
```powershell
# Make script executable and run
.\scripts\build.sh
```

**On Linux/Mac:**
```bash
chmod +x scripts/build.sh
./scripts/build.sh
```

#### Option B: Manual Build

```powershell
# Create bin directory
mkdir bin

# Build main node
go build -o bin/apex.exe ./cmd/apex

# Build CLI tool
go build -o bin/apexctl.exe ./cmd/apexctl

# Build genesis tool
go build -o bin/genesis.exe ./cmd/genesis
```

### 4. Verify Build

```powershell
.\bin\apex.exe --help
.\bin\apexctl.exe version
.\bin\genesis.exe --help
```

## 🐛 Fixing the Service Worker Storage Error

The error you're seeing:
```
[5752:0123/020746.965:ERROR:service_worker_storage.cc(2016)] Failed to delete the database: Database IO error
```

This is related to Chrome/Chromium's service worker cache and can be fixed:

### Solution 1: Clear Browser Data (Quick Fix)

1. Open Chrome
2. Press `Ctrl + Shift + Delete`
3. Select "All time"
4. Check:
   - Cached images and files
   - Site settings
5. Click "Clear data"

### Solution 2: Delete Service Worker Cache Manually

**On Windows:**
```powershell
# Close all Chrome instances first

# Navigate to Chrome user data
cd "$env:LOCALAPPDATA\Google\Chrome\User Data\Default"

# Delete Service Worker directory
Remove-Item -Recurse -Force "Service Worker"

# Also clear IndexedDB
Remove-Item -Recurse -Force "IndexedDB"
```

**On Linux:**
```bash
rm -rf ~/.config/google-chrome/Default/Service\ Worker
rm -rf ~/.config/google-chrome/Default/IndexedDB
```

### Solution 3: Run Chrome with Clean Profile

```powershell
# Create new profile directory
mkdir C:\temp\chrome-profile

# Run Chrome with clean profile
"C:\Program Files\Google\Chrome\Application\chrome.exe" --user-data-dir="C:\temp\chrome-profile"
```

### Solution 4: For Development - Disable Service Workers

If this is a development environment, you can disable service workers:

1. Open Chrome DevTools (`F12`)
2. Go to Application tab
3. Click "Service Workers"
4. Check "Update on reload"
5. Or unregister all service workers

## 🏃 Running the Blockchain

### 1. Initialize Configuration

```powershell
# Copy default config
copy config\config.yaml.example config\config.yaml

# Or create custom config
notepad config\config.yaml
```

### 2. Generate Genesis Block

```powershell
.\bin\genesis.exe --chain-id apex-mainnet-1 --output config\genesis.json
```

### 3. Run the Node

#### As Full Node:
```powershell
.\bin\apex.exe --config config\config.yaml
```

#### As Validator:
```powershell
# First, generate validator keys
.\bin\apexctl.exe keys generate --output validator_key.json

# Edit config.yaml to enable validator mode
# Then start node
.\bin\apex.exe --config config\config.yaml
```

### 4. Verify Node is Running

Open another terminal:
```powershell
# Check block number
curl -X POST http://localhost:8545 -H "Content-Type: application/json" -d "{\"jsonrpc\":\"2.0\",\"method\":\"apex_blockNumber\",\"params\":[],\"id\":1}"
```

## 🔍 Testing the API

### Using PowerShell:

```powershell
# Get current block number
$body = @{
    jsonrpc = "2.0"
    method = "apex_blockNumber"
    params = @()
    id = 1
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8545" -Method Post -Body $body -ContentType "application/json"

# Get account balance
$body = @{
    jsonrpc = "2.0"
    method = "apex_getBalance"
    params = @("0x0000000000000000000000000000000000000001")
    id = 1
} | ConvertTo-Json

Invoke-RestMethod -Uri "http://localhost:8545" -Method Post -Body $body -ContentType "application/json"
```

### Using curl:

```bash
# Get validators
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"apex_getValidators","params":[],"id":1}'
```

## 📦 Common Build Issues

### Issue 1: "go: command not found"

**Solution:**
- Reinstall Go and ensure it's in PATH
- Restart terminal after installation

### Issue 2: "package not found"

**Solution:**
```powershell
go clean -modcache
go mod download
go mod tidy
```

### Issue 3: "cannot find module"

**Solution:**
- Check go.mod file exists
- Run `go mod init github.com/apex-blockchain/apex` if missing

### Issue 4: Build errors with dependencies

**Solution:**
```powershell
# Update all dependencies
go get -u ./...
go mod tidy
```

### Issue 5: Port 8545 already in use

**Solution:**
```powershell
# Find process using port
netstat -ano | findstr :8545

# Kill the process (replace PID)
taskkill /PID <PID> /F

# Or change port in config.yaml
```

## 🧪 Running Tests

```powershell
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package
go test ./pkg/consensus/...

# Verbose output
go test -v ./...
```

## 📊 Monitoring

### Enable Metrics

Edit `config/config.yaml`:
```yaml
metrics:
  enabled: true
  port: 9091
```

Access metrics at: `http://localhost:9091/metrics`

### View Logs

```powershell
# Real-time logs
Get-Content logs\apex.log -Wait -Tail 50

# Filter errors only
Get-Content logs\apex.log | Select-String "ERROR"
```

## 🐳 Docker Alternative

If you prefer Docker:

```powershell
# Build image
docker build -t apex-blockchain .

# Run node
docker run -d -p 8545:8545 -p 30303:30303 `
  -v ${PWD}\data:/data `
  -v ${PWD}\config:/config `
  apex-blockchain
```

## 🔐 Security Checklist

- [ ] Change default RPC port
- [ ] Enable TLS for RPC
- [ ] Use firewall to restrict access
- [ ] Keep validator keys secure
- [ ] Regular backups
- [ ] Monitor for suspicious activity

## 📚 Next Steps

1. **Set up a validator**: Follow the validator guide
2. **Join testnet**: Connect to test network
3. **Stake tokens**: Start earning rewards
4. **Monitor performance**: Use Prometheus + Grafana
5. **Join community**: Discord, Telegram

## 🆘 Getting Help

- **Documentation**: https://docs.apex.network
- **Discord**: https://discord.gg/apex
- **GitHub Issues**: https://github.com/apex-blockchain/apex/issues
- **Email**: support@apex.network

## ✅ Verification Checklist

After building, verify:

- [ ] All 3 binaries created in `bin/` folder
- [ ] Node starts without errors
- [ ] RPC responds on port 8545
- [ ] Can query blockchain height
- [ ] Can query validators
- [ ] Logs are being written

## 🎯 Performance Tips

1. **Use SSD** for blockchain data
2. **Allocate sufficient RAM** (16GB+)
3. **Fast internet** for sync
4. **Keep system updated**
5. **Monitor resource usage**

---

**Build Date**: January 2026  
**Version**: 1.0.0  
**License**: MIT

For the latest updates, visit: https://github.com/apex-blockchain/apex
