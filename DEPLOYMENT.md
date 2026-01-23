# Apex Blockchain Production Deployment Guide

## 🎯 Overview

This guide covers deploying Apex Blockchain nodes in production environments.

## 🔧 Infrastructure Requirements

### Hardware Requirements

#### Validator Node (Recommended)
- **CPU**: 8+ cores (AMD/Intel)
- **RAM**: 32 GB
- **Storage**: 1 TB NVMe SSD
- **Network**: 1 Gbps connection
- **Bandwidth**: Unlimited or 10+ TB/month

#### Full Node (Minimum)
- **CPU**: 4 cores
- **RAM**: 16 GB
- **Storage**: 500 GB SSD
- **Network**: 100 Mbps connection
- **Bandwidth**: 5+ TB/month

### Software Requirements

- **OS**: Ubuntu 22.04 LTS or later
- **Go**: 1.21+
- **Docker**: 24.0+ (optional)
- **PostgreSQL**: 15+ (for indexer, optional)

## 🚀 Deployment Methods

### Method 1: Binary Deployment

#### 1. System Preparation

```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y build-essential git curl wget

# Install Go
wget https://go.dev/dl/go1.21.6.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.6.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
```

#### 2. Build Apex

```bash
# Clone repository
git clone https://github.com/apex-blockchain/apex.git
cd apex

# Build binaries
./scripts/build.sh

# Install systemd service
sudo cp apex.service /etc/systemd/system/
sudo systemctl daemon-reload
```

#### 3. Configuration

```bash
# Create data directory
sudo mkdir -p /var/lib/apex
sudo chown $USER:$USER /var/lib/apex

# Copy configuration
cp config/config.yaml /etc/apex/config.yaml

# Edit configuration
nano /etc/apex/config.yaml
```

#### 4. Start Node

```bash
# Start service
sudo systemctl start apex

# Enable auto-start
sudo systemctl enable apex

# Check status
sudo systemctl status apex

# View logs
sudo journalctl -u apex -f
```

### Method 2: Docker Deployment

#### 1. Docker Installation

```bash
# Install Docker
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# Install Docker Compose
sudo apt install docker-compose-plugin
```

#### 2. Prepare Docker Environment

```bash
# Create docker-compose.yml
cat > docker-compose.yml <<EOF
version: '3.8'

services:
  apex:
    image: apexblockchain/apex:latest
    container_name: apex-node
    restart: unless-stopped
    ports:
      - "8545:8545"
      - "30303:30303"
    volumes:
      - ./data:/data
      - ./config:/config
    environment:
      - APEX_CONFIG=/config/config.yaml
    networks:
      - apex-network

networks:
  apex-network:
    driver: bridge
EOF
```

#### 3. Deploy

```bash
# Start container
docker compose up -d

# View logs
docker compose logs -f apex

# Stop container
docker compose down
```

### Method 3: Kubernetes Deployment

See `k8s/` directory for Kubernetes manifests.

## 🔐 Validator Setup

### Generate Validator Keys

```bash
# Generate new keys
./bin/apexctl keys generate \
  --output /secure/path/validator_key.json

# Backup keys securely!
# Store in multiple secure locations
```

### Register Validator

```bash
# Create validator
./bin/apexctl validator create \
  --moniker "My Validator" \
  --commission 10 \
  --self-stake 100000 \
  --key /secure/path/validator_key.json \
  --details "Professional validator" \
  --website "https://myvalidator.com"
```

### Validator Configuration

```yaml
# config/config.yaml
validator:
  enabled: true
  key_file: "/secure/path/validator_key.json"
  rewards_address: "0xYourRewardsAddress"
```

## 🔒 Security Hardening

### Firewall Configuration

```bash
# Allow SSH
sudo ufw allow 22/tcp

# Allow P2P
sudo ufw allow 30303/tcp
sudo ufw allow 30303/udp

# Allow RPC (localhost only recommended)
sudo ufw allow from 127.0.0.1 to any port 8545

# Enable firewall
sudo ufw enable
```

### SSH Hardening

```bash
# Disable password auth
sudo nano /etc/ssh/sshd_config
# Set: PasswordAuthentication no
# Set: PermitRootLogin no

# Restart SSH
sudo systemctl restart sshd
```

### Key Management

- **Never expose private keys**
- **Use hardware wallets** for validator keys
- **Implement key rotation** policies
- **Backup keys** to multiple secure locations
- **Use encrypted storage**

### DDoS Protection

```bash
# Install fail2ban
sudo apt install fail2ban

# Configure for Apex
sudo nano /etc/fail2ban/jail.local
```

