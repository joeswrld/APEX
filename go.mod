module github.com/apex-blockchain/apex

go 1.21

require (
	github.com/dgraph-io/badger/v4 v4.2.0
	github.com/libp2p/go-libp2p v0.32.2
	github.com/libp2p/go-libp2p-kad-dht v0.25.2
	github.com/libp2p/go-libp2p-pubsub v0.10.0
	github.com/multiformats/go-multiaddr v0.12.0
	github.com/spf13/cobra v1.8.0
	github.com/spf13/viper v1.18.2
	google.golang.org/grpc v1.60.1
	google.golang.org/protobuf v1.32.0
	gopkg.in/yaml.v3 v3.0.1
	github.com/btcsuite/btcd/btcec/v2 v2.3.2
	github.com/ethereum/go-ethereum v1.13.8
	github.com/golang/protobuf v1.5.3
	github.com/google/uuid v1.5.0
	github.com/stretchr/testify v1.8.4
	go.uber.org/zap v1.26.0
	golang.org/x/crypto v0.18.0
)
EOF
echo "go.mod created"