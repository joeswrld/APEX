package crypto

import (
	"crypto/ecdsa"
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"math/big"

	"github.com/apex/pkg/types"
)

// SignData signs data with private key
func SignData(data []byte, privKey *ecdsa.PrivateKey) (types.Signature, error) {
	hash := sha256.Sum256(data)
	
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return nil, err
	}
	
	// Serialize signature
	signature := append(r.Bytes(), s.Bytes()...)
	return types.Signature(signature), nil
}

// VerifySignature verifies signature with public key
func VerifySignature(data []byte, signature types.Signature, pubKey *ecdsa.PublicKey) bool {
	if len(signature) < 64 {
		return false
	}
	
	hash := sha256.Sum256(data)
	
	// Deserialize signature
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	
	return ecdsa.Verify(pubKey, hash[:], r, s)
}

// RecoverPublicKey recovers public key from signature (simplified)
func RecoverPublicKey(data []byte, signature types.Signature) (*ecdsa.PublicKey, error) {
	// This is a simplified implementation
	// In production, use proper ECDSA recovery like in Ethereum
	return nil, errors.New("not implemented")
}

// SignHash signs a hash directly
func SignHash(hash types.Hash, privKey *ecdsa.PrivateKey) (types.Signature, error) {
	r, s, err := ecdsa.Sign(rand.Reader, privKey, hash[:])
	if err != nil {
		return nil, err
	}
	
	signature := append(r.Bytes(), s.Bytes()...)
	return types.Signature(signature), nil
}

// VerifyHashSignature verifies hash signature
func VerifyHashSignature(hash types.Hash, signature types.Signature, pubKey *ecdsa.PublicKey) bool {
	if len(signature) < 64 {
		return false
	}
	
	r := new(big.Int).SetBytes(signature[:32])
	s := new(big.Int).SetBytes(signature[32:64])
	
	return ecdsa.Verify(pubKey, hash[:], r, s)
}
