// Kademlia

package main

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"sort"
)

const BucketSize = 20

type Entry struct {
	ZifAddress    Address
	Name          string
	Desc          string
	PublicAddress string
	PublicKey     []byte

	// The owner of this entry should have signed it, we need to store the
	// sigature. It's actually okay as we can verify that a peer owns a public
	// key by generating an address from it - if the address is not the peers,
	// then Mallory is just using someone elses entry for their own address.
	Signature []byte
	Port      int

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
	LongBuckets  []*list.List
}

func (rt *RoutingTable) Setup(addr Address) {
	rt.LocalAddress = addr

	rt.Buckets = make([]*list.List, len(rt.LocalAddress.Bytes)*8)
	rt.LongBuckets = make([]*list.List, len(rt.LocalAddress.Bytes)*8)

	for i := 0; i < len(rt.LocalAddress.Bytes)*8; i++ {
		rt.Buckets[i] = list.New()
		rt.LongBuckets[i] = list.New()
	}
}

func (rt *RoutingTable) SaveBuckets(buckets []*list.List, filename string) error {
	all_buckets := make([][]*Entry, 0, len(buckets))

	for _, b := range buckets {
		slice := make([]*Entry, b.Len())

		index := 0
		for i := b.Front(); i != nil; i = i.Next() {
			slice[index] = i.Value.(*Entry)
			index++
		}

		all_buckets = append(all_buckets, slice)
	}

	json, err := json.Marshal(all_buckets)

	if err != nil {
		return err
	}

	err = ioutil.WriteFile(filename, json, 0600)

	if err != nil {
		return err
	}

	return nil
}

func (rt *RoutingTable) NumPeers() int {
	count := 0

	for _, b := range rt.Buckets {
		for i := b.Front(); i != nil; i = i.Next() {
			count += 1
		}
	}

	for _, b := range rt.LongBuckets {
		for i := b.Front(); i != nil; i = i.Next() {
			count += 1
		}
	}

	go rt.SaveBuckets(rt.Buckets, "dht")

	return count
}

func (rt *RoutingTable) UpdateBucket(buckets []*list.List, entry Entry) bool {
	zero_count := entry.ZifAddress.Xor(&rt.LocalAddress).LeadingZeroes()
	bucket := buckets[zero_count]

	// TODO: Ping peers, starting from back. If none reply, remove them.
	// Ensures only active peers are stored.
	if bucket.Len() == BucketSize {
		return false
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

	return true
}

func (rt *RoutingTable) Update(entry Entry) bool {
	var success bool

	closest := rt.FindClosest(entry.ZifAddress, 1)
	if len(closest) == 1 {
		closest_entry := closest[0]

		dist_closest := closest_entry.ZifAddress.Xor(&entry.ZifAddress)
		dist_this := rt.LocalAddress.Xor(&entry.ZifAddress)

		if dist_this.Less(dist_closest) {
			success = rt.UpdateBucket(rt.LongBuckets, entry)
		}
	}

	success = rt.UpdateBucket(rt.Buckets, entry)

	return success
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
