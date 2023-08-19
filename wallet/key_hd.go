package wallet

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg"

	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/tyler-smith/go-bip39"
)

// The code below is based on:
// github.com/miguelmota/go-ethereum-hdwallet
// github.com/ethereum/go-ethereum/blob/master/accounts

// Mnemonic is a mnemonic phrase with a password used to derive private keys.
type Mnemonic struct {
	// Internally, we only need a master key to derive private keys, but to
	// keep package API easier to understand, the whole structure is called
	// mnemonic because it is the name that users are familiar with.
	masterKey *hdkeychain.ExtendedKey
}

// DerivationPath represents derivation path as internal binary format.
//
// Derivation path allows to derive multiple child keys from a single parent
// key.
//
// The path is defined as:
//
//	m / purpose' / coin_type' / account' / change / address_index
//
// Where:
//
//   - purpose is a constant set to 44' (or 0x8000002C) following the BIP43
//     recommendation.
//
//   - coin_type indicates the coin type, as defined in SLIP-0044. For Ethereum,
//     it is 60' (or 0x8000003C).
//
//   - account is an index that allows users to create multiple identities from a
//     single seed.
//
//   - change is a constant, set to 0 (or 0x80000000) for external chain and 1
//     (or 0x80000001) for internal chain.
//
//   - address_index is an address index that is incremented for each new
//     address.
//
// Reference:
// BIP-32: https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki
// BIP-44: https://github.com/bitcoin/bips/blob/master/bip-0044.mediawiki
// SLIP-44: https://github.com/satoshilabs/slips/blob/master/slip-0044.md
type DerivationPath []uint32

// RootDerivationPath is the root derivation path used to derive the child keys.
// It is set to m/44'/60'/0'/0.
var RootDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0}

// DefaultDerivationPath is the default derivation path for the first key.
// It is set to m/44'/60'/0'/0/0.
var DefaultDerivationPath = DerivationPath{0x80000000 + 44, 0x80000000 + 60, 0x80000000 + 0, 0, 0}

// Indices of the components in the derivation path.
const (
	PurposeComponent      = 0
	CoinTypeComponent     = 1
	AccountComponent      = 2
	ChangeComponent       = 3
	AddressIndexComponent = 4
)

// NewKeyFromMnemonic creates a new private key from a mnemonic phrase.
// The derivation path is set to m/44'/60'/account'/0/index.
func NewKeyFromMnemonic(mnemonic, password string, account, index uint32) (*PrivateKey, error) {
	m, err := NewMnemonic(mnemonic, password)
	if err != nil {
		return nil, err
	}
	dp := make(DerivationPath, len(DefaultDerivationPath))
	copy(dp, DefaultDerivationPath)
	_ = dp.SetAccount(account)
	_ = dp.SetAddressIndex(index)
	return m.Derive(dp)
}

// NewMnemonic creates a new mnemonic that can be used to derive private keys.
func NewMnemonic(mnemonic, password string) (Mnemonic, error) {
	if mnemonic == "" {
		return Mnemonic{}, errors.New("mnemonic is required")
	}
	seed, err := bip39.NewSeedWithErrorChecking(mnemonic, password)
	if err != nil {
		return Mnemonic{}, err
	}
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return Mnemonic{}, err
	}
	return Mnemonic{masterKey: masterKey}, nil
}

// Derive derives a private key from the mnemonic using given derivation path.
func (m Mnemonic) Derive(path DerivationPath) (*PrivateKey, error) {
	var err error
	key := m.masterKey
	for _, n := range path {
		key, err = key.Derive(n)
		if err != nil {
			return nil, err
		}
	}
	privKey, err := key.ECPrivKey()
	privKeyECDSA := privKey.ToECDSA()
	if err != nil {
		return nil, err
	}
	return NewKeyFromECDSA(privKeyECDSA), nil
}

