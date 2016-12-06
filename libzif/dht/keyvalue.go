package dht

type KeyValue struct {
	Key   Address
	Value []byte // Max size of 64kbs

	// Used for sorting, compare on key-value result to another.
	distance Address
}

func NewKeyValue(key Address, value []byte) *KeyValue {
	ret := &KeyValue{}

	ret.Key = key
	copy(ret.Value, value)

	return ret
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
