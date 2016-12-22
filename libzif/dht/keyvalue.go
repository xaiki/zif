package dht

const (
	MaxValueSize = 10 * 1024
)

type KeyValue struct {
	Key   Address
	Value []byte // Max size of 64kbs

	distance Address
}

func NewKeyValue(key Address, value []byte) *KeyValue {
	ret := &KeyValue{}

	ret.Key = key
	ret.Value = make([]byte, len(value))
	copy(ret.Value, value)

	return ret
}

func (kv *KeyValue) Valid() bool {
	return len(kv.Value) <= MaxValueSize && len(kv.Key.Raw) == AddressBinarySize
}

type Pairs []*KeyValue

func (e Pairs) Len() int {
	return len(e)
}

func (e Pairs) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Pairs) Less(i, j int) bool {
	return e[i].distance.Less(&e[j].distance)
}
