package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"

	"github.com/apex/pkg/types"
	"github.com/btcsuite/btcd/btcec/v2"
	"golang.org/x/crypto/ripemd160"
)

// GenerateKeyPair generates a new ECDSA key pair
func GenerateKeyPair() (*ecdsa.PrivateKey, *ecdsa.PublicKey, error) {
	privKey, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, nil, err
	}
	return privKey, &privKey.PublicKey, nil
}

// PublicKeyToAddress converts public key to address
func PublicKeyToAddress(pubKey *ecdsa.PublicKey) types.Address {
	// Serialize public key
	pubBytes := elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
	
	// Hash with SHA256
	sha := sha256.Sum256(pubBytes)
	
	// Hash with RIPEMD160
	ripemd := ripemd160.New()
	ripemd.Write(sha[:])
	hash := ripemd.Sum(nil)
	
	var addr types.Address
	copy(addr[:], hash[:20])
	return addr
}

// PrivateKeyToHex exports private key to hex
func PrivateKeyToHex(privKey *ecdsa.PrivateKey) string {
	return hex.EncodeToString(privKey.D.Bytes())
}

// HexToPrivateKey imports private key from hex
func HexToPrivateKey(hexKey string) (*ecdsa.PrivateKey, error) {
	bytes, err := hex.DecodeString(hexKey)
	if err != nil {
		return nil, err
	}
	
	privKey := new(ecdsa.PrivateKey)
	privKey.PublicKey.Curve = elliptic.P256()
	privKey.D = new(btcec.ModNScalar)
	privKey.D.SetByteSlice(bytes)
	
	privKey.PublicKey.X, privKey.PublicKey.Y = privKey.PublicKey.Curve.ScalarBaseMult(bytes)
	
	return privKey, nil
}

// PublicKeyToBytes serializes public key
func PublicKeyToBytes(pubKey *ecdsa.PublicKey) []byte {
	return elliptic.Marshal(pubKey.Curve, pubKey.X, pubKey.Y)
}

// BytesToPublicKey deserializes public key
func BytesToPublicKey(bytes []byte) (*ecdsa.PublicKey, error) {
	curve := elliptic.P256()
	x, y := elliptic.Unmarshal(curve, bytes)
	if x == nil {
		return nil, errors.New("invalid public key")
	}
	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}

// HashData hashes data with SHA256
func HashData(data []byte) types.Hash {
	hash := sha256.Sum256(data)
	var h types.Hash
	copy(h[:], hash[:])
	return h
}

// ValidateAddress validates an address format
func ValidateAddress(addr types.Address) bool {
	// Basic validation - check if not all zeros
	zero := types.Address{}
	return addr != zero
}