## 📊 Monitoring

### Prometheus Setup

```bash
# Install Prometheus
wget https://github.com/prometheus/prometheus/releases/download/v2.45.0/prometheus-2.45.0.linux-amd64.tar.gz
tar xvf prometheus-2.45.0.linux-amd64.tar.gz
cd prometheus-2.45.0.linux-amd64

# Configure prometheus.yml
cat > prometheus.yml <<EOF
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'apex'
    static_configs:
      - targets: ['localhost:9091']
EOF

# Start Prometheus
./prometheus --config.file=prometheus.yml
```

### Grafana Dashboard

```bash
# Install Grafana
sudo apt install -y software-properties-common
sudo add-apt-repository "deb https://packages.grafana.com/oss/deb stable main"
wget -q -O - https://packages.grafana.com/gpg.key | sudo apt-key add -
sudo apt update
sudo apt install grafana

# Start Grafana
sudo systemctl start grafana-server
sudo systemctl enable grafana-server
```

### Key Metrics to Monitor

- Block height
- Peer count
- Transaction throughput
- Mempool size
- CPU/RAM usage
- Disk I/O
- Network bandwidth
- Validator uptime
- Missed blocks

## 🔄 Backup and Recovery

### Automated Backups

```bash
#!/bin/bash
# backup.sh

BACKUP_DIR="/backup/apex"
DATA_DIR="/var/lib/apex"
DATE=$(date +%Y%m%d_%H%M%S)

# Create backup
tar -czf $BACKUP_DIR/apex_backup_$DATE.tar.gz $DATA_DIR

# Keep only last 7 backups
find $BACKUP_DIR -name "apex_backup_*.tar.gz" -mtime +7 -delete
```

### Recovery Procedure

```bash
# Stop node
sudo systemctl stop apex

# Restore from backup
tar -xzf apex_backup_YYYYMMDD_HHMMSS.tar.gz -C /

# Start node
sudo systemctl start apex
```

## 🔄 Upgrade Process

### Zero-Downtime Upgrade

1. **Prepare new version**
```bash
git fetch --tags
git checkout v1.1.0
./scripts/build.sh
```

2. **Test on staging**
```bash
# Run tests
./scripts/test.sh
```

3. **Backup current state**
```bash
./backup.sh
```

4. **Upgrade**
```bash
sudo systemctl stop apex
sudo cp bin/apex /usr/local/bin/
sudo systemctl start apex
```

5. **Verify**
```bash
./bin/apex version
sudo systemctl status apex
```

## 📡 Network Configuration

### Sentry Node Architecture

```
Internet
    |
[Sentry Nodes] ← Public-facing
    |
[Validator Node] ← Private, firewalled
```

### Sentry Node Setup

```yaml
# Sentry node config
network:
  listen_address: "0.0.0.0:30303"
  private_peer_ids: []
  persistent_peers:
    - "/ip4/VALIDATOR_IP/tcp/30303/p2p/VALIDATOR_PEER_ID"
```

### Validator Node Setup

```yaml
# Validator node config
network:
  listen_address: "127.0.0.1:30303"
  private_peer_ids: ["VALIDATOR_PEER_ID"]
  persistent_peers:
    - "/ip4/SENTRY1_IP/tcp/30303/p2p/SENTRY1_PEER_ID"
    - "/ip4/SENTRY2_IP/tcp/30303/p2p/SENTRY2_PEER_ID"
```

## 🚨 Incident Response

### Node is Down

1. Check logs: `sudo journalctl -u apex -n 100`
2. Check disk space: `df -h`
3. Check connectivity: `netstat -tlnp`
4. Restart if needed: `sudo systemctl restart apex`

### Validator Jailed

1. Check why: `./bin/apexctl query validator STATUS`
2. Wait unbonding period
3. Unjail: `./bin/apexctl validator unjail`

### Database Corruption

1. Stop node
2. Restore from backup
3. Sync from peers
4. Restart node

## 📞 Support Channels

- **Emergency**: emergency@apex.network
- **Discord**: #validator-support
- **Documentation**: https://docs.apex.network
- **Status Page**: https://status.apex.network

## ✅ Deployment Checklist

- [ ] Hardware meets requirements
- [ ] OS updated and hardened
- [ ] Firewall configured
- [ ] SSH secured
- [ ] Node built and tested
- [ ] Configuration reviewed
- [ ] Validator keys secured
- [ ] Backups automated
- [ ] Monitoring setup
- [ ] Alerts configured
- [ ] Documentation reviewed
- [ ] Team trained

---

**Need help?** Contact the Apex team or join our Discord!