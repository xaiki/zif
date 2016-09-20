package main

import "golang.org/x/crypto/ed25519"
import "errors"

const ProtocolHeaderSize = 2 + 2 + ed25519.PublicKeySize

// Non-exported members are NOT sent over the wire.
// For instance, the zif address can be generated from the public key quite
// easily.
type ProtocolHeader struct {
	// This is Zif.
	Zif [2]byte

	// Protocol versions, ignores peers where this differs.
	Version [2]byte

	// Address from this, also used for verficication of other things.
	PublicKey [ed25519.PublicKeySize]byte

	zifAddress Address
}

func ProtocolHeaderFromBytes(raw []byte) (ProtocolHeader, error) {
	var ph ProtocolHeader

	if len(raw) < ProtocolHeaderSize {
		return ph, errors.New("Incorrect header size")
	}

	copy(ph.Zif[:], raw[:2])
	copy(ph.Version[:], raw[2:4])
	copy(ph.PublicKey[:], raw[4:4+ed25519.PublicKeySize])

	ph.zifAddress.Generate(ph.PublicKey[:])

	return ph, nil
}

func (ph *ProtocolHeader) Bytes() []byte {
	ret := make([]byte, 0, ProtocolHeaderSize)

	ret = append(ret, ph.Zif[:]...)
	ret = append(ret, ph.Version[:]...)
	ret = append(ret, ph.PublicKey[:]...)

	return ret
}
