package data

type Encodable interface {
	Bytes() ([]byte, error)
	String() (string, error)

	// The latter two may be equivelant
	Json() ([]byte, error)
	JsonString() ([]byte, error)
}

type Signer interface {
	Sign([]byte) []byte
	PublicKey() []byte
}
