// Kademlia

package dht

import (
	"container/list"
	"encoding/json"
	"io/ioutil"
	"os"
	"sort"

	log "github.com/sirupsen/logrus"
)

const MaxBucketSize = 20

type DhtFile struct {
	entryCount int
	kv         [][]KeyValue
}

type DHTSave struct {
	Buckets     [][]*KeyValue
	LongBuckets [][]*KeyValue
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

func (rt *RoutingTable) Save(filename string) error {
	save := DHTSave{
		make([][]*KeyValue, len(rt.LocalAddress.Bytes)*8),
		make([][]*KeyValue, len(rt.LocalAddress.Bytes)*8),
	}

	for n, b := range rt.Buckets {
		for i := b.Front(); i != nil; i = i.Next() {
			save.Buckets[n] = append(save.Buckets[n], i.Value.(*KeyValue))
		}
	}

	for n, b := range rt.LongBuckets {
		for i := b.Front(); i != nil; i = i.Next() {
			save.LongBuckets[n] = append(save.LongBuckets[n], i.Value.(*KeyValue))
		}
	}

	json, err := json.Marshal(save)

	log.Info(string(json))

	if err != nil {
		log.Error(err.Error())
		return err
	}

	err = ioutil.WriteFile(filename, json, 0600)

	if err != nil {
		return err
	}

	return nil
}

func LoadRoutingTable(path string, addr Address) (*RoutingTable, error) {
	var ret RoutingTable
	ret.Setup(addr)
	var save DHTSave

	if _, err := os.Stat(path); os.IsNotExist(err) {
		log.Info("Creating new routing table")
		return &ret, nil
	}

	data, err := ioutil.ReadFile(path)

	if err != nil {
		log.Info("An error occured, creating new routing table")
		return &ret, nil
	}

	err = json.Unmarshal(data, &save)

	if err != nil {
		log.Info("An error occured, creating new routing table")
		return &ret, nil
	}

	for n, i := range save.Buckets {
		for _, j := range i {
			ret.Buckets[n].PushBack(j)
		}
	}

	for n, i := range save.LongBuckets {
		for _, j := range i {
			ret.LongBuckets[n].PushBack(j)
		}
	}

	return &ret, nil
}

func (rt *RoutingTable) NumPeers() int {
	return BucketSize(rt.LongBuckets) + BucketSize(rt.Buckets)
}

func BucketSize(bucket []*list.List) int {
	count := 0

	for _, b := range bucket {
		for i := b.Front(); i != nil; i = i.Next() {
			count += 1
		}
	}

	return count
}

func (rt *RoutingTable) UpdateBucket(buckets []*list.List, kv *KeyValue) bool {
	if len(kv.Key.Bytes) < AddressBinarySize {
		return false
	}

	zero_count := kv.Key.Xor(&rt.LocalAddress).LeadingZeroes()
	bucket := buckets[zero_count]

	// TODO: Ping peers, starting from back. If none reply, remove them.
	// Ensures only active peers are stored.
	if bucket.Len() == MaxBucketSize {
		return false
	}

	var foundEntry *list.Element = nil
	for i := bucket.Front(); i != nil; i = i.Next() {
		if i.Value.(*KeyValue).Key.Equals(&kv.Key) {
			foundEntry = i
		}
	}

	if foundEntry == nil {
		bucket.PushFront(kv)
	} else {
		// Update the value as well
		copy(foundEntry.Value.(*KeyValue).Value, kv.Value)
		bucket.MoveToFront(foundEntry)
	}

	return true
}

func (rt *RoutingTable) Update(kv *KeyValue) bool {
	var success bool

	closest := rt.FindClosest(kv.Key, MaxBucketSize)

	// If this peer is the closest known peer, then store it.
	if len(closest) > 0 {

		nearest := true
		dist_this := rt.LocalAddress.Xor(&kv.Key)

		for _, i := range closest {
			dist := i.Key.Xor(&kv.Key)

			if !dist_this.Less(dist) {
				nearest = false
			}
		}

		if nearest {
			success = rt.UpdateBucket(rt.LongBuckets, kv)
		}

	} else if len(closest) == 0 {
		success = rt.UpdateBucket(rt.LongBuckets, kv)
	}

	success = rt.UpdateBucket(rt.Buckets, kv)

	return success
}

func copyToEntrySlice(slice *Pairs, begin *list.Element, count int) {

	for i := begin; i != nil && len(*slice) < count; i = i.Next() {
		*slice = append(*slice, i.Value.(*KeyValue))
	}

}

func (rt *RoutingTable) FindClosestInBuckets(buckets []*list.List, target Address, count int) Pairs {
	if len(target.Bytes) != AddressBinarySize {
		return nil
	}

	ret := make(Pairs, 0, count)

	bucket_num := target.Xor(&rt.LocalAddress).LeadingZeroes()
	bucket := buckets[bucket_num]

	copyToEntrySlice(&ret, bucket.Front(), count)

	// If the bucket is not filled, look the the buckets either side.
	for i := 1; (bucket_num-i >= 0 || bucket_num+i <= len(target.Bytes)*8) &&
		len(ret) < count; i++ {

		if bucket_num-i >= 0 {
			bucket = buckets[bucket_num-i]
			copyToEntrySlice(&ret, bucket.Front(), count)
		}

		if bucket_num+i < len(target.Bytes)*8 {
			bucket = buckets[bucket_num+i]
			copyToEntrySlice(&ret, bucket.Front(), count)
		}

	}

	for _, e := range ret {
		e.distance = *e.Key.Xor(&target)
	}

	sort.Sort(Pairs(ret))

	return ret
}

func (rt *RoutingTable) FindClosest(target Address, count int) Pairs {
	entries := make(Pairs, 0)

	entries = append(entries, rt.FindClosestInBuckets(rt.Buckets, target, count)...)
	entries = append(entries, rt.FindClosestInBuckets(rt.LongBuckets, target, count)...)

	sort.Sort(Pairs(entries))

	// then remove duplicates, as the two bucket lists may contain the same
	// entries

	if len(entries) < 1 {
		return entries
	}

	last := entries[0]
	j := 1
	for i := 1; i < len(entries); i++ {
		if entries[i].Key.Equals(&last.Key) {
			continue
		}

		entries[j] = entries[i]
		last = entries[i]
		j++
	}

	return entries[:j]
}
