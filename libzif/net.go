// a few network helpers

package libzif

import "golang.org/x/crypto/ed25519"

type ConnHeader struct {
	cl Client
	pk ed25519.PublicKey
}
