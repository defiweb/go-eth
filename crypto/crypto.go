package crypto

import (
	"fmt"
)

// AddMessagePrefix adds the Ethereum message prefix to the given data.
func AddMessagePrefix(data []byte) []byte {
	return []byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d%s", len(data), data))
}
