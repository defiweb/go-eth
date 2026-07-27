package kzg4844

import (
	"crypto/sha256"
	"encoding/hex"
	"math/big"
	"testing"

	"github.com/defiweb/go-eth/crypto/primitives"
	"github.com/defiweb/go-eth/hexutil"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func getRandBlob(seed int64) *primitives.KZGBlob {
	blob := primitives.KZGBlob{}
	for i := 0; i < primitives.KZGBlobSize; i += primitives.KZGScalarSize {
		h := sha256.Sum256([]byte{byte(seed + int64(i))})
		p := new(big.Int).SetBytes(h[:])
		p = p.Mod(p, BLSModulus)
		p.FillBytes(blob[i : i+primitives.KZGScalarSize])
	}
	return &blob
}

func getPoint(blob *primitives.KZGBlob) (p primitives.KZGPoint) {
	h := sha256.Sum256(blob[:])
	i := new(big.Int).SetBytes(h[:])
	i = i.Mod(i, BLSModulus)
	i.FillBytes(p[:])
	return p
}

func TestVerifyProof(t *testing.T) {
	blob := getRandBlob(0)
	point := getPoint(blob)

	// Get commitment for the blob.
	commitment, err := BlobToCommitment(blob)
	require.NoError(t, err)

	// Compute proof and claim.
	proof, claim, err := ComputeProof(blob, point)
	require.NoError(t, err)

	// Verify the proof.
	err = VerifyProof(commitment, point, claim, proof)
	assert.NoError(t, err)
}

func TestBlobVerifyProof(t *testing.T) {
	// Create a test blob
	blob := getRandBlob(0)

	// Get commitment for the blob.
	commitment, err := BlobToCommitment(blob)
	require.NoError(t, err)

	// Compute proof.
	proof, err := ComputeBlobProof(blob, commitment)
	require.NoError(t, err)

	// Verify the proof.
	err = VerifyBlobProof(blob, commitment, proof)
	assert.NoError(t, err)
}

func TestComputeBlobHashV1(t *testing.T) {
	tc := []struct {
		name       string
		commitment string
		wantHash   string
	}{
		{
			name:       "zero commitment",
			commitment: "000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000",
			wantHash:   "01b0761f87b081d5cf10757ccc89f12be355c70e2e29df288b65b30710dcbcd1",
		},
		{
			name:       "test commitment",
			commitment: "123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef0",
			wantHash:   "012a8194ef18215aa6278923a1b143e6cf3a28087a3d26bcf6f887befc62cb47",
		},
	}
	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			var commitment primitives.KZGCommitment
			copy(commitment[:], hexutil.MustHexToBytes(tt.commitment))

			hash := ComputeBlobHashV1(commitment)
			assert.Equal(t, tt.wantHash, hex.EncodeToString(hash[:]))
		})
	}
}
