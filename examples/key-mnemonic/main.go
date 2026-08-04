package main

import (
	"fmt"

	"github.com/defiweb/go-eth/wallet"
)

func main() {
	// Parse mnemonic.
	mnemonic, err := wallet.NewMnemonic("gravity trophy shrimp suspect sheriff avocado label trust dove tragic pitch title network myself spell task protect smooth sword diary brain blossom under bulb", "")
	if err != nil {
		panic(err)
	}

	// Parse derivation path.
	path, err := wallet.ParseDerivationPath("m/44'/60'/0'/10/10")
	if err != nil {
		panic(err)
	}

	// Derive private key.
	key, err := mnemonic.Derive(path)
	if err != nil {
		panic(err)
	}

	// Print the address of the derived private key.
	fmt.Println("Private key:", key.Address().String())
}
