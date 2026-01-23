# Apex Blockchain

<div align="center">

![Apex Logo](https://via.placeholder.com/200x200/4f46e5/ffffff?text=APEX)

**A High-Performance Layer-1 Blockchain with Delegated Proof of Stake**

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Version](https://img.shields.io/badge/Go-1.21+-blue.svg)](https://golang.org)
[![Build Status](https://img.shields.io/badge/build-passing-brightgreen.svg)]()

[Documentation](https://docs.apex.network) | [Explorer](https://explorer.apex.network) | [Website](https://apex.network)

</div>

## 🌟 Overview

**Apex Blockchain** is a next-generation Layer-1 blockchain built with Go, featuring:

- **Native Token**: APX with 500M total supply
- **Consensus**: Delegated Proof of Stake (DPoS) with 21 validators
- **Block Time**: 3 seconds for fast finality
- **Smart Contracts**: EVM-compatible execution environment
- **Staking**: Flexible delegation and reward system
- **Governance**: On-chain voting and proposals

## 🚀 Key Features

### Delegated Proof of Stake (DPoS)
- **21 Active Validators**: Selected by voting power
- **Fast Finality**: 3-second block time
- **Democratic**: Token holders vote for validators
- **Rewards**: Block producers earn APX rewards
- **Slashing**: Penalties for misbehavior

### Staking & Delegation
- **Minimum Stake**: 10 APX for delegation
- **Validator Stake**: 100K APX minimum
- **Flexible**: Delegate to multiple validators
- **Rewards**: Earn staking rewards proportionally
- **Unbonding**: 7-day unbonding period

### Security Features
- **Slashing Mechanisms**:
  - 5% slash for double signing
  - 1% slash for prolonged downtime
  - 3% slash for invalid blocks
- **Jail System**: Temporary validator suspension
- **Validator Monitoring**: Real-time uptime tracking

## 📋 Prerequisites

- **Go**: 1.21 or higher
- **Git**: For cloning the repository
- **Make**: For building (optional)

## 🛠️ Installation

### From Source

```bash
# Clone the repository
git clone https://github.com/apex-blockchain/apex.git
cd apex

# Install dependencies
go mod download

# Build the node
go build -o bin/apex ./cmd/apex

# Build additional tools
go build -o bin/apexctl ./cmd/apexctl
go build -o bin/genesis ./cmd/genesis
```

### Using Docker

```bash
# Build Docker image
docker build -t apex-blockchain .

# Run node
docker run -d -p 8545:8545 -p 30303:30303 \
  -v /path/to/data:/data \
  apex-blockchain
```

## 🎯 Quick Start

### 1. Initialize Configuration

```bash
# Copy default config
cp config/config.yaml.example config/config.yaml

# Edit configuration as needed
nano config/config.yaml
```

### 2. Run a Full Node

```bash
# Start the node
./bin/apex --config config/config.yaml
```

### 3. Run a Validator Node

```bash
# Generate validator keys
./bin/apexctl keys generate --output validator_key.json

# Configure validator in config.yaml
validator:
  enabled: true
  key_file: "./validator_key.json"

# Start validator node
./bin/apex --config config/config.yaml
```

## 📡 JSON-RPC API

The Apex node exposes a JSON-RPC API on port 8545 (configurable).

### Available Methods

#### Blockchain Methods
- `apex_blockNumber` - Get current block height
- `apex_getBlockByNumber` - Get block by number
- `apex_getBlockByHash` - Get block by hash
- `apex_getTransaction` - Get transaction details

#### Account Methods
- `apex_getBalance` - Get account balance
- `apex_getNonce` - Get account nonce
- `apex_sendTransaction` - Submit transaction

#### Staking Methods
- `apex_stake` - Stake tokens to validator
- `apex_unstake` - Initiate unstaking
- `apex_getStakingInfo` - Get staking information
- `apex_getValidators` - List active validators

### Example Usage

```bash
# Get current block number
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "apex_blockNumber",
    "params": [],
    "id": 1
  }'

# Get account balance
curl -X POST http://localhost:8545 \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "apex_getBalance",
    "params": ["0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0"],
    "id": 1
  }'
```

## 🔐 Staking Guide

### Delegate to a Validator

```bash
# Using apexctl
./bin/apexctl stake \
  --validator 0x742d35Cc6634C0532925a3b844Bc9e7595f0bEb0 \
  --amount 1000 \
  --from your_wallet_address
```

### Become a Validator

```bash
# Create validator
./bin/apexctl validator create \
  --moniker "My Validator" \
  --commission 10 \
  --self-stake 100000 \
  --key validator_key.json
```

### Check Rewards

```bash
# Query staking rewards
./bin/apexctl query staking rewards \
  --delegator your_address \
  --validator validator_address
```

## 🏗️ Architecture

```
apex/
├── cmd/               # Command-line tools
│   ├── apex/          # Main node executable
│   ├── apexctl/       # CLI management tool
│   └── genesis/       # Genesis block generator
├── pkg/               # Core packages
│   ├── api/           # API layer (JSON-RPC, gRPC)
│   ├── consensus/     # DPoS consensus engine
│   ├── core/          # Blockchain core logic
│   ├── crypto/        # Cryptographic functions
│   ├── mempool/       # Transaction pool
│   ├── network/       # P2P networking
│   ├── staking/       # Staking mechanisms
│   ├── storage/       # Database layer
│   └── types/         # Common types
├── config/            # Configuration files
└── scripts/           # Build and deployment scripts
```

## 🔧 Configuration

### Node Configuration (`config/config.yaml`)

```yaml
node:
  name: "apex-node-1"
  chain_id: "apex-mainnet-1"

rpc:
  enabled: true
  port: 8545

consensus:
  type: "dpos"
  block_time: 3
  max_validators: 21

validator:
  enabled: false
  key_file: "./validator_key.json"
```

### Genesis Configuration (`config/genesis.json`)

Configure initial validators, accounts, and consensus parameters.

## 📊 Tokenomics

- **Total Supply**: 500,000,000 APX
- **Decimals**: 18
- **Initial Distribution**:
  - Foundation Treasury: 20% (100M APX)
  - Community Fund: 10% (50M APX)
  - Development Fund: 6% (30M APX)
  - Marketing: 4% (20M APX)
  - Validators & Staking: 40% (200M APX)
  - Public Sale: 20% (100M APX)

### Block Rewards

- **Initial**: 2 APX per block
- **Halving**: Every 4 years
- **Minimum**: 0.1 APX per block
- **Estimated APY**: 8-15% (varies with staking participation)

## 🧪 Testing

```bash
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific package tests
go test ./pkg/consensus/...

# Run integration tests
go test -tags=integration ./tests/...
```

## 🔒 Security

### Audit Status
- [ ] Phase 1: Code Review (Planned Q2 2026)
- [ ] Phase 2: Security Audit (Planned Q3 2026)
- [ ] Phase 3: Penetration Testing (Planned Q4 2026)

### Security Best Practices
1. Always use hardware wallets for validator keys
2. Enable firewall and rate limiting
3. Keep your node software updated
4. Use TLS for RPC endpoints
5. Implement DDoS protection

### Bug Bounty
We run a bug bounty program. Report vulnerabilities to security@apex.network.

## 🤝 Contributing

We welcome contributions! Please see [CONTRIBUTING.md](CONTRIBUTING.md) for details.

### Development Workflow

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📝 License

This project is licensed under the MIT License - see [LICENSE](LICENSE) file for details.

## 🌐 Community

- **Website**: https://apex.network
- **Documentation**: https://docs.apex.network
- **Explorer**: https://explorer.apex.network
- **Discord**: https://discord.gg/apex
- **Twitter**: https://twitter.com/apexblockchain
- **Telegram**: https://t.me/apexblockchain

## 🗺️ Roadmap

### Phase 1: Foundation (Q1 2026) ✅
- [x] Core blockchain implementation
- [x] DPoS consensus mechanism
- [x] Basic staking functionality

### Phase 2: Enhancement (Q2 2026)
- [ ] Smart contract support (EVM)
- [ ] Cross-chain bridges
- [ ] Enhanced governance

### Phase 3: Scaling (Q3 2026)
- [ ] Layer-2 solutions
- [ ] Sharding implementation
- [ ] Advanced privacy features

### Phase 4: Ecosystem (Q4 2026)
- [ ] DEX integration
- [ ] NFT marketplace
- [ ] Developer grants program

## 📞 Support

- **Email**: support@apex.network
- **Discord**: #support channel
- **Documentation**: https://docs.apex.network

## 🙏 Acknowledgments

Built with ❤️ using:
- [Go](https://golang.org)
- [libp2p](https://libp2p.io)
- [BadgerDB](https://github.com/dgraph-io/badger)
- [Cobra](https://github.com/spf13/cobra)

---

<div align="center">

Made with ❤️ by the Apex Team

**[Website](https://apex.network)** | **[Docs](https://docs.apex.network)** | **[Explorer](https://explorer.apex.network)**

</div>