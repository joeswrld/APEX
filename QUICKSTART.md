# Apex Blockchain - Quick Fix Guide

## 🔥 Quick Fix for Service Worker Error

If you see this error:
```
[5752:0123/020746.965:ERROR:service_worker_storage.cc(2016)] Failed to delete the database: Database IO error
```

### Windows - Instant Fix

Run this script:
```powershell
.\scripts\fix-service-worker.bat
```

Or manually:
```powershell
# Close Chrome
taskkill /F /IM chrome.exe

# Clear cache
Remove-Item -Recurse -Force "$env:LOCALAPPDATA\Google\Chrome\User Data\Default\Service Worker"
Remove-Item -Recurse -Force "$env:LOCALAPPDATA\Google\Chrome\User Data\Default\IndexedDB"
```

### Linux/Mac - Instant Fix

```bash
# Close Chrome
pkill chrome

# Clear cache
rm -rf ~/.config/google-chrome/Default/Service\ Worker
rm -rf ~/.config/google-chrome/Default/IndexedDB
```

## 🚀 Build and Run (2 Commands)

### Windows

```powershell
# 1. Build everything (includes fix)
.\scripts\build.bat

# 2. Run the node
.\bin\apex.exe
```

### Linux/Mac

```bash
# 1. Build everything
chmod +x scripts/build.sh && ./scripts/build.sh

# 2. Run the node
./bin/apex
```

## ✅ Verify It Works

```powershell
# Test the RPC endpoint
curl -X POST http://localhost:8545 -H "Content-Type: application/json" -d '{"jsonrpc":"2.0","method":"apex_blockNumber","params":[],"id":1}'
```

Expected response:
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "blockNumber": 0
  }
}
```

## 📁 Project Structure

```
apex/
├── bin/                    # Built binaries
│   ├── apex.exe           # Main node
│   ├── apexctl.exe        # CLI tool
│   └── genesis.exe        # Genesis generator
├── cmd/                    # Command-line applications
│   ├── apex/              # Main node
│   ├── apexctl/           # Management CLI
│   └── genesis/           # Genesis tool
├── pkg/                    # Core packages
│   ├── api/               # JSON-RPC & gRPC
│   ├── consensus/         # DPoS consensus
│   ├── core/              # Blockchain core
│   ├── crypto/            # Cryptography
│   ├── mempool/           # Transaction pool
│   ├── network/           # P2P networking
│   ├── staking/           # Staking logic
│   ├── storage/           # Database
│   └── types/             # Common types
├── config/                 # Configuration
│   ├── config.yaml        # Node config
│   └── genesis.json       # Genesis block
├── scripts/               # Build scripts
│   ├── build.bat          # Windows build (with fix)
│   ├── build.sh           # Linux/Mac build
│   └── fix-service-worker.bat  # Error fix only
└── BUILD_GUIDE.md         # Detailed guide
```

## 🎯 Common Tasks

### Create a Validator

```powershell
# 1. Generate keys
.\bin\apexctl.exe keys generate --output validator_key.json

# 2. Create validator
.\bin\apexctl.exe validator create --moniker "My Validator" --self-stake 100000
```

### Stake Tokens

```powershell
.\bin\apexctl.exe stake delegate --validator <address> --amount 1000
```

### Query Blockchain

```powershell
# Get validators
.\bin\apexctl.exe query validators

# Get balance
.\bin\apexctl.exe query balance --address <address>
```

## 🐛 Troubleshooting

### Build fails with "go: command not found"
```powershell
# Install Go from https://golang.org/dl/
# Add to PATH and restart terminal
```

### Port 8545 already in use
```powershell
# Find and kill process
netstat -ano | findstr :8545
taskkill /PID <PID> /F
```

### Dependencies error
```powershell
go clean -modcache
go mod download
go mod tidy
```

## 📚 Documentation

- **Full Build Guide**: `BUILD_GUIDE.md`
- **Deployment Guide**: `DEPLOYMENT.md`
- **API Reference**: Check JSON-RPC methods in `pkg/api/jsonrpc/handlers.go`
- **Configuration**: See `config/config.yaml`

## 🆘 Need Help?

1. **Check BUILD_GUIDE.md** for detailed instructions
2. **Run the fix script**: `.\scripts\fix-service-worker.bat`
3. **GitHub Issues**: https://github.com/apex-blockchain/apex/issues
4. **Discord**: https://discord.gg/apex

## ⚡ Quick Commands Reference

```powershell
# Build
.\scripts\build.bat

# Fix Service Worker error
.\scripts\fix-service-worker.bat

# Run node
.\bin\apex.exe

# Generate genesis
.\bin\genesis.exe --output config\genesis.json

# Get help
.\bin\apexctl.exe --help

# Run tests
go test ./...
```

## 🔐 Important

- **Never share** validator private keys
- **Backup** your keys and data
- **Secure** your RPC endpoint
- **Monitor** your node regularly

---

**Version**: 1.0.0  
**License**: MIT  
**Website**: https://apex.network

For the latest updates, see: https://github.com/apex-blockchain/apex