// ParseDerivationPath converts a BIP-33 derivation path string into the
// internal binary format.
//
// The path is expected to be of the form:
//
//	m / purpose' / coin_type' / account' / change / address_index
//
// The single quotes are used to indicate hardened derivation.
//
// If m/ prefix is omitted, then RootDerivationPath (m/44'/60'/0'/0) is
// prepended to the path.
func ParseDerivationPath(path string) (result DerivationPath, err error) {
	const (
		stateStart       = iota // Start of the path, before the first non-whitespace character.
		stateEnd                // End of the path.
		stateAfterSlash         // After a slash, before the next non-whitespace character.
		stateBeforeSlash        // Before a slash, expecting a whitespace character or the end of the path.
		stateNumber             // Inside a number.
		stateAfterNumber        // After a number, expecting a whitespaces, slash or hardening character.
	)
	var (
		pos      = 0          // Current position in the string.
		numStart = 0          // Start position of the current number.
		numEnd   = 0          // End position of the current number.
		relative = true       // Whether the path is relative or absolute.
		state    = stateStart // Current state of the parser.
	)
	for {
		var char byte
		if pos < len(path) {
			char = path[pos]
		}
		switch state {
		case stateStart:
			switch {
			case isWhitespace(char):
				// Ignore whitespace.
			case char == 'm':
				state = stateBeforeSlash
				relative = false
			case isNibble(char):
				state = stateNumber
				numStart = pos
			default:
				return nil, fmt.Errorf("invalid character %q at position %d", char, pos)
			}
		case stateBeforeSlash:
			switch {
			case isWhitespace(char):
				// Ignore whitespace.
			case char == 0 && pos == len(path):
				state = stateEnd
			case char == '/':
				state = stateAfterSlash
			default:
				return nil, fmt.Errorf("invalid character %q at position %d", char, pos)
			}
		case stateAfterSlash:
			switch {
			case isWhitespace(char):
				// Ignore whitespace.
			case isNibble(char):
				state = stateNumber
				numStart = pos
			default:
				return nil, fmt.Errorf("invalid character %q at position %d", char, pos)
			}
		case stateNumber:
			switch {
			case isNibble(char):
				// Continue reading the number.
			default:
				state = stateAfterNumber
				numEnd = pos
				pos-- // Re-read the current character in the next state.
			}
		case stateAfterNumber:
			switch {
			case isWhitespace(char):
			// Ignore whitespace.
			case char == 0 && pos == len(path):
				state = stateEnd
				result = append(result, parseNumber(path[numStart:numEnd]))
			case char == '/' || (char == 0 && pos == len(path)):
				state = stateAfterSlash
				result = append(result, parseNumber(path[numStart:numEnd]))
			case char == '\'':
				// Hardened derivation.
				state = stateBeforeSlash
				number := parseNumber(path[numStart:numEnd])
				if number&0x80000000 != 0 {
					return nil, fmt.Errorf("component overflows int32")
				}
				result = append(result, parseNumber(path[numStart:numEnd])|0x80000000)
			default:
				return nil, fmt.Errorf("invalid character %q at position %d", char, pos)
			}
		}
		if state == stateEnd {
			break
		}
		pos++
	}
	if len(result) == 0 {
		return nil, errors.New("derivation path is empty")
	}
	if relative {
		result = append(RootDerivationPath, result...)
	}
	return result, nil
}

// Purpose returns the purpose component of the derivation path.
// If the path is empty, then 0 is returned.
func (dp DerivationPath) Purpose() uint32 {
	if len(dp) == 0 {
		return 0
	}
	return dp[PurposeComponent] & 0x7fffffff
}

// CoinType returns the coin type component of the derivation path.
// If coin type component is missing, then 0 is returned.
func (dp DerivationPath) CoinType() uint32 {
	if len(dp) <= CoinTypeComponent {
		return 0
	}
	return dp[CoinTypeComponent] & 0x7fffffff
}

// Account returns the account component of the derivation path.
// If account component is missing, then 0 is returned.
func (dp DerivationPath) Account() uint32 {
	if len(dp) <= AccountComponent {
		return 0
	}
	return dp[AccountComponent] & 0x7fffffff
}

// Change returns the change component of the derivation path.
// If change component is missing, then 0 is returned.
func (dp DerivationPath) Change() uint32 {
	if len(dp) <= ChangeComponent {
		return 0
	}
	return dp[ChangeComponent]
}

// AddressIndex returns the address index component of the derivation path.
// If address index component is missing, then 0 is returned.
func (dp DerivationPath) AddressIndex() uint32 {
	if len(dp) <= AddressIndexComponent {
		return 0
	}
	return dp[AddressIndexComponent]
}

// SetAccount sets the account component of the derivation path.
func (dp DerivationPath) SetAccount(account uint32) error {
	if len(dp) <= AccountComponent {
		return errors.New("derivation path has no account component")
	}
	if account&0x80000000 != 0 {
		return errors.New("account number is too large")
	}
	dp[AccountComponent] = account | 0x80000000
	return nil
}

// SetChange sets the change component of the derivation path.
func (dp DerivationPath) SetChange(change uint32) error {
	if len(dp) <= ChangeComponent {
		return errors.New("derivation path has no change component")
	}
	dp[ChangeComponent] = change
	return nil
}

// SetAddressIndex sets the address index component of the derivation path.
func (dp DerivationPath) SetAddressIndex(index uint32) error {
	if len(dp) <= AddressIndexComponent {
		return errors.New("derivation path has no address index component")
	}
	dp[AddressIndexComponent] = index
	return nil
}

// IncreaseAccount increases the account number in the derivation path.
func (dp DerivationPath) IncreaseAccount() error {
	if len(dp) <= AccountComponent {
		return errors.New("derivation path has no account component")
	}
	dp[AccountComponent]++
	return nil
}

// IncreaseAddressIndex increases the index number in the derivation path.
func (dp DerivationPath) IncreaseAddressIndex() error {
	if len(dp) <= AddressIndexComponent {
		return errors.New("derivation path has no address index component")
	}
	dp[AddressIndexComponent]++
	return nil
}

// Increase increases the last component of the derivation path.
func (dp DerivationPath) Increase() error {
	if len(dp) == 0 {
		return errors.New("derivation path is empty")
	}
	dp[len(dp)-1]++
	return nil
}

func isWhitespace(c byte) bool {
	return c == ' ' || c == '\t' || c == '\r' || c == '\n'
}

func isNibble(c byte) bool {
	return c >= '0' && c <= '9' || c >= 'a' && c <= 'f' || c >= 'A' && c <= 'F' || c == 'x'
}

func parseNumber(s string) uint32 {
	result, _ := strconv.ParseUint(s, 0, 32)
	return uint32(result)
}
