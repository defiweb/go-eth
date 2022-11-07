package crypto

import (
	"fmt"

	"github.com/btcsuite/btcd/btcec"

	"web3rpc/types"
)

func Ecrecover(hash types.Hash, sig types.Signature) (types.Address, error) {
	v := sig[types.SignatureLength-1]
	copy(sig[1:], sig[:types.SignatureLength-1])
	sig[0] = v
	if v < 27 {
		sig[0] += 27
	}
	pub, _, err := btcec.RecoverCompact(s256, sig.Bytes(), hash.Bytes())
	if err != nil {
		return types.Address{}, err
	}
	return pubKeyToAddress(pub.ToECDSA()), nil
}

func EcrecoverMessage(data []byte, sig types.Signature) (types.Address, error) {
	return Ecrecover(Keccak256(formatMessage(data)), sig)
}

func EcrecoverTransaction(tx *types.Transaction) (types.Address, error) {
	d, err := tx.SigningData()
	if err != nil {
		return types.Address{}, err
	}
	addr, err := Ecrecover(Keccak256(d), tx.Signature)
	if err != nil {
		return types.Address{}, err
	}
	return addr, nil
}

func formatMessage(data []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
}
