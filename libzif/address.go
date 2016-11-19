// Generates a zif address given a public key
// Similar to the method Bitcoin uses
// see: https://en.bitcoin.it/wiki/Technical_background_of_version_1_Bitcoin_addresses

package zif

import (
	"bytes"
	"errors"

	"github.com/prettymuchbryce/hellobitcoin/base58check"
	"golang.org/x/crypto/ripemd160"
	"golang.org/x/crypto/sha3"
)

const AddressBinarySize = 20

type Address struct {
	Bytes []byte
}

// Generates an Address from a PublicKey.
func NewAddress(key []byte) (addr Address) {
	addr = Address{}
	addr.Generate(key)

	return
}

// Returns Address.Bytes Base58 encoded and prepended with a Z.
// Base58 removes ambiguous characters, reducing the chances of address confusion.
func (a *Address) Encode() string {
	return base58check.Encode("51", a.Bytes)
}

// Decodes a string address into address bytes.
func DecodeAddress(value string) Address {
	var addr Address
	addr.Bytes = base58check.Decode(value)

	return addr
}

// Generate a Zif address from a public key.
// This process involves one SHA3-256 iteration, followed by RIPEMD160. This is
// similar to bitcoin, and the RIPEMD160 makes the resulting address a bit shorted.
func (a *Address) Generate(key []byte) string {
	ripemd := ripemd160.New()

	if len(key) != 32 {
		panic(errors.New("Local peer public key is not 32 bytes"))
	}

	// Why hash and not just use the pub key?
	// This way we can change curve or algorithm entirely, and still have
	// the same format for addresses.

	firstHash := sha3.Sum256(key)
	ripemd.Write(firstHash[:])

	secondHash := ripemd.Sum(nil)

	a.Bytes = secondHash

	return a.Encode()
}

func (a *Address) Less(other *Address) bool {

	for i := 0; i < len(a.Bytes); i++ {
		if a.Bytes[i] != other.Bytes[i] {
			return a.Bytes[i] < other.Bytes[i]
		}
	}

	return false
}

func (a *Address) Xor(other *Address) *Address {
	var ret Address
	ret.Bytes = make([]byte, len(a.Bytes))

	for i := 0; i < len(a.Bytes); i++ {
		ret.Bytes[i] = a.Bytes[i] ^ other.Bytes[i]
	}

	return &ret
}

// Counts the number of leading zeroes this address has.
// The address should be the result of an Xor.
// This shows the k-bucket that this address should go into.
func (a *Address) LeadingZeroes() int {

	for i := 0; i < len(a.Bytes); i++ {
		for j := 0; j < 8; j++ {
			if (a.Bytes[i]>>uint8(7-j))&0x1 != 0 {
				return i*8 + j
			}
		}
	}

	return len(a.Bytes)*8 - 1
}

func (a *Address) Equals(other *Address) bool {
	return bytes.Equal(a.Bytes, other.Bytes)
}
