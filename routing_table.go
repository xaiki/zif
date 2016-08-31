// Kademlia

package main

import (
	"container/list"
	"sort"
)

const BucketSize = 20

type Entry struct {
	ZifAddress    Address
	Name          string
	Desc          string
	PublicAddress string
	PublicKey     []byte
	Port          int

	// Used in the FindClosest function, for sorting.
	distance Address
}

type Entries []*Entry

func (e Entries) Len() int {
	return len(e)
}

func (e Entries) Swap(i, j int) {
	e[i], e[j] = e[j], e[i]
}

func (e Entries) Less(i, j int) bool {
	return e[i].distance.Less(&e[j].distance)
}

type RoutingTable struct {
	LocalAddress Address
	Buckets      []*list.List
}

func (rt *RoutingTable) Setup(addr Address) {
	rt.LocalAddress = addr
	rt.Buckets = make([]*list.List, len(rt.LocalAddress.Bytes)*8)

	for i := 0; i < len(rt.LocalAddress.Bytes)*8; i++ {
		rt.Buckets[i] = list.New()
	}
}

func (rt *RoutingTable) Update(entry Entry) {
	zero_count := entry.ZifAddress.Xor(&rt.LocalAddress).LeadingZeroes()
	bucket := rt.Buckets[zero_count]

	// TODO: Ping peers, starting from back. If none reply, remove them.
	// Ensures only active peers are stored.
	if bucket.Len() == BucketSize {
		return
	}

	var foundEntry *list.Element = nil
	for i := bucket.Front(); i != nil; i = i.Next() {
		if i.Value.(*Entry).ZifAddress.Equals(&entry.ZifAddress) {
			foundEntry = i
		}
	}

	if foundEntry == nil {
		bucket.PushFront(&entry)
	} else {
		bucket.MoveToFront(foundEntry)
	}
}

func copyToEntrySlice(slice *[]*Entry, begin *list.Element, count int) {

	for i := begin; i != nil && len(*slice) < count; i = i.Next() {
		*slice = append(*slice, i.Value.(*Entry))
	}

}

func (rt *RoutingTable) FindClosest(target Address, count int) []*Entry {
	ret := make([]*Entry, 0, count)

	bucket_num := target.Xor(&rt.LocalAddress).LeadingZeroes()
	bucket := rt.Buckets[bucket_num]

	copyToEntrySlice(&ret, bucket.Front(), count)

	// If the bucket is not filled, look the the buckets either side.
	for i := 1; (bucket_num-i >= 0 || bucket_num+i <= len(target.Bytes)*8) &&
		len(ret) < count; i++ {

		if bucket_num-i >= 0 {
			bucket = rt.Buckets[bucket_num-i]
			copyToEntrySlice(&ret, bucket.Front(), count)
		}

		if bucket_num+i < len(target.Bytes)*8 {
			bucket = rt.Buckets[bucket_num+i]
			copyToEntrySlice(&ret, bucket.Front(), count)
		}

	}

	for _, e := range ret {
		e.distance = *e.ZifAddress.Xor(&target)
	}

	sort.Sort(Entries(ret))

	return ret
}
